package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Video struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	IsDownloaded bool   `json:"is_downloaded,omitempty"`
}

func ListVideos(videoDir string, downloading map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		videoDir, err := filepath.Abs(videoDir)
		if err != nil {
			http.Error(w, "Failed to resolve video directory", http.StatusInternalServerError)
			return
		}
		files, err := os.ReadDir(videoDir)
		if err != nil {
			http.Error(w, "Failed to list videos", http.StatusInternalServerError)
			return
		}

		sort.Slice(files, func(i, j int) bool {
			extractEpisodeNumber := func(name string) int {
				parts := strings.FieldsFunc(name, func(r rune) bool {
					return r < '0' || r > '9'
				})
				for _, part := range parts {
					if len(part) > 0 {
						if num, err := strconv.Atoi(part); err == nil {
							return num
						}
					}
				}
				return 0
			}
			return extractEpisodeNumber(files[i].Name()) < extractEpisodeNumber(files[j].Name())
		})

		var videos []Video
		for _, file := range files {
			switch filepath.Ext(file.Name()) {
			case ".mp4", ".webm", ".ogg", ".mkv", ".m3u8":
				video := Video{Name: file.Name(), Path: "/videos/" + file.Name(), IsDownloaded: downloading[file.Name()]}
				videos = append(videos, video)
			}
		}
		if videos == nil {
			json.NewEncoder(w).Encode([]Video{})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)
	}
}
