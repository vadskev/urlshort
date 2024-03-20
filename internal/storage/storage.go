package storage

import (
	"errors"
)

type MemStorage struct {
	data map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func (s *MemStorage) AddURL(key string, value string) error {
	s.data[key] = value
	return nil
}

func (s *MemStorage) GetURL(key string) (string, error) {
	value, ok := s.data[key]
	if !ok {
		return "", errors.New("url not found: " + key)
	}
	return value, nil
}
