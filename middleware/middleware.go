package middleware

import tele "gopkg.in/telebot.v3"

var bot *tele.Bot

func Init(telebot *tele.Bot) {
	bot = telebot
}

func Serve() {
	AutoDeleteInit()
}
