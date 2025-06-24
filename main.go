package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ulxng/yamlbotconf/configurator"
	"ulxng/yamlbotconf/email"
	"ulxng/yamlbotconf/form/greeting"

	"github.com/jessevdk/go-flags"
	tele "gopkg.in/telebot.v4"
)

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
	mailer := email.NewMailer(opts.SmtpConfig)

	bot.Handle(tele.OnCallback, func(c tele.Context) error {
		key := c.Callback().Data
		if key == "" {
			return nil
		}
		return sender.Edit(c, key)
	})
	bot.Handle("/send", func(c tele.Context) error {
		m, err := c.Bot().Send(c.Recipient(), "Отправляю сообщение на email...")
		if err != nil {
			return fmt.Errorf("bot.Send: %w", err)
		}
		go func() {
			if err := mailer.Send("Demo message", "Test"); err != nil {
				log.Printf("Failed to send email: %v", err)
				_, err := c.Bot().Edit(m, "Failed to send email")
				if err != nil {
					log.Printf("Failed to send message: %v", err)
					return
				}
			}
			_, err := c.Bot().Edit(m, "Заявка успешно отправлена")
			if err != nil {
				return
			}
		}()
		return nil
	})
	bot.Handle("/start", func(c tele.Context) error {
		return sender.Send(c, "main.menu")
	}, greetingForm.CheckIsFormCompleted)
	bot.Handle(tele.OnText, greetingForm.HandleStep, greetingForm.CheckIsFormCompleted)

	bot.Start()
	return nil
}
