package messages

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Loader struct {
	data map[string]Message
}

func NewLoader(source string) *Loader {
	data, err := os.ReadFile(source)
	if err != nil {
		panic(err)
	}

	var raw map[string]Message
	err = yaml.Unmarshal(data, &raw)
	if err != nil {
		panic(err)
	}
	return &Loader{data: raw}
}

func (l *Loader) GetByKey(key string) Message {
	return l.data[key]
}
