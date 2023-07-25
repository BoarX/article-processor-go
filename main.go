package main

import (
	"github.com/SkaisgirisMarius/article-processor/articles"
	"github.com/SkaisgirisMarius/article-processor/config"
	"github.com/SkaisgirisMarius/article-processor/server"
	log "github.com/sirupsen/logrus"
)

// init initializes the application configuration by reading it from the "conf.yaml" file.
func init() {
	config.GetConfig("conf.yaml")
}

func main() {
	log.Println("Starting Article Processor")

	// Create a new router for handling HTTP requests
	r := server.NewRouter()

	// Initialize the article retriever to periodically fetch new articles
	articles.InitializeArticleRetriever()

	// Start the HTTP server with the provided router
	server.StartServer(r)
}
