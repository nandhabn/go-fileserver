package main

import (
	"fileserver/controllers"
	"fmt"
	"log"
	"net/http"
)

func startServer() {
	// Frontend Mux for serving React build and static files
	frontendMux := http.NewServeMux()
	frontendMux.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(videoDir))))
	frontendMux.HandleFunc("/", controllers.ServeIndex(uiDir)) // Serve index.html for root path

	// API Mux for handling API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/ws", controllers.HandleWebSocket(watchFile))
	apiMux.HandleFunc("/api/videos/download-next-10", controllers.DownloadNext10Handler(downloadQueue))
	apiMux.HandleFunc("/api/videos", controllers.ListVideos(videoDir, downloading))
	apiMux.HandleFunc("/api/videos/{id}", controllers.DeleteVideo(videoDir))
	apiMux.HandleFunc("/api/videos/download", controllers.DownloadVideoHandler(downloadQueue))

	// Root Mux to combine frontend and API
	rootMux := http.NewServeMux()
	rootMux.Handle("/api/", apiMux)  // Route API requests to apiMux
	rootMux.Handle("/", frontendMux) // Route all other requests to frontendMux

	fmt.Printf("Serving on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, rootMux))
}
