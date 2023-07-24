package main

import (
	"github.com/SkaisgirisMarius/article-processor/articles"
	"github.com/SkaisgirisMarius/article-processor/config"
	"github.com/SkaisgirisMarius/article-processor/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.GetConfig("conf.yaml")
}

func main() {
	log.Println("Starting Article Processor")

	r := server.NewRouter()

	articles.InitializeArticleRetriever()
	server.StartServer(r)
}
