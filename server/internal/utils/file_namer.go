package utils

import "github.com/google/uuid"

func GenerateFileName(ext string) string {
	return uuid.New().String() + ext
}