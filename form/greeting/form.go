package greeting

import (
	tele "gopkg.in/telebot.v4"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/session"
)

const flowID = "greet"

const (
	StateIdle     session.State = "idle"
	StateName     session.State = "waiting_name"
	StateEmail    session.State = "waiting_email"
	StateStatus   session.State = "waiting_status"
	StateComplete session.State = "complete"
)

type Form struct {
	sessionStore *session.Store
	sender       configurator.MessageSender
}

func NewForm(sender configurator.MessageSender) *Form {
	return &Form{
		sessionStore: session.NewStore(),
		sender:       sender,
	}
}

func (f Form) HandleStep(c tele.Context) error {
	userID := c.Message().Sender.ID // проверять что это юзер, а не бот
	s := f.sessionStore.Get(userID)
	if s == nil {
		return nil
	}
	if s.FlowID != flowID {
		return nil
	}

	text := c.Message().Text
	switch s.State {
	case StateIdle:
		s.State = StateName
		if err := f.sender.Send(c, "form.start.greeting"); err != nil {
			return err
		}
		return f.sender.Send(c, "form.start.name")
	case StateName:
		s.Data["name"] = text
		s.State = StateEmail
		return f.sender.Send(c, "form.start.email")
	case StateEmail:
		s.Data["email"] = text
		s.State = StateStatus
		return f.sender.Send(c, "form.start.status")
	case StateStatus:
		s.Data["status"] = text
		s.State = StateComplete
		// отправить форму админу, сохранить в бд
		f.sessionStore.Delete(userID)
		return f.sender.Send(c, "form.start.complete")
	default:
		return nil
	}
}

func (f Form) CheckIsFormCompleted(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		userID := c.Message().Sender.ID
		s := f.sessionStore.Get(userID)
		if s != nil {
			//принудительно перенаправить на заполнение формы
			return f.HandleStep(c)
		}
		//todo посмотреть в бд - заполнял ли уже этот пользователь форму. Пока тут будет заглушка
		userExists := true
		if !userExists {
			f.sessionStore.Create(userID, StateIdle, flowID)
		}
		return next(c)
	}
}
