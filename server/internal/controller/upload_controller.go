package controller

import (
	"oauth/internal/database"
	"oauth/internal/services"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot open file"})
		return
	}

	defer file.Close()

	service := services.FileService{
		Collection: database.DB.Collection("files"),
	}

	result, err := service.SaveFile(
		file,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		fileHeader.Size,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
