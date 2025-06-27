package flow

import (
	"fmt"
	"ulxng/blueprintbot/lib/config"

	"gopkg.in/yaml.v3"
)

// для проверки на этапе компиляции
var _ config.Loader[Flow] = (*Loader)(nil)

type Loader struct {
	Flows map[string]Flow
}

func NewLoader(source string) (*Loader, error) {
	l := &Loader{Flows: make(map[string]Flow)}
	if err := l.Load(source); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return l, nil
}

func (l *Loader) Load(path string) error {
	return config.LoadYamlFiles(path, l.ParseData)
}

func (l *Loader) GetByKey(key string) Flow {
	return l.Flows[key]
}

func (l *Loader) ParseData(data []byte, path string) error {
	var parsed map[string]Flow
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	for k, v := range parsed {
		if _, exists := l.Flows[k]; exists {
			return fmt.Errorf("duplicate key: %s in file %s", k, path)
		}
		l.Flows[k] = v
	}
	return nil
}
