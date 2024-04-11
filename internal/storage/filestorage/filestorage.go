package filestorage

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

type FileStore struct {
	filePath string
	storage  sync.Map
}

func New(fPath string) *FileStore {
	return &FileStore{
		filePath: fPath,
		storage:  sync.Map{},
	}
}

func (fs *FileStore) Add(link entity.Links) (entity.Links, error) {
	if _, ok := fs.storage.Load(link.Slug); ok {
		return entity.Links{}, entity.ErrSlugExists
	}

	fs.storage.Store(link.Slug, link)

	err := fs.save()
	if err != nil {
		return entity.Links{}, err
	}
	return link, nil
}

func (fs *FileStore) save() error {
	var byteFile []byte

	fs.storage.Range(func(k, v interface{}) bool {
		link := v.(entity.Links)

		data, err := link.MarshalJSON()
		if err != nil {
			return false
		}

		data = append(data, '\n')
		byteFile = append(byteFile, data...)

		return true
	})

	fileName := filepath.FromSlash(fs.filePath)
	directory, _ := filepath.Split(fileName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(fileName, byteFile, 0666)
}

func (fs *FileStore) Load(ms *memstorage.MemStorage) error {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err
	}
	splitData := bytes.Split(data, []byte("\n"))

	for _, item := range splitData {
		link := entity.Links{}
		err = link.UnmarshalJSON(item)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
		}
		if _, ok := fs.storage.Load(link.Slug); ok {
			return nil
		}
		fs.storage.Store(link.Slug, link)

		_, err = ms.Add(link)
		if err != nil {
			return err
		}
	}

	return nil
}
