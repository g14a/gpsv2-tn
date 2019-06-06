package db

import (
	"context"
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/errcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

func getMongoClient() *mongo.Client {
	once.Do(func() {
		connectDBOfficial()
	})

	return mongoClient
}

func GetMongoCollectionWithContext(collectionName string) (*mongo.Collection, context.Context) {
	mongoClient = getMongoClient()
	collection := mongoClient.Database("gpsgolang").Collection(collectionName)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return collection, ctx
}

func connectDBOfficial() {
	appConfigInstance := config.GetAppConfig()

	mClient, err := mongo.NewClient(options.Client().ApplyURI(appConfigInstance.Mongoconfig.URL))

	errcheck.CheckError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mClient.Connect(ctx)

	errcheck.CheckError(err)
	mongoClient = mClient
}
