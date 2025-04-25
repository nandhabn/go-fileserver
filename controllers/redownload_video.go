package controllers

import (
	"net/http"
	"os"
	"path/filepath"
)

func RedownloadVideo(videoDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		videoName := r.URL.Query().Get("name")
		if videoName == "" {
			http.Error(w, "Missing video name", http.StatusBadRequest)
			return
		}

		videoPath := filepath.Join(videoDir, videoName)
		if _, err := os.Stat(videoPath); os.IsNotExist(err) {
			http.Error(w, "Video not found", http.StatusNotFound)
			return
		}

		err := os.Remove(videoPath)
		if err != nil {
			http.Error(w, "Failed to delete video", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
