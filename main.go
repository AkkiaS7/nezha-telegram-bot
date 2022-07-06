package main

import (
	"github.com/AkkiaS7/nezha-telegram-bot/controller"
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	"github.com/AkkiaS7/nezha-telegram-bot/utils/config"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

var bot *tele.Bot

func init() {
	config.Init()
	var err error
	pref := tele.Settings{
		Token: config.Conf.Token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	log.Println("Authorized on account " + bot.Me.Username)

	controller.Init(bot)
	middleware.Init(bot)
	model.Init()
	service.Init()
}

func main() {
	controller.Serve()
	middleware.Serve()
	service.Serve()

	bot.Start()
}
