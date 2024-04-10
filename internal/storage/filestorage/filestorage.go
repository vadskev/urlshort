package filestorage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/vadskev/urlshort/internal/entity"
	"github.com/vadskev/urlshort/internal/storage/memstorage"
)

type FileStore struct {
	file         *os.File
	inMemoryData *memstorage.MemStorage
	filePath     string
}

func New(filePath string) (*FileStore, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	store := &FileStore{
		inMemoryData: memstorage.New(),
		file:         file,
		filePath:     filePath,
	}
	return store, nil
}

func (fs *FileStore) SaveToFileStorage(link *entity.Links) error {
	file, err := os.OpenFile(fs.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	line, err := link.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = writer.Write(append(line, '\n'))
	if err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	return nil
}

func (fs *FileStore) ReadFileStorage(memstore *memstorage.MemStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var link entity.Links
		if err = json.Unmarshal(line, &link); err != nil {
			return err
		}

		_, err := memstore.Add(link)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}