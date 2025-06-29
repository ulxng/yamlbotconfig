package main

import (
	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	a.bot.Handle("/start", func(c tele.Context) error {
		return a.sender.Send(c, "main.menu")
	})
	a.bot.Handle(tele.OnText, a.handleError, func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if err := a.sender.NavigateFor(c, c.Text()); err != nil {
				return next(c)
			}
			return nil
		}
	})
}

func (a *App) registerFlows() {

}
