package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ulxng/blueprintbot/app/email"
	"ulxng/blueprintbot/app/sender"
	"ulxng/blueprintbot/app/storage"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/messages"
	"ulxng/blueprintbot/lib/state"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	bot    *tele.Bot
	store  *state.Store
	sender sender.MessageSender
	mailer *email.Mailer

	flowRegistry *flow.Registry

	userRepository storage.UserRepository
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

	loader, err := messages.NewLoader("data/config/responses")
	if err != nil {
		return fmt.Errorf("messages.NewLoader: %w", err)
	}
	a.sender = sender.NewConfigurableSenderAdapter(loader)
	a.mailer = email.NewMailer(opts.SmtpConfig)
	a.store = state.NewStore()
	a.userRepository = storage.NewUserMemoryStorage()

	flowLoader, err := flow.NewLoader("data/config/flows")
	if err != nil {
		return fmt.Errorf("flow.NewLoader: %w", err)
	}
	a.flowRegistry = flow.NewRegistry(flowLoader)

	a.registerRoutes()
	a.registerFlows()

	a.bot.Start()
	return nil
}
