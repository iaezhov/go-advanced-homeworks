package file

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileStorage struct {
	fileName string
}

func NewFileStorage(name string) *FileStorage {
	return &FileStorage{fileName: name}
}

func NewJSONFileStorage(name string) (*FileStorage, error) {
	if filepath.Ext(name) != ".json" {
		return nil, fmt.Errorf("Файл должен иметь расширение .json")
	}
	return NewFileStorage(name), nil
}

func (fs *FileStorage) Read() ([]byte, error) {
	data, err := os.ReadFile(fs.fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fs *FileStorage) Write(content []byte) error {
	file, err := os.Create(fs.fileName)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}
