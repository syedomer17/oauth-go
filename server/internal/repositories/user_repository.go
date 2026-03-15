package repositories

import (
	"context"
	"oauth/internal/database"
	"oauth/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct{}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	collection := database.DB.Collection("users")

	var user models.User

	err := collection.FindOne(context.Background(), bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	collection := database.DB.Collection("users")

	user.CreatedAt = time.Now()

	user.ID = primitive.ObjectID{}

	_, err := collection.InsertOne(context.Background(), user)

	return err
}

// FindByID fetches a single user by their MongoDB ObjectID.
func (r *UserRepository) FindByID(id primitive.ObjectID) (*models.User, error) {
	collection := database.DB.Collection("users")

	var user models.User

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateOAuthProfile keeps user profile fields in sync with the OAuth provider.
func (r *UserRepository) UpdateOAuthProfile(userID primitive.ObjectID, name, avatar string, provider models.Provider, providerID, username string) error {
	collection := database.DB.Collection("users")

	updates := bson.M{
		"name":        name,
		"avatar":      avatar,
		"provider":    provider,
		"provider_id": providerID,
	}

	if username != "" {
		updates["username"] = username
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": updates},
	)

	return err
}
