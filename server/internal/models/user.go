package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderGitHub Provider = "github"
	ProviderLocal  Provider = "local"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Avatar   string             `bson:"avatar" json:"avatar"`
	Username string             `bson:"username,omitempty" json:"username,omitempty"`

	Provider   Provider `bson:"provider" json:"provider"`
	ProviderID string   `bson:"provider_id" json:"provider_id"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
