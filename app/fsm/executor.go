package fsm

import (
	"ulxng/blueprintbot/app/sender"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/fsm"
	"ulxng/blueprintbot/lib/state"

	tele "gopkg.in/telebot.v4"
)

type TelebotExecutor struct {
	fsm.BaseExecutor
	bapi *botAPI
}

func NewExecutor(store state.Store, sender sender.Sender, registry *flow.Registry, bot *tele.Bot) *TelebotExecutor {
	api := &botAPI{
		bot:    bot,
		sender: sender,
	}
	return &TelebotExecutor{
		BaseExecutor: fsm.NewBaseExecutor(store, registry, api),
		bapi:         api,
	}
}

func (e *TelebotExecutor) Middleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			e.bapi.SetContext(c)
			userID := c.Message().Chat.ID
			if err := e.RunFSM(userID); err != nil {
				return next(c)
			}
			return nil
		}
	}
}
