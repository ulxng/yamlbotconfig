package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ulxng/blueprintbot/app/config"
	"ulxng/blueprintbot/app/fsm"
	"ulxng/blueprintbot/app/resolver"
	"ulxng/blueprintbot/app/sender"
	"ulxng/blueprintbot/lib/flow"
	"ulxng/blueprintbot/lib/state"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	bot    *tele.Bot
	sender sender.RoutableSender

	fsmExecutor  *fsm.TelebotExecutor
	flowRegistry *flow.Registry
	loader       config.NavigationLoader
}

const (
	flowsConfigPath    = "data/config/flows.yaml" // can be directory too
	messagesConfigPath = "data/config/messages.yaml"
)

type options struct {
	BotToken string `long:"token" env:"BOT_TOKEN" required:"true" description:"telegram bot token"`
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

	loader, err := config.NewNavigableLoader(messagesConfigPath)
	if err != nil {
		return fmt.Errorf("config.NewLoader: %w", err)
	}
	a.loader = loader
	lt := resolver.NewNavigableResolver(loader)
	a.sender = sender.NewSimpleRoutableSender(lt)

	flowLoader, err := flow.NewLoader(flowsConfigPath)
	if err != nil {
		return fmt.Errorf("flow.NewLoader: %w", err)
	}
	a.flowRegistry = flow.NewRegistry(flowLoader)
	a.fsmExecutor = fsm.NewExecutor(state.NewMemoryStore(), a.sender, a.flowRegistry, a.bot)

	a.registerRoutes()
	a.registerFlows()

	a.bot.Start()
	return nil
}
