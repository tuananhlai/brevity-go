package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brevity",
	Short: "Brevity is a content management system",
	Long:  `A modern content management system with AI-powered features.`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(generateArticleCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Brevity server",
	Long:  `Start the Brevity web server that serves the API and web interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

var generateArticleCmd = &cobra.Command{
	Use:   "generate-article",
	Short: "Generate a new article using AI",
	Long:  `Generate a new article using AI and save it to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGenerateArticle()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
