package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/form/greeting"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
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

	if err := run(opts); err != nil {
		log.Printf("run: %s", err)
	}

	log.Println("bot stopped")
}

func run(opts options) error {
	pref := tele.Settings{
		Token:  opts.BotToken,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("tele.NewBot: %w", err)
	}

	loader := configurator.NewLoader("responses")
	sender := configurator.NewConfigurableSenderAdapter(loader)
	greetingForm := greeting.NewForm(sender)

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		key := c.Callback().Data
		if key == "" {
			return nil
		}
		return sender.Edit(c, key)
	})
	bot.Handle("/start", func(c tele.Context) error {
		return sender.Send(c, "main.menu")
	}, greetingForm.CheckIsFormCompleted)
	bot.Handle(tele.OnText, greetingForm.HandleStep, greetingForm.CheckIsFormCompleted)

	bot.Start()
	return nil
}
