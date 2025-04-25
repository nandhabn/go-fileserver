package main

import "fileserver/controllers"

var (
	downloadQueue = make(chan string)
	downloading   = make(map[string]bool)
)

// Start the download queue in a separate goroutine
func backgroundServices() {
	sem := make(chan struct{}, 10) // Limit to 10 goroutines
	for episodeId := range downloadQueue {
		sem <- struct{}{}
		go func(id string) {
			defer func() { <-sem }()
			controllers.DownloadVideo(id, downloading)
		}(episodeId)
	}
}

func init() {
	go backgroundServices()
}
