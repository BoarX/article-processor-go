package articles

import (
	"encoding/xml"
	"fmt"
	"github.com/SkaisgirisMarius/article-processor/config"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

// InitializeArticleRetriever sets up and starts the article retrieval scheduler.
// It schedules the getNewArticles function to run at the specified interval defined in the configuration.
func InitializeArticleRetriever() {
	log.Println("Initializing article scheduler")
	scheduler := gocron.NewScheduler()
	log.Println(config.Conf.ArticleInterval)

	err := scheduler.Every(uint64(config.Conf.ArticleInterval)).Seconds().Do(getNewArticles)
	if err != nil {
		log.Fatal("Error scheduling job:", err)
		return
	}
	scheduler.Start()
}

// getNewArticles fetches the latest article list from the specified ArticleListURL,
// reads the XML content, identifies missing articles from the database,
// and inserts them in batch if there are any new articles.
func getNewArticles() {
	log.Println("Scanning for new articles")
	response, err := http.Get(config.Conf.ArticleListURL)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return
	}
	defer response.Body.Close()

	bodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	articleList, err := readXMLContent(string(bodyContent))
	if err != nil {
		log.Println("Error reading XML: ", err)
		return
	}
	missingArticles, err := getMissingArticlesFromDatabase(articleList)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	if len(missingArticles) > 0 {
		insertArticlesToDatabaseInBatch(missingArticles)
	} else {
		log.Println("There are no new articles to be added")
	}
}

// readXMLContent takes the XML content as input, unmarshals it into the ExternalArticleListData struct,
// and transforms the received XML feeds into a slice of Article structs.
func readXMLContent(xmlContent string) ([]Article, error) {
	var result ExternalArticleListData
	err := xml.Unmarshal([]byte(xmlContent), &result)
	if err != nil {
		log.Println("Failed to unmarshal XML: ", err)
		return nil, err
	}

	// Transform the received XML feeds into the Article struct
	var articles []Article
	for _, item := range result.NewsletterNewsItems.Items {
		publishDate, err := time.Parse("2006-01-02 15:04:05", item.PublishDate)
		if err != nil {
			log.Println("Failed to parse publish date: ", err)
			return nil, err
		}
		article := Article{
			ArticleID:   item.NewsArticleID,
			TeamID:      result.ClubName,
			OptaMatchID: nil,
			Title:       item.Title,
			Type:        []string{item.Taxonomies},
			Teaser:      item.TeaserText,
			URL:         item.ArticleURL,
			ImageURL:    item.ThumbnailImageURL,
			Published:   publishDate,
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// readXMLContentForSingleArticle takes the XML content for a single article as input, unmarshals it into the ExternalArticleData struct,
// and creates an Article struct from the parsed data.
func readXMLContentForSingleArticle(xmlContent string) (*Article, error) {
	var result *ExternalArticleData
	err := xml.Unmarshal([]byte(xmlContent), &result)
	if err != nil {
		log.Println("Failed to unmarshal XML: ", err)
		return nil, err
	}

	// Check if there are any items in the XML
	if result == nil {
		return nil, fmt.Errorf("no articles found in the XML")
	}

	publishDate, err := time.Parse("2006-01-02 15:04:05", result.NewsArticle.PublishDate)
	if err != nil {
		log.Println("Failed to parse publish date: ", err)
		return nil, err
	}

	// Create and return the Article struct
	article := &Article{
		ArticleID:   result.NewsArticle.NewsArticleID,
		TeamID:      result.ClubName,
		OptaMatchID: result.NewsArticle.OptaMatchID,
		Title:       result.NewsArticle.Title,
		Type:        []string{result.NewsArticle.Taxonomies},
		Teaser:      result.NewsArticle.TeaserText,
		Content:     result.NewsArticle.BodyText,
		URL:         result.NewsArticle.ArticleURL,
		ImageURL:    result.NewsArticle.ThumbnailImageURL,
		GalleryURLs: []string{result.NewsArticle.GalleryImageURLs},
		VideoURL:    result.NewsArticle.VideoURL,
		Published:   publishDate,
	}

	return article, nil
}

// getArticleByID retrieves an article from the provided API URL by sending a GET request with the given ID.
// It reads the XML response body, parses it into an Article struct using the readXMLContentForSingleArticle function,
// and returns the resulting article or an error if any occurred during the process.
func getArticleByID(id string) (*Article, error) {
	response, err := http.Get(config.Conf.ArticleURL + id)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body content as a string
	bodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	article, err := readXMLContentForSingleArticle(string(bodyContent))
	if err != nil {
		log.Println("Error reading XML: ", err)
		return nil, err
	}
	return article, nil
}
