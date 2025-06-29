package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Loader[T any] interface {
	Load(path string) error
	GetByKey(key string) T
	All() map[string]T
}

type ParseFunc func(data []byte, path string) error

func LoadYamlFiles(path string, parseFunc ParseFunc) error {
	return filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkDir: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") {
			return nil
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("os.ReadFile %s: %w", filePath, err)
		}

		return parseFunc(data, filePath)
	})
}
