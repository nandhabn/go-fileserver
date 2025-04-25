package main

import (
	"fileserver/cmd"
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
var watchFile string
var uiDir string

func init() {
	cmd.AddGetUrlsCmd(rootCmd)
	cmd.AddServeCmd(startServer, rootCmd, &videoDir, &uiDir, &watchFile, &port)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
