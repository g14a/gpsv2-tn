// Package db assists in getting mongo clients, connections
// and contexts of the databases needed
package db

import (
	"context"
	"gitlab.com/gpsv2-withoutrm/config"
	"gitlab.com/gpsv2-withoutrm/errorcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	mongoClient        *mongo.Client
	once               sync.Once
	liveDB = config.GetAppConfig().Mongoconfig.DBName
	rawDB  = config.GetAppConfig().Mongoconfig.RawDB
)

// getMongoClient returns the client
func getMongoClient() *mongo.Client {
	once.Do(func() {
		connectLiveDB()
	})

	return mongoClient
}

// GetMongoCollectionWithContext is mostly used by other functions
// to insert data into the collection which this function returns.
func GetMongoCollectionWithContext(collectionName string) (*mongo.Collection, context.Context) {
	mongoClient = getMongoClient()
	collection := mongoClient.Database(liveDB).Collection(collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	return collection, ctx
}

func GetRawCollectionWithContext() (*mongo.Collection, context.Context) {
	mongoClient = getMongoClient()
	collection := mongoClient.Database(rawDB).Collection("raw_data")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	return collection, ctx
}

// connectDBOfficial connects to the Mongo db and assigns a Mongo client
// using the db url mentioned in the config file.
func connectLiveDB() {
	appConfigInstance := config.GetAppConfig()

	mClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.Mongoconfig.URL))

	err = mClient.Connect(context.TODO())
	errorcheck.CheckError(err)

	mongoClient = mClient
}
