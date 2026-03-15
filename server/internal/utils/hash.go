package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func GenerateHash(file io.Reader) (string, error){
	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}