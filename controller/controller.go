package controller

import (
	tele "gopkg.in/telebot.v3"
)

var (
	bot     *tele.Bot
	cmdList = make([]tele.Command, 0)
)

func Init(telebot *tele.Bot) {
	bot = telebot
	userInit()
	statusInit()
	rankListInit()
}

func Serve() {
	bot.SetCommands(cmdList)
}

func AddCommand(text string, desc string) {
	cmdList = append(cmdList, tele.Command{
		Text:        text,
		Description: desc,
	})
}
