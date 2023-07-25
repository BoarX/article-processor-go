package articles

import (
	"context"
	"fmt"
	"github.com/SkaisgirisMarius/article-processor/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestInsertArticlesToDatabaseInBatch(t *testing.T) {
	// Connect to the MongoDB database
	config.GetConfig("../conf_test.yaml")
	mongoURI := fmt.Sprintf("%s://%s:%s", config.Conf.MongoDb.DriverName, config.Conf.MongoDb.Host, config.Conf.MongoDb.Port)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer client.Disconnect(context.Background())

	// Create a collection for testing
	collection := getArticlesCollection()
	if err != nil {
		t.Fatal("Failed to create a collection for testing: ", err)
	}

	// Cleanup after the test
	defer func() {
		_, err := collection.DeleteMany(context.Background(), bson.M{})
		if err != nil {
			t.Fatal("Failed to clean up test data: ", err)
		}
	}()

	// Define the test articles
	article1 := Article{
		ArticleID: "123",
		Title:     "Test Article 1",
	}

	article2 := Article{
		ArticleID: "456",
		Title:     "Test Article 2",
	}

	article3 := Article{
		ArticleID: "789",
		Title:     "Test Article 3",
	}

	// Insert the articles to the database
	insertArticlesToDatabaseInBatch([]Article{article1, article2, article3})

	// Check if the articles were inserted correctly
	result, err := getArticleByIDFromDatabase("123")
	assert.NoError(t, err)
	assert.Equal(t, "Test Article 1", result.Title)

	// Call getArticleListFromDatabase() to get all articles
	articles, err := getArticleListFromDatabase()
	assert.NoError(t, err)

	// Check if the article count matches
	assert.Equal(t, 3, len(articles))
}
