package database

import (
	"context"
	"log"

	// "oauth/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(db *mongo.Database) {

	userCollection := db.Collection("users")

	index := mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(context.Background(), index)

	if err != nil {
		log.Println("Index creation error:", err)
	}
}