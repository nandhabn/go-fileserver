package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func deleteVideo(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := os.Remove(filepath.Join(videoDir, filepath.Base(request.Path)))
	if err != nil {
		http.Error(w, "Failed to delete video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
