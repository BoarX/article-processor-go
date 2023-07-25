package articles

import (
	"github.com/SkaisgirisMarius/article-processor/db"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const statusSuccess = "success"

//Internal Structures

type Article struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ArticleID   string             `bson:"articleID,omitempty" json:"articleID"`
	TeamID      string             `bson:"teamId" json:"teamId"`
	OptaMatchID *string            `bson:"optaMatchId,omitempty" json:"optaMatchId"`
	Title       string             `bson:"title" json:"title"`
	Type        []string           `bson:"type" json:"type"`
	Teaser      *string            `bson:"teaser" json:"teaser"`
	Content     string             `bson:"content" json:"content"`
	URL         string             `bson:"url" json:"url"`
	ImageURL    string             `bson:"imageUrl" json:"imageUrl"`
	GalleryURLs []string           `bson:"galleryUrls,omitempty" json:"galleryUrls,omitempty"`
	VideoURL    *string            `bson:"videoUrl,omitempty" json:"videoUrl,omitempty"`
	Published   time.Time          `bson:"published" json:"published"`
}

type SingleArticleResponse struct {
	Status string   `json:"status"`
	Data   *Article `json:"data"`
}

type MultipleArticlesResponse struct {
	Status string     `json:"status"`
	Data   []*Article `json:"data"`
}

// External XML structures

type ExternalArticleItem struct {
	ArticleURL        string  `xml:"ArticleURL"`
	NewsArticleID     string  `xml:"NewsArticleID"`
	PublishDate       string  `xml:"PublishDate"`
	Taxonomies        string  `xml:"Taxonomies"`
	TeaserText        *string `xml:"TeaserText"`
	Subtitle          string  `xml:"Subtitle"`
	ThumbnailImageURL string  `xml:"ThumbnailImageURL"`
	Title             string  `xml:"Title"`
	BodyText          string  `xml:"BodyText"`
	GalleryImageURLs  string  `xml:"GalleryImageURLs"`
	VideoURL          *string `xml:"VideoURL"`
	OptaMatchID       *string `xml:"OptaMatchId"`
	LastUpdateDate    string  `xml:"LastUpdateDate"`
	IsPublished       string  `xml:"IsPublished"`
}

type ExternalArticleItems struct {
	Items []ExternalArticleItem `xml:"NewsletterNewsItem"`
}

type ExternalArticleListData struct {
	ClubName            string               `xml:"ClubName"`
	ClubWebsiteURL      string               `xml:"ClubWebsiteURL"`
	NewsletterNewsItems ExternalArticleItems `xml:"NewsletterNewsItems"`
}

type ExternalArticleData struct {
	ClubName       string              `xml:"ClubName"`
	ClubWebsiteURL string              `xml:"ClubWebsiteURL"`
	NewsArticle    ExternalArticleItem `xml:"NewsArticle"`
}

// getMissingArticlesFromDatabase retrieves missing articles by comparing the given list of articles
// with the existing articles in the database.
func getMissingArticlesFromDatabase(articles []Article) ([]Article, error) {
	var articleIDs []string
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ArticleID)
	}
	collection := getArticlesCollection()
	ctx, _ := db.GetTimeoutContext()

	// Create a filter to find documents with ArticleIDs in the given list
	filter := bson.M{"articleID": bson.M{"$in": articleIDs}}

	// Execute the find operation with the filter
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Error("Error while querying the db: ", err)
		return nil, err
	}
	defer cur.Close(ctx)

	// Create a map to store existing ArticleIDs
	existingArticleIDs := make(map[string]bool)

	// Iterate through the cursor and store the existing ArticleIDs in the map
	for cur.Next(ctx) {
		var existingArticle Article
		if err := cur.Decode(&existingArticle); err != nil {
			return nil, err
		}
		existingArticleIDs[existingArticle.ArticleID] = true
	}

	// Create a slice to store the missing articles
	var missingArticles []Article

	// Iterate through the input articles and check which ones are not present in the existingArticleIDs map
	for _, article := range articles {
		if _, found := existingArticleIDs[article.ArticleID]; !found {
			// Article is missing, add it to the missingArticles slice
			missingArticles = append(missingArticles, article)
		}
	}

	// Range through the missingArticles and call getArticleByID for each of them
	var updatedMissingArticles []Article
	for _, missingArticle := range missingArticles {
		updatedArticle, err := getArticleByID(missingArticle.ArticleID)
		if err != nil {
			return nil, err
		}
		updatedMissingArticles = append(updatedMissingArticles, *updatedArticle)
	}

	// Return the updated missing articles
	return updatedMissingArticles, nil
}

// insertArticlesToDatabaseInBatch inserts a batch of articles into the database using bulk write operations.
// It checks if each article already exists in the database based on its ArticleID before adding them.
func insertArticlesToDatabaseInBatch(articles []Article) {
	collection := getArticlesCollection()
	ctx, _ := db.GetTimeoutContext()
	// Prepare the bulk write operations
	var bulkOps []mongo.WriteModel

	for _, article := range articles {
		// Create a filter to check if the article already exists in the database
		filter := bson.M{"articleID": article.ArticleID}

		// Create the update operation. Here, we use upsert to insert the document if it doesn't exist.
		update := bson.M{"$set": article}

		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		bulkOps = append(bulkOps, updateModel)
	}

	// Execute the bulk write operation
	_, err := collection.BulkWrite(ctx, bulkOps)
	if err != nil {
		log.Println("Failed to insert article to the DB: ", err)
		return
	}
	log.Println("Added ", len(bulkOps), " new articles.")
}

// getArticleByIDFromDatabase retrieves an article from the database based on its ID.
func getArticleByIDFromDatabase(id string) (*Article, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error("could not get primitive.ObjectID from provided id. ", err)
		return nil, err
	}
	var article Article
	filter := bson.M{"_id": objID}
	ctx, _ := db.GetTimeoutContext()
	singleResult := getArticlesCollection().FindOne(ctx, filter)
	if err := singleResult.Decode(&article); err != nil {
		log.Errorf("could not find article with ID: %v, error: %v", id, err)
		return nil, err
	}
	return &article, nil
}

// getArticleListFromDatabase retrieves a list of articles from the database.
func getArticleListFromDatabase() ([]*Article, error) {
	filter := bson.M{}
	ctx, _ := db.GetTimeoutContext()

	articles := make([]*Article, 0)
	result, err := getArticlesCollection().Find(ctx, filter)
	if err != nil {
		log.Error("Could not get articles from the database. Error: ", err)
		return nil, err
	}
	// Iterate through the cursor to process each retrieved article.
	for result.Next(ctx) {
		var a Article
		if err = result.Decode(&a); err != nil {
			log.Fatal(err)
		}
		articles = append(articles, &a)
	}
	return articles, nil
}

func getArticlesCollection() *mongo.Collection {
	return db.GetMongoCollection("articles")
}
