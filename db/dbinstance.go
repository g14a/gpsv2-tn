package db

import (
	"context"
	"gitlab.com/gpsv2/config"
	"gitlab.com/gpsv2/errcheck"
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

func getMongoClient() *mongo.Client {
	once.Do(func() {
		connectDBOfficial()
	})

	return mongoClient
}

func getHistoryMongoClient() *mongo.Client {
	once.Do(func() {
		connectDBOfficial()
	})

	return historyMongoClient
}

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

func connectDBOfficial() {
	appConfigInstance := config.GetAppConfig()

	mClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.Mongoconfig.URL))
	historyClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.HistoryMongoConfig.BackupURL))

	errcheck.CheckError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// connect to live db
	err = mClient.Connect(ctx)
	errcheck.CheckError(err)

	// connect to history db
	err = historyClient.Connect(ctx)
	errcheck.CheckError(err)

	mongoClient = mClient
	historyMongoClient = historyClient

}
