package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("wordcrafter")
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
