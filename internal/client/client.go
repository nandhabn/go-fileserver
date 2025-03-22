package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const pageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Video Gallery</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; }
        .container { max-width: 800px; margin: auto; }
        video { width: 100%; margin-top: 10px; }
    </style>
</head>
<body>
    <h1>Video Gallery</h1>
    <div class="container">
        {{range .Videos}}
        <video controls>
            <source src="{{.}}" type="video/mp4">
            <source src="{{.}}" type="video/webm">
            <source src="{{.}}" type="video/ogg">
            <source src="{{.}}" type="video/x-matroska">
            Your browser does not support the video tag.
        </video>
        {{end}}
    </div>
</body>
</html>
`

type PageData struct {
	Videos []string
}

func main() {
	dir := "./videos" // Change this to your video directory
	port := "8081"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(dir)
		if err != nil {
			http.Error(w, "Failed to list videos", http.StatusInternalServerError)
			return
		}

		var videos []string
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".mp4" || filepath.Ext(file.Name()) == ".webm" || filepath.Ext(file.Name()) == ".ogg" || filepath.Ext(file.Name()) == ".mkv" {
				videos = append(videos, "/videos/"+file.Name())
			}
		}

		tmpl := template.Must(template.New("page").Parse(pageTemplate))
		tmpl.Execute(w, PageData{Videos: videos})
	})

	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(dir))))

	fmt.Printf("Serving on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
