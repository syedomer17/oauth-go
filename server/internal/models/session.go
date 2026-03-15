package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id,omitempty" `
	RefreshToken  string        `bson:"refresh_token" `

	UserAgent string `bson:"user_agent" `
	IPAddress string `bson:"ip_address" `

	ExpiresAt time.Time `bson:"expires_at" `
	CreatedAt time.Time `bson:"created_at" `
}