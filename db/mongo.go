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

func MongoConnect() *mongo.Client {
	if mongoDb != nil {
		return mongoDb
	}
	mongoURI := fmt.Sprintf("%s://%s:%s/%s", config.Conf.MongoDb.DriverName, config.Conf.MongoDb.Host, config.Conf.MongoDb.Port, config.Conf.MongoDb.DbName)
	clientOptions := options.Client().ApplyURI(mongoURI)

	mongoDb, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = mongoDb.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
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
