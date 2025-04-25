package controllers

import (
	"net/http"
	"os"
	"path/filepath"
)

func DeleteVideo(videoDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := os.Remove(filepath.Join(videoDir, filepath.Base(r.PathValue("id"))))
		if err != nil {
			http.Error(w, "Failed to delete video", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
