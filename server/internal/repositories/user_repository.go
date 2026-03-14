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