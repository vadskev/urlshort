package storage

import "github.com/vadskev/urlshort/internal/entity"

// Storage
type Storage interface {
	Add(link entity.Links) (*entity.Links, error)
	Get(key string) (*entity.Links, error)
}
