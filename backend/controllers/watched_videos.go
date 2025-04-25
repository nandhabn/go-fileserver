package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

var (
	watchedVideos = make(map[string]bool)
	watchedFile   = "watched_videos.json"
	mu            sync.Mutex
)

func init() {
	loadWatchedVideos()
}

func loadWatchedVideos() {
	file, err := os.Open(watchedFile)
	if err != nil {
		if !os.IsNotExist(err) {
			panic("Failed to open watched videos file: " + err.Error())
		}
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&watchedVideos); err != nil {
		panic("Failed to decode watched videos file: " + err.Error())
	}
}

func saveWatchedVideos() {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(watchedFile)
	if err != nil {
		panic("Failed to save watched videos file: " + err.Error())
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(watchedVideos); err != nil {
		panic("Failed to encode watched videos to file: " + err.Error())
	}
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var video Video
	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	watchedVideos[video.Path] = true
	mu.Unlock()

	saveWatchedVideos()
	w.WriteHeader(http.StatusOK)
}

func GetWatchedVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(watchedVideos)
}
