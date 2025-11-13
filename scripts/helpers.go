package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func ValidateToken(token string, expiry time.Time) error {
	if token == "" {
		return ErrInvalidToken
	}

	if time.Now().After(expiry) {
		return ErrExpiredToken
	}

	return nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}