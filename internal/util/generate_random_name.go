package util

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
)

func GenerateRandomFilename(originalFilename string) string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return originalFilename
	}
	randomString := hex.EncodeToString(randomBytes)
	ext := filepath.Ext(originalFilename)
	return randomString + ext
}
