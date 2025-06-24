package configurator

import (
	tele "gopkg.in/telebot.v4"
)

type MessageSender interface {
	Send(c tele.Context, key string) error
	Edit(c tele.Context, key string) error
}

// ConfigurableSenderAdapter ищет конфигурацию сообщения по ключу и готовит его к отправке.
// Методы - обертки над с.Send(), c.Edit() и тд
type ConfigurableSenderAdapter struct {
	loader *Loader
}

func NewConfigurableSenderAdapter(loader *Loader) *ConfigurableSenderAdapter {
	return &ConfigurableSenderAdapter{loader: loader}
}

func (b *ConfigurableSenderAdapter) Send(c tele.Context, messageKey string) error {
	msg := b.loader.GetByKey(messageKey)
	if msg.Text == "" {
		return nil
	}

	markup := &tele.ReplyMarkup{}
	for _, button := range msg.Buttons {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []tele.InlineButton{
			{Text: button.Text, Data: button.Code},
		})
	}
	// только для fsm режима, где все исходящие отправляются методом edit. В edit mode никаких reply кнопок
	if len(msg.Answers) == 0 {
		markup.RemoveKeyboard = true
	} else {
		markup.OneTimeKeyboard = true
		for _, button := range msg.Answers {
			markup.ReplyKeyboard = append(markup.ReplyKeyboard, []tele.ReplyButton{
				{Text: button},
			})
		}
	}
	return c.Send(msg.Text, markup)
}

func (b *ConfigurableSenderAdapter) Edit(c tele.Context, messageKey string) error {
	msg := b.loader.GetByKey(messageKey)
	if msg.Text == "" {
		return nil
	}

	markup := &tele.ReplyMarkup{}
	for _, button := range msg.Buttons {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []tele.InlineButton{
			{Text: button.Text, Data: button.Code},
		})
	}
	return c.Edit(msg.Text, markup)
}
