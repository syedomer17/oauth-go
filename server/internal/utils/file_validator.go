package utils

import (
	"errors"
	"mime/multipart"
)

var allowedTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

func ValidateFile(fileHeader *multipart.FileHeader, maxSize int64) error {

	if fileHeader.Size > maxSize {
		return errors.New("File too large")
	}

	if !allowedTypes[fileHeader.Header.Get("Content-Type")] {
		return errors.New("Invalid file type")
	}

	return nil
}