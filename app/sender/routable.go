package sender

import (
	tele "gopkg.in/telebot.v4"
	"ulxng/blueprintbot/app/resolver"
)

type SimpleRoutableSender struct {
	resolver resolver.RoutableResolver
	*DefaultSender
}

func NewSimpleRoutableSender(resolver resolver.RoutableResolver) *SimpleRoutableSender {
	return &SimpleRoutableSender{
		DefaultSender: NewDefaultSender(resolver),
		resolver:      resolver,
	}
}

func (s *SimpleRoutableSender) Route(c tele.Context, message string) error {
	m, markup, err := s.resolver.FindNextByText(message)
	if err != nil {
		return err
	}
	return c.Send(m, markup)
}
