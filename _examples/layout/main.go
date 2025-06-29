package main

import (
    "fmt"
    tele "gopkg.in/telebot.v4"
    "gopkg.in/telebot.v4/layout"
    "log"
)

func main() {
    log.Println("bot started")

    if err := run(); err != nil {
        log.Printf("run: %s", err)
    }

    log.Println("bot stopped")
}

func run() error {
    lt, err := layout.NewDefault("layout/bot.yaml", "en")
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
    bot.Handle(lt.Callback("help"), func(c tele.Context) error {
        return c.Send(lt.Text("menu"))
    })

    bot.Start()
    return nil
}
