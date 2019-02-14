package app

import (
	"encoding/base64"
	"crypto/rand"
)

func generateSecureBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomStringURLSafe(n int) (string, error) {
	b, err := generateSecureBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
