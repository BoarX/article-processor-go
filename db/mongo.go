package db

import (
	"context"
	"fmt"
	"github.com/SkaisgirisMarius/article-processor/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var mongoDb *mongo.Client

const MongoTimeout = 15 * time.Second

// MongoConnect establishes a connection to the MongoDB database and returns the client instance.
func MongoConnect() *mongo.Client {
	if mongoDb != nil {
		return mongoDb
	}

	// Construct the MongoDB connection URI using the configuration settings.
	mongoURI := fmt.Sprintf("%s://%s:%s/%s", config.Conf.MongoDb.DriverName, config.Conf.MongoDb.Host, config.Conf.MongoDb.Port, config.Conf.MongoDb.DbName)
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to the MongoDB server using the specified client options.
	mongoDb, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection to the MongoDB server is successful.
	err = mongoDb.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return mongoDb
}

func getDbName() string {
	dbName := config.Conf.MongoDb.DbName
	return dbName
}

func getInCrowdDb() *mongo.Database {
	return MongoConnect().Database(getDbName())
}

func GetMongoCollection(name string) *mongo.Collection {
	return getInCrowdDb().Collection(name)
}

func GetTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), MongoTimeout)
}
