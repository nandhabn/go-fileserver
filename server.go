package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "videoserver",
	Short: "A simple video file server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var videoDir string
var port string

func init() {
	rootCmd.Flags().StringVarP(&videoDir, "directory", "d", "./videos", "Directory to serve videos from")
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
