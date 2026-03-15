package database

import (
	"context"
	"log"

	// "oauth/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(db *mongo.Database) {

	userCollection := db.Collection("users")
	sessionCollection := db.Collection("sessions")

	_, err := userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"expires_at": 1,
			},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	)

	if err != nil {
		log.Println("Index creation error:", err)
	}

	sessionRefreshTokenIndex := mongo.IndexModel{
		Keys:    map[string]int{"refresh_token": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = sessionCollection.Indexes().CreateOne(context.Background(), sessionRefreshTokenIndex)
	if err != nil {
		log.Println("Session refresh token index creation error:", err)
	}

	sessionExpiresTTLIndex := mongo.IndexModel{
		Keys: bson.M{
			"expires_at": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err = sessionCollection.Indexes().CreateOne(context.Background(), sessionExpiresTTLIndex)
	if err != nil {
		log.Println("Session TTL index creation error:", err)
	}
}
