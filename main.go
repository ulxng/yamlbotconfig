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
	"ulxng/yamlbotconf/storage"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	bot    *tele.Bot
	store  *state.Store
	sender configurator.MessageSender
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

	loader := configurator.NewLoader("responses")
	a.sender = configurator.NewConfigurableSenderAdapter(loader)
	a.mailer = email.NewMailer(opts.SmtpConfig)
	a.store = state.NewStore()
	a.userRepository = storage.NewUserMemoryStorage()

	flowLoader := flow.NewLoader("flow")
	a.flowRegistry = flow.NewRegistry(flowLoader)

	a.registerRoutes()
	a.registerFlows()

	a.bot.Start()
	return nil
}
