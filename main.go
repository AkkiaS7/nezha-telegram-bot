package main

import (
	"flag"
	"github.com/AkkiaS7/nezha-telegram-bot/controller"
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

var bot *tele.Bot

func init() {
	var err error
	token := flag.String("token", "", "Telegram bot token")
	flag.Parse()
	if *token == "" {
		log.Println("必须指定token，请使用参数 --token=<token>")
	}
	pref := tele.Settings{
		Token: *token,
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
