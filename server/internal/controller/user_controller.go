package controller

import (
	"net/http"
	"oauth/internal/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMe returns the authenticated user's profile.
// The userID is injected by AuthRequired middleware — no token parsing here.
func GetMe(c *gin.Context) {
	userIDHex, _ := c.Get("userID")

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	repo := repositories.UserRepository{}
	user, err := repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
