package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/email"
	"ulxng/yamlbotconf/flow"
	"ulxng/yamlbotconf/state"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	bot    *tele.Bot
	store  *state.Store
	sender configurator.MessageSender
	mailer *email.Mailer

	flowRegistry *flow.Registry
}

type options struct {
	BotToken string `long:"token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`

	SmtpConfig email.SmtpConfig `group:"smtp" namespace:"smtp" env-namespace:"SMTP" required:"true"`
}

func main() {
	var opts options
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		os.Exit(1)
	}

	log.Println("bot started")

	app := App{}
	if err := app.run(opts); err != nil {
		log.Printf("run: %s", err)
	}

	log.Println("bot stopped")
}

func (a *App) run(opts options) error {
	pref := tele.Settings{
		Token:  opts.BotToken,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("tele.NewBot: %w", err)
	}
	a.bot = bot

	loader := configurator.NewLoader("responses")
	a.sender = configurator.NewConfigurableSenderAdapter(loader)
	a.mailer = email.NewMailer(opts.SmtpConfig)
	a.store = state.NewStore()

	flowLoader := flow.NewLoader("flow")
	a.flowRegistry = flow.NewRegistry(flowLoader)

	a.handleButtons()
	a.handleSendMailCommand()
	a.handleFlows()
	a.handleStart()

	a.bot.Start()
	return nil
}

func (a *App) FindFSM() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Message().Sender.ID
			session := a.store.Get(userID)
			var fsm *flow.FSM
			if session != nil {
				fsm = a.flowRegistry.FindUserActiveFlow(session)
			} else {
				fsm = a.flowRegistry.FindFlowToStart(c) // todo rename
				if fsm != nil {
					session = fsm.Start(userID)
					a.store.Save(userID, session)
				}
			}
			if fsm == nil {
				return nil
			}

			c.Set("fsm", fsm)
			c.Set("session", session)
			//передать управление на этот flow
			return next(c)
		}
	}
}
