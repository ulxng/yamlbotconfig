package resolver

import (
	"errors"
	"ulxng/blueprintbot/lib/config"
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

type BaseResolver struct {
	loader config.Loader[messages.Message]
}

func NewBaseResolver(loader config.Loader[messages.Message]) *BaseResolver {
	return &BaseResolver{loader: loader}
}

func (lt *BaseResolver) Get(key string) (interface{}, *tele.ReplyMarkup, error) {
	msgLayout := lt.loader.GetByKey(key)
	if msgLayout.Text == "" {
		return nil, nil, messages.ErrMessageNotFound
	}
	return lt.Convert(msgLayout)
}

func (lt *BaseResolver) Convert(msgLayout messages.Message) (interface{}, *tele.ReplyMarkup, error) {
	var response interface{}
	if msgLayout.Image != "" {
		response = &tele.Photo{File: tele.FromDisk(msgLayout.Image), Caption: msgLayout.Text}
	}
	if msgLayout.File != "" {
		response = &tele.Document{File: tele.FromDisk(msgLayout.File), Caption: msgLayout.Text}
	}
	if msgLayout.Text != "" {
		response = msgLayout.Text
	}

	if msgLayout.Buttons != nil && msgLayout.Answers != nil {
		return nil, nil, errors.New("cannot use reply keyboard and inline keyboard together")
	}

	markup := &tele.ReplyMarkup{}
	if msgLayout.Buttons != nil {
		for _, button := range msgLayout.Buttons {
			btn := tele.InlineButton{
				Text: button.Text,
				Data: button.Code,
				URL:  button.Link,
			}
			markup.InlineKeyboard = append(markup.InlineKeyboard, []tele.InlineButton{btn})
		}
	} else {
		markup.OneTimeKeyboard = true
		for _, button := range msgLayout.Answers {
			btn := tele.ReplyButton{
				Text:    button.Text,
				Contact: button.Contact,
			}
			markup.ReplyKeyboard = append(markup.ReplyKeyboard, []tele.ReplyButton{btn})
		}
	}

	return response, markup, nil
}
