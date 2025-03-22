package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Video struct {
	Name string
	Path string
}

type PageData struct {
	Videos []Video
}

func listVideos(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	var videos []Video
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp4" || filepath.Ext(file.Name()) == ".webm" || filepath.Ext(file.Name()) == ".ogg" || filepath.Ext(file.Name()) == ".mkv" {
			videos = append(videos, Video{Name: file.Name(), Path: "/videos/" + file.Name()})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func servePage(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		http.Error(w, "Failed to list videos", http.StatusInternalServerError)
		return
	}

	var videos []Video
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".mp4" || filepath.Ext(file.Name()) == ".webm" || filepath.Ext(file.Name()) == ".ogg" || filepath.Ext(file.Name()) == ".mkv" {
			videos = append(videos, Video{Name: file.Name(), Path: "/videos/" + file.Name()})
		}
	}
	tmpl := template.Must(template.New("page").Parse(pageTemplate))
	tmpl.Execute(w, PageData{Videos: videos})
}

func startServer() {
	http.HandleFunc("/videos", listVideos)
	http.HandleFunc("/delete", deleteVideo)
	http.HandleFunc("/", servePage)
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(videoDir))))

	fmt.Printf("Serving on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
