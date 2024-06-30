package memstorage

import (
	"context"
	"errors"
	"sync"

	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type MemStorage struct {
	store sync.Map
	log   *zap.Logger
}

var _ storage.Storage = (*MemStorage)(nil)

func NewMemStorage(log *zap.Logger) *MemStorage {
	return &MemStorage{
		store: sync.Map{},
		log:   log,
	}
}

func (s *MemStorage) SaveURL(ctx context.Context, data storage.URLData) error {
	if _, ok := s.store.Load(data.Alias); ok {
		return errors.New("url exists")
	}
	s.store.Store(data.Alias, data)
	s.log.Info("added to storage")
	return nil
}

func (s *MemStorage) SaveBatchURL(ctx context.Context, data []storage.URLData) error {
	for _, v := range data {
		if _, ok := s.store.Load(v.Alias); ok {
			return errors.New("url exists")
		}
		s.store.Store(v.Alias, v)
	}
	s.log.Info("added to mem storage")
	return nil
}

func (s *MemStorage) GetURL(ctx context.Context, alias string) (storage.URLData, error) {
	value, ok := s.store.Load(alias)
	if !ok {
		return storage.URLData{}, errors.New("url not found")
	}
	s.log.Info("return url storage")
	return value.(storage.URLData), nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	//TODO implement me
	return nil
}
