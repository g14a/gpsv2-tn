package db

import (
	"context"
	"fmt"
	"gitlab.com/gps2.0/config"
	"gitlab.com/gps2.0/server"
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
	collection := mongoClient.Database("gpsV2").Collection(collectionName)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancelFunc()

	return collection, ctx
}

func connectDBOfficial() {
	appConfigInstance := config.GetAppConfig()

	fmt.Println(appConfigInstance.Mongoconfig.URL)

	mClient, err := mongo.NewClient(&options.ClientOptions{
		Hosts: []string{appConfigInstance.Mongoconfig.URL},
	})

	server.CheckError(err)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mClient.Connect(ctx)

	server.CheckError(err)
	mongoClient = mClient
}
