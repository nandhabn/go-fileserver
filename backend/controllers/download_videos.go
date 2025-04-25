package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/grafov/m3u8"
)

// Path to the JSON file
var jsonFilePath = "./files/1p.json"

type URLList map[string]string

func DownloadNext10Handler(downloadQueue chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var lastDownloaded struct {
			LastDownloaded string `json:"lastDownloaded,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&lastDownloaded); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if lastDownloaded.LastDownloaded == "" {
			lastDownloaded.LastDownloaded = "0"
		}

		episodeId, err := strconv.Atoi(lastDownloaded.LastDownloaded)
		if err != nil {
			http.Error(w, "Invalid last downloaded episode ID "+lastDownloaded.LastDownloaded, http.StatusBadRequest)
			return
		}

		// Add the next 10 videos to the download queue
		for i := episodeId; i < episodeId+10; i++ {
			episodeId := strconv.Itoa(i + 1)
			downloadQueue <- episodeId
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Next 10 videos added to download queue successfully"))
	}
}

func DownloadVideoHandler(downloadQueue chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var video struct {
			EpisodeId string `json:"episode_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add the video URL to the download queue
		downloadQueue <- video.EpisodeId

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Video added to download queue successfully"))
	}
}

func DownloadVideo(episode_id string, downloading map[string]bool) {
	// Open the JSON file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Printf("Error opening JSON file: %v\n", err)
		return
	}
	defer file.Close()

	// Decode the JSON file
	var urlList URLList
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&urlList); err != nil {
		fmt.Printf("Error decoding JSON file: %v\n", err)
		return
	}

	// Create a directory to save the videos
	outputDir := "downloaded_videos"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	url := urlList[episode_id]
	fileName := "one piece EP-" + episode_id
	downloading[episode_id] = true
	outputPath := filepath.Join(outputDir, fileName)
	if err := downloadVideo(url, outputPath); err != nil {
		fmt.Printf("Error downloading video %s: %v\n", episode_id, err)
	}
	delete(downloading, episode_id)
}

func downloadVideo(url, outputPath string) error {
	// Send GET request
	resp, err := fetchURL(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is a valid M3U8 playlist
	if resp.Header.Get("Content-Type") == "application/vnd.apple.mpegurl" {
		playlist := getPlaylistUrl(resp, url)
		if playlist == "" {
			return fmt.Errorf("failed to get playlist URL")
		}
		return downloadMediaPlaylist(playlist, outputPath)
	}

	outputPath += getExtensionByContentType(url)

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	return nil
}

func fetchURL(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	return resp, nil
}

func getExtensionByContentType(contentType string) string {
	switch contentType {
	case "video/mp4":
		return "-.mp4"
	case "video/webm":
		return "-.webm"
	case "video/ogg":
		return "-.ogg"
	case "video/x-matroska":
		return "-.mkv"
	default:
		return "-.mp4" // Default to .mp4 if content type is unknown
	}
}

func getPlaylistUrl(res *http.Response, url string) string {
	p, listType, err := m3u8.DecodeFrom(res.Body, true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		return url
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		for _, v := range masterpl.Variants {
			if v.Resolution == "1920x1080" {
				resp, err := fetchURL(v.URI)
				if err != nil {
					fmt.Println("Error fetching playlist URL:", err)
					return ""
				}
				defer resp.Body.Close()
				return getPlaylistUrl(resp, v.URI)
			}
		}
	}
	return ""
}

func downloadMediaPlaylist(playlist string, outputPath string) error {
	// Create a directory to store the segments
	segmentsDir := outputPath + "_segments"
	if err := os.MkdirAll(segmentsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create segments directory: %w", err)
	}

	// Fetch the playlist
	resp, err := fetchURL(playlist)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	// Parse the playlist
	p, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return fmt.Errorf("failed to decode playlist: %w", err)
	}

	if listType != m3u8.MEDIA {
		return fmt.Errorf("expected a media playlist but got a different type")
	}

	mediaPlaylist := p.(*m3u8.MediaPlaylist)
	newMP, err := m3u8.NewMediaPlaylist(mediaPlaylist.WinSize(), mediaPlaylist.Count())
	newMP.SeqNo = mediaPlaylist.SeqNo
	newMP.TargetDuration = mediaPlaylist.TargetDuration
	if err != nil {
		return fmt.Errorf("failed to create new media playlist: %w", err)
	}

	// Download each segment
	for _, segment := range mediaPlaylist.Segments {
		if segment == nil {
			continue
		}

		segmentURL := segment.URI
		segmentFileName := filepath.Join(segmentsDir, filepath.Base(segmentURL))
		encodedSegmentFileName := filepath.Join("/videos/"+url.PathEscape(strings.Split(segmentsDir, "/")[1]), filepath.Base(segmentURL))

		// Fetch and save the segment
		segmentResp, err := fetchURL(segmentURL)
		if err != nil {
			return fmt.Errorf("failed to fetch segment %s: %w", segmentURL, err)
		}
		defer segmentResp.Body.Close()

		newMP.Append(encodedSegmentFileName, segment.Duration, segment.Title)

		outFile, err := os.Create(segmentFileName)
		if err != nil {
			return fmt.Errorf("failed to create segment file %s: %w", segmentFileName, err)
		}

		_, err = io.Copy(outFile, segmentResp.Body)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to save segment %s: %w", segmentFileName, err)
		}
	}

	// Create a new playlist file to serve the downloaded segments
	localPlaylistPath := outputPath + ".m3u8"
	localPlaylistFile, err := os.Create(localPlaylistPath)
	if err != nil {
		return fmt.Errorf("failed to create local playlist file: %w", err)
	}
	defer localPlaylistFile.Close()

	// Create a new playlist file to serve the downloaded segments
	localPlaylistPath1 := outputPath + ".m3u8"
	localPlaylistFile1, err := os.Create(localPlaylistPath1)
	if err != nil {
		return fmt.Errorf("failed to create local playlist file: %w", err)
	}
	defer localPlaylistFile1.Close()

	buf := newMP.Encode()
	buf.WriteTo(localPlaylistFile1)

	buf = newMP.Encode()
	buf.WriteTo(localPlaylistFile)
	fmt.Printf("Playlist saved to %s\n", localPlaylistPath)

	return nil
}
