package persistence

import (
	"io/ioutil"
	"os"
)

type fileRepository struct {}

func (f fileRepository) Save(file []byte, path string) error {
	err := ioutil.WriteFile(path, file, 0644)
	return err
}

func (f fileRepository) Delete(path string) error {
	err := os.Remove(path)
	return err
}

