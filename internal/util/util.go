package util

import (
	"math/rand"
	"net/url"
)

// GenerateRandomString Алгоритм сокращения URL
func GenerateRandomString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := make([]byte, 7)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}

// IsURL проверяем валидность URL
func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
