package resolver

import (
	"fmt"
	"ulxng/blueprintbot/app/config"

	tele "gopkg.in/telebot.v4"
)

type BaseNavigableResolver struct {
	*BaseResolver
	loader config.NavigationLoader
}

func NewNavigableResolver(loader config.NavigationLoader) *BaseNavigableResolver {
	return &BaseNavigableResolver{
		loader:       loader,
		BaseResolver: NewBaseResolver(loader),
	}
}

func (lt *BaseNavigableResolver) FindNextByText(response string) (interface{}, *tele.ReplyMarkup, error) {
	m, err := lt.loader.GetNextByText(response)
	if err != nil {
		return nil, nil, fmt.Errorf("getNextByText: %w", err)
	}
	return lt.Convert(m)
}
