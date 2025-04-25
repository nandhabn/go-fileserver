package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func AddServeCmd(startServer func(), rootCmd *cobra.Command, videoDir *string, uiDir *string, watchFile *string, port *string) {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve video files",
		Run: func(cmd *cobra.Command, args []string) {
			checkAbs(videoDir)
			checkAbs(uiDir)
			checkAbs(watchFile)

			if _, err := os.Stat(*videoDir); os.IsNotExist(err) {
				err := os.MkdirAll(*videoDir, os.ModePerm)
				if err != nil {
					fmt.Printf("Error creating video directory %s: %v\n", videoDir, err)
					os.Exit(1)
				}
			}
			startServer()
		},
	}
	rootCmd.Flags().StringVarP(port, "port", "p", "8080", "Port to run the server on")
	serveCmd.Flags().StringVarP(uiDir, "ui-dir", "u", "./frontend/build/client", "Directory to serve the UI from")
	serveCmd.Flags().StringVarP(videoDir, "directory", "d", "./videos", "Directory to serve videos from")
	serveCmd.Flags().StringVarP(watchFile, "watchfile", "w", "watched_videos.json", "File to store watched video data")
	rootCmd.AddCommand(serveCmd)

}

func checkAbs(variable *string) {
	if !filepath.IsAbs(*variable) {
		absPath, err := filepath.Abs(*variable)
		if err != nil {
			fmt.Printf("Error converting %s to absolute path: %v\n", *variable, err)
			os.Exit(1)
		}
		*variable = absPath
	}
}
