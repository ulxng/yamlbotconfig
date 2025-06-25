package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	Flows map[string]Flow
}

func NewLoader(source string) *Loader {
	l := Loader{Flows: make(map[string]Flow)}
	err := l.loadYamlFiles(source)
	if err != nil {
		return nil
	}
	return &l
}

func (l *Loader) GetByKey(key string) Flow {
	return l.Flows[key]
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

		//todo тут немного другой механизм парсинга нужен
		//верхний уровень - flow. Под ним
		var parsed map[string]Flow
		if err := yaml.Unmarshal(data, &parsed); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		for k, v := range parsed {
			if _, exists := t.Flows[k]; exists {
				return fmt.Errorf("duplicate key: %s in file %s", k, path)
			}
			t.Flows[k] = v
		}

		return nil
	})
}
