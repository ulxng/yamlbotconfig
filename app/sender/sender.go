package sender

import (
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

type MessageSender interface {
	Send(c tele.Context, key string) error
	Edit(c tele.Context, key string) error
	SendRaw(c tele.Context, message messages.Message) error
}

// ConfigurableSenderAdapter ищет конфигурацию сообщения по ключу и готовит его к отправке.
// Методы - обертки над с.Send(), c.Edit() и тд
type ConfigurableSenderAdapter struct {
	loader *messages.Loader
}

func NewConfigurableSenderAdapter(loader *messages.Loader) *ConfigurableSenderAdapter {
	return &ConfigurableSenderAdapter{loader: loader}
}

func (b *ConfigurableSenderAdapter) Send(c tele.Context, messageKey string) error {
	msg := b.loader.GetByKey(messageKey)
	if msg.Text == "" {
		return messages.ErrMessageNotFound
	}

	return b.SendRaw(c, msg)
}

func (b *ConfigurableSenderAdapter) SendRaw(c tele.Context, msg messages.Message) error {
	markup := &tele.ReplyMarkup{}
	for _, button := range msg.Buttons {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []tele.InlineButton{
			{Text: button.Text, Data: button.Code, URL: button.Link},
		})
	}
	if msg.Image != "" {
		photo := &tele.Photo{File: tele.FromDisk(msg.Image), Caption: msg.Text}
		return c.Send(photo, markup)
	}
	if msg.File != "" {
		file := &tele.Document{File: tele.FromDisk(msg.File), Caption: msg.Text}
		return c.Send(file, markup)
	}
	// только для fsm режима, где все исходящие отправляются методом send. В edit mode никаких reply кнопок
	if len(msg.Answers) == 0 {
		markup.RemoveKeyboard = true
	} else {
		markup.OneTimeKeyboard = true
		for _, button := range msg.Answers {
			markup.ReplyKeyboard = append(markup.ReplyKeyboard, []tele.ReplyButton{
				{Text: button.Text, Contact: button.Contact},
			})
		}
	}
	return c.Send(msg.Text, markup)
}

func (b *ConfigurableSenderAdapter) Edit(c tele.Context, messageKey string) error {
	msg := b.loader.GetByKey(messageKey)
	if msg.Text == "" {
		return messages.ErrMessageNotFound
	}

	markup := &tele.ReplyMarkup{}
	for _, button := range msg.Buttons {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []tele.InlineButton{
			{Text: button.Text, Data: button.Code, URL: button.Link},
		})
	}
	if msg.Image != "" {
		photo := &tele.Photo{File: tele.FromDisk(msg.Image), Caption: msg.Text}
		return c.Send(photo, markup)
	}
	if msg.File != "" {
		file := &tele.Document{File: tele.FromDisk(msg.File), Caption: msg.Text}
		return c.Send(file, markup)
	}
	if c.Message().Text == "" {
		// для предотвращения ошибки telegram: Bad Request: there is no text in the message to edit (400)
		// происходит, когда пытаешься сделать edit на сообщениях с кнопками и файлом/картинкой
		return c.Send(msg.Text, markup)
	}
	return c.Edit(msg.Text, markup)
}
