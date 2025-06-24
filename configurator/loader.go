package configurator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	data map[string]Message
}

func NewLoader(source string) *Loader {
	l := Loader{data: make(map[string]Message)}
	err := l.loadYamlFiles(source)
	if err != nil {
		return nil
	}
	return &l
}

func (l *Loader) GetByKey(key string) Message {
	return l.data[key]
}

func (t *Loader) loadYamlFiles(path string) error {
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		var parsed map[string]Message
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		for k, v := range parsed {
			if _, exists := t.data[k]; exists {
				return fmt.Errorf("duplicate key: %s in file %s", k, path)
			}
			t.data[k] = v
		}

		return nil
	})
}
