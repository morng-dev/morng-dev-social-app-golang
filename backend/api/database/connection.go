package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoURI := "mongodb://localhost:27017"

	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err == nil {
		fmt.Println("Connected to MongoDB")
		DB = Client.Database("social-app")
	} else {
		fmt.Println("Failed to connect to MongoDB")
		return err
	}
	return nil
}
