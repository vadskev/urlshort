package random

import "math/rand"

// GenerateRandomString Алгоритм сокращения URL
func GenerateRandomString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := make([]byte, 7)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}
