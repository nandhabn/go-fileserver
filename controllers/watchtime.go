package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type WatchDetail struct {
	watched   bool
	watchTime int
}

var (
	upgrader       = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	lastWatchedMap = make(map[string]WatchDetail) // Map to store last watched time for each video
)

func sendInitialData(conn *websocket.Conn) error {
	mu.Lock()
	defer mu.Unlock()

	err := conn.WriteJSON(lastWatchedMap)
	if err != nil {
		return fmt.Errorf("error sending initial data: %v", err)
	}

	return nil
}

func HandleWebSocket(watchFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading to WebSocket:", err)
			return
		}

		defer func() {
			saveDataToFile(watchFile)
			conn.Close()
		}()

		err = sendInitialData(conn)
		if err != nil {
			fmt.Println("Error sending initial data:", err)
			return
		}

		for {
			var message struct {
				VideoID string  `json:"videoId"`
				Time    float64 `json:"time"`
			}

			// Read message from client
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println("Error reading JSON:", err)
				break
			}
			videoDetails := lastWatchedMap[message.VideoID]
			videoDetails.watchTime = int(message.Time)

			// Store the last watched time
			mu.Lock()
			lastWatchedMap[message.VideoID] = videoDetails
			mu.Unlock()

			fmt.Printf("Updated last watched time for video %s: %f seconds\n", message.VideoID, message.Time)
		}
	}
}

func saveDataToFile(filename string) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	for videoID, details := range lastWatchedMap {
		_, err := file.WriteString(fmt.Sprintf("%s,%d,%t\n", videoID, details.watchTime, details.watched))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
