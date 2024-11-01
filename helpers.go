package celeritas

import (
	"fmt"
	"os"
)

// CreateDirIfNotExist takes a path string
// and returns an error. It creates the directory if it does not exist.
func (c *Celeritas) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, mode)
	}
	return nil
}

// CreateFileIfNotExists creates a new empty file at the specified path if one doesn't exist.
// If the file already exists, it does nothing and returns nil.
// The created file is automatically closed after creation.
// It returns an error if the file creation fails.
func (c *Celeritas) CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				err = fmt.Errorf("error closing file: %w", err)
				// Handle or log the error appropriately
				fmt.Println(err)
			}
		}(file)
	}
	return nil
}
