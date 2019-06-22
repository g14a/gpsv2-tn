package db

import (
	"context"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errorcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	mongoClient        *mongo.Client
	historyMongoClient *mongo.Client
	once               sync.Once
)

// getMongoClient returns the client.
func getMongoClient() *mongo.Client {
	once.Do(func() {
		connectDBOfficial()
	})

	return mongoClient
}

func GetSessionFromClient() (mongo.Session, error) {
	once.Do(func() {
		connectDBOfficial()
	})

	session, err := mongoClient.StartSession()

	return session, err
}

// getHistoryMongoClient returns the history db client.
func getHistoryMongoClient() *mongo.Client {
	once.Do(func() {
		connectDBOfficial()
	})

	return historyMongoClient
}

// GetMongoCollectionWithContext is mostly used by other functions
// to insert data into the collection which this function returns.
func GetMongoCollectionWithContext(collectionName string) (*mongo.Collection, context.Context) {
	mongoClient = getMongoClient()
	collection := mongoClient.Database("gpsgolang").Collection(collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return collection, ctx
}

func GetHistoryCollectionsWithContext(collectionName string) (*mongo.Collection, context.Context) {
	historyMongoClient = getHistoryMongoClient()
	collection := historyMongoClient.Database("gpsgolang").Collection(collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return collection, ctx
}

// connectDBOfficial connects to the Mongo db and assigns a Mongo client
// using the db url mentioned in the config file.
func connectDBOfficial() {
	appConfigInstance := config.GetAppConfig()

	mClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.Mongoconfig.URL))
	historyClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.HistoryMongoConfig.BackupURL))

	errorcheck.CheckError(err)

	// Pings the database for a max of 10 seconds. Afterwards it gives an error.
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// connect to live db
	err = mClient.Connect(ctx)
	errorcheck.CheckError(err)

	// connect to history db
	err = historyClient.Connect(ctx)
	errorcheck.CheckError(err)

	mongoClient = mClient
	historyMongoClient = historyClient

}
