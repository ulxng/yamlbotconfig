package main

import (
	"fmt"
	"log"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/layout"
)

func main() {
	log.Println("bot started")

	if err := run(); err != nil {
		log.Printf("run: %s", err)
	}

	log.Println("bot stopped")
}

func run() error {
	lt, err := layout.NewDefault("layout/bot.yaml", "ru")
	if err != nil {
		log.Fatal("layout.New:", err)
	}

	bot, err := tele.NewBot(lt.Settings())
	if err != nil {
		return fmt.Errorf("tele.NewBot: %w", err)
	}

	bot.Handle("/start", func(c tele.Context) error {
		return c.Send(lt.Text("start"), lt.Markup("menu"))
	})
	bot.Handle(lt.Callback("back"), func(c tele.Context) error {
		return c.Send(lt.Text("start"), lt.Markup("menu"))
	})

	bot.Handle("/help", func(c tele.Context) error {
		return c.Send(lt.Text("start"), lt.Markup("back"))
	})
	bot.Handle(lt.Callback("help"), func(c tele.Context) error {
		return c.Send(lt.Text("help"), lt.Markup("back"))
	})

	bot.Handle(lt.Callback("settings"), func(c tele.Context) error {
		return c.Send(lt.Text("menu"), lt.Markup("settings"))
	})

	bot.Handle(lt.Callback("account"), func(c tele.Context) error {
		return c.Send(lt.Text("account"), lt.Markup("back"))
	})
	bot.Handle(lt.Callback("notification"), func(c tele.Context) error {
		return c.Send(lt.Text("notification"), lt.Markup("notification"))
	})

	bot.Start()
	return nil
}
