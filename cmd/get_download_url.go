package cmd

import (
	"encoding/json"
	"fileserver/internal"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func AddGetUrlsCmd(rootCmd *cobra.Command) {
	var indexDir string
	getUrlsCmd := &cobra.Command{
		Use:   "get-urls",
		Short: "Get download URLs for videos",
		Run: func(cmd *cobra.Command, args []string) {
			updateAPI(indexDir)
		},
	}
	getUrlsCmd.Flags().StringVarP(&indexDir, "index-dir", "i", "./index", "Directory to save the index file")
	rootCmd.AddCommand(getUrlsCmd)
}

func updateAPI(file string) {
	animeID, err := findAnimeID("one piece", "sub", "1P")
	if err != nil {
		fmt.Println("Error finding anime ID:", err)
		return
	}

	urlsList, err := getAllEpisodeLinks(animeID)
	if err != nil {
		fmt.Printf("Error fetching episode links: %v\n", err)
		return
	}

	if err := saveURLsToFile(file, urlsList); err != nil {
		fmt.Printf("Error saving URLs to file: %v\n", err)
	}
}

func findAnimeID(title, language, keyword string) (string, error) {
	anime, err := internal.SearchAnime(title, language)
	if err != nil {
		return "", fmt.Errorf("error searching anime: %w", err)
	}

	for key, value := range anime {
		if strings.Contains(value, keyword) {
			return key, nil
		}
	}

	return "", fmt.Errorf("anime ID not found")
}

func getAllEpisodeLinks(id string) (map[string]string, error) {
	episodeList, err := internal.EpisodesList(id, "sub")
	if err != nil {
		return nil, fmt.Errorf("error fetching episodes list: %w", err)
	}

	downloadURLs := make(map[string]string, len(episodeList))
	retryCount := 0

	isRoundNumber := 1
	if len(episodeList)%20 == 0 {
		isRoundNumber = 0
	}
	c := make(chan []string)

	parallelAPI := 60

	for i := range (len(episodeList) / parallelAPI) + isRoundNumber {
		if i%2 == 0 {
			fmt.Println("Next set of episode, count:", len(downloadURLs))
		}
		for j := range parallelAPI {
			if i*parallelAPI+j+1 >= len(episodeList) {
				break
			}
			go getDownloadURL(id, c, i*parallelAPI+j+1, retryCount)
		}

		for url := range c {
			if url[1] == "" {
				fmt.Println("Failed to get download URL for episode:", i+1)
				return downloadURLs, nil
			}
			downloadURLs[url[0]] = url[1]
			if len(downloadURLs) == i*parallelAPI+parallelAPI || len(downloadURLs) >= len(episodeList)-1 {
				break
			}
		}
		retryCount = 0 // Reset retry count on success
	}

	return downloadURLs, nil
}

func getDownloadURL(id string, c chan []string, episode, retry int) {
	downloadURL, err := internal.GetEpisodeURL(id, episode)
	if err != nil || !strings.Contains(downloadURL[0], ".m3u8") {
		u := downloadURL[0]
		if retry < 5 {
			fmt.Println("Error fetching download URL, retrying:", episode)
			time.Sleep(10 * time.Second)
			getDownloadURL(id, c, episode, retry+1)
		} else {
			fmt.Println("Max retry count reached for episode:", episode)
			c <- []string{strconv.Itoa(episode), u}
		}
	} else {
		c <- []string{strconv.Itoa(episode), downloadURL[0]}
	}
}

func saveURLsToFile(filepath string, urlsList map[string]string) error {
	marshalled, err := json.Marshal(urlsList)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	dir := filepath[:strings.LastIndex(filepath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	if err := os.WriteFile(filepath, marshalled, 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Println("URLs saved to file:", filepath)
	return nil
}
