package config

import (
	"fmt"
	"ulxng/blueprintbot/lib/config"
	"ulxng/blueprintbot/lib/messages"
)

type LoaderWithNav struct {
	// dictionary answer.text => message.code
	replyButtonsCodes map[string]string
	config.Loader[messages.Message]
}

func NewLoaderWithNav(source string) (*LoaderWithNav, error) {
	l, err := messages.NewLoader(source)
	if err != nil {
		return nil, fmt.Errorf("config.NewLoaderWithNav: %w", err)
	}
	loader := &LoaderWithNav{
		replyButtonsCodes: make(map[string]string),
		Loader:            l,
	}
	loader.buildTextToNextMap()
	return loader, nil
}

func (l *LoaderWithNav) buildTextToNextMap() {
	for _, msg := range l.All() {
		for _, btn := range msg.Answers {
			if btn.Code != "" {
				l.replyButtonsCodes[btn.Text] = btn.Code
			}
		}
	}
}

func (l *LoaderWithNav) GetNextByText(text string) (messages.Message, error) {
	if val, ok := l.replyButtonsCodes[text]; ok {
		return l.GetByKey(val), nil
	} else {
		return messages.Message{}, messages.ErrMessageNotFound
	}
}
