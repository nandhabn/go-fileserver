package controllers

import (
	"net/http"
	"os"
)

func ServeIndex(uiDir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serve static files if they exist, otherwise fallback to `index.html`
		filePath := uiDir + r.URL.Path
		if _, err := os.Stat(filePath); os.IsNotExist(err) || r.URL.Path == "/" {
			http.ServeFile(w, r, filePath+"/index.html")
		} else {
			http.FileServer(http.Dir(uiDir)).ServeHTTP(w, r)
		}
	}
}
