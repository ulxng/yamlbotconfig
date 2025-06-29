package messages

import (
	"fmt"
	"ulxng/blueprintbot/lib/config"

	"gopkg.in/yaml.v3"
)

var _ config.Loader[Message] = (*Loader)(nil)

type Loader struct {
	data map[string]Message
}

func NewLoader(source string) (*Loader, error) {
	l := &Loader{data: make(map[string]Message)}

	if err := l.Load(source); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return l, nil
}

func (l *Loader) Load(path string) error {
	return config.LoadYamlFiles(path, l.ParseData)
}

func (l *Loader) GetByKey(key string) Message {
	return l.data[key]
}

func (l *Loader) All() map[string]Message {
	return l.data
}

func (l *Loader) ParseData(data []byte, path string) error {
	var parsed map[string]Message
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return fmt.Errorf("yaml.Unmarshal %s: %w", path, err)
	}

	for k, v := range parsed {
		if _, exists := l.data[k]; exists {
			return fmt.Errorf("duplicate key: %s in file %s", k, path)
		}
		l.data[k] = v
	}

	return nil
}
