package services

import (
	"context"
	"io"
	"oauth/internal/models"
	"oauth/internal/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileService struct {
	Collection *mongo.Collection
}

func (s *FileService) SaveFile(file io.Reader, filename string, mime string, size int64) (*models.File, error) {

	hash, err := utils.GenerateHash(file)
	if err != nil {
		return nil, err
	}

	fileName := uuid.New().String() + filepath.Ext(filename)

	path := "uploads/" + fileName

	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	fileDoc := models.File{
		ID:        uuid.New().String(),
		Filename:  fileName,
		Path:      path,
		Size:      size,
		MimeType:  mime,
		Hash:      hash,
		CreatedAt: time.Now().Unix(),
	}

	_, err = s.Collection.InsertOne(context.Background(), fileDoc)
	if err != nil {
		return nil, err
	}

	return &fileDoc, nil
}
