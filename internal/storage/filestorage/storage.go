package filestorage

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/vadskev/urlshort/internal/lib/logger/zp"
	"github.com/vadskev/urlshort/internal/storage"
	"go.uber.org/zap"
)

type FileStore struct {
	filePath string
	log      *zap.Logger
}

var _ storage.Storage = (*FileStore)(nil)

func NewFileStorage(filePath string, logger *zap.Logger) (*FileStore, error) {
	fileName := filepath.FromSlash(filePath)
	directory, _ := filepath.Split(fileName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			logger.Info("Error to create file", zp.Err(err))
			return nil, err
		}
		return nil, err
	}
	logger.Info("Create file", zap.String("patch", filepath.Dir(directory)))
	return &FileStore{
		filePath: filePath,
		log:      logger,
	}, nil
}

func (fs *FileStore) Get(ctx context.Context, ms storage.Storage) error {
	if _, err := os.Stat(fs.filePath); errors.Is(err, os.ErrNotExist) {
		fs.log.Info("Error to open file", zp.Err(err))
		return err
	}

	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		fs.log.Info("Error to read file", zp.Err(err))
		return err
	}
	splitData := bytes.Split(data, []byte("\n"))

	for _, item := range splitData {
		link := storage.URLData{}
		if json.Valid(item) {
			err = json.Unmarshal(item, &link)
			if err != nil {
				fs.log.Info("Error to Unmarshal file", zp.Err(err))
				return err
			}
			err = ms.SaveURL(ctx, link)
			if err != nil {
				fs.log.Info("Error to save memory in file", zp.Err(err))
				return err
			}
		}

	}
	return nil
}

func (fs *FileStore) GetURL(ctx context.Context, alias string) (storage.URLData, error) {
	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0774)
	defer func() {
		err = file.Close()
		if err != nil {
			fs.log.Info("Error to close file", zp.Err(err))
		}
	}()

	if err != nil {
		return storage.URLData{}, err
	}
	scanner := bufio.NewScanner(file)
	var shortUrls []storage.URLData

	for scanner.Scan() {
		var url storage.URLData
		err = json.Unmarshal(scanner.Bytes(), &url)
		if err != nil {
			return storage.URLData{}, err
		}
		shortUrls = append(shortUrls, url)
	}

	for _, v := range shortUrls {
		if v.Alias == alias {
			return v, nil
		}
	}

	return storage.URLData{}, nil
}

func (fs *FileStore) SaveURL(ctx context.Context, data storage.URLData) error {
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0774)
	defer func() {
		err = file.Close()
		if err != nil {
			fs.log.Info("Error to close file", zp.Err(err))
		}
	}()
	if err != nil {
		fs.log.Info("Error to open file", zp.Err(err))
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		fs.log.Info("Error to encode file", zp.Err(err))
		return err
	}
	return nil
}

func (fs *FileStore) SaveBatchURL(ctx context.Context, data []storage.URLData) error {
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0774)
	defer func() {
		err = file.Close()
		if err != nil {
			fs.log.Info("Error to close file", zp.Err(err))
		}
	}()
	if err != nil {
		fs.log.Info("Error to open file", zp.Err(err))
		return err
	}
	encoder := json.NewEncoder(file)
	for _, v := range data {
		err = encoder.Encode(v)
		if err != nil {
			fs.log.Info("Error to encode file", zp.Err(err))
			return err
		}
	}
	return nil
}

func (fs *FileStore) Ping(ctx context.Context) error {
	//TODO implement me
	return nil
}

func (fs *FileStore) GetURLbyURL(ctx context.Context, url string) (storage.URLData, bool) {
	if _, err := os.Stat(fs.filePath); errors.Is(err, os.ErrNotExist) {
		fs.log.Info("Error to open file", zp.Err(err))
		return storage.URLData{}, false
	}
	fdata, err := os.ReadFile(fs.filePath)
	if err != nil {
		fs.log.Info("Error to read file", zp.Err(err))
		return storage.URLData{}, false
	}
	splitData := bytes.Split(fdata, []byte("\n"))
	for _, item := range splitData {
		link := storage.URLData{}
		if json.Valid(item) {
			err = json.Unmarshal(item, &link)
			if err != nil {
				fs.log.Info("Error to Unmarshal file", zp.Err(err))
				return storage.URLData{}, false
			}
			if link.URL == url {
				return link, true
			}
		}
	}
	return storage.URLData{}, false
}
