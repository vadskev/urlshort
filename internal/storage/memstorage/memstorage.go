package memstorage

import (
	"sync"

	"github.com/vadskev/urlshort/internal/entity"
)

type MemStorage struct {
	data sync.Map
}

func New() *MemStorage {
	return &MemStorage{data: sync.Map{}}
}

func (store *MemStorage) Add(link entity.Links) (entity.Links, error) {
	if _, ok := store.data.Load(link.Slug); ok {
		return entity.Links{}, entity.ErrSlugExists
	}
	store.data.Store(link.Slug, link)
	return link, nil
}

func (store *MemStorage) Get(key string) (entity.Links, error) {
	value, ok := store.data.Load(key)
	if !ok {
		return entity.Links{}, entity.ErrNotFound
	}
	link := value.(entity.Links)
	return link, nil
}
