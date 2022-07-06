package controller

import (
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	tele "gopkg.in/telebot.v3"
	"strconv"
)

func userInit() {
	bot.Handle("/seturl", seturl)

	AddCommand("/seturl", "设置URL 需要包含http/https标识")
}

func seturl(c tele.Context) error {
	userID := c.Sender().ID
	args := c.Args()
	msg := ""
	if len(args) < 1 || len(args) > 2 {
		msg = "参数错误 请使用 /seturl <url> 需要包含http/https标识"
	} else {
		url := args[0]
		token := ""
		if len(args) == 2 {
			token = args[1]
		}

		u := service.UserMgr{
			ID:    userID,
			URL:   url,
			Token: token,
		}

		if err := u.SetURL(); err == nil {
			msg = "设置成功"
		} else {
			msg = "设置失败" + err.Error()
		}
	}
	reply, err := bot.Reply(c.Message(), msg)
	if err != nil {
		return err
	}
	replyMsg := &model.Message{}
	replyMsg.StoredMessage = tele.StoredMessage{
		MessageID: strconv.Itoa(reply.ID),
		ChatID:    c.Chat().ID,
	}
	replyMsg.Save()
	middleware.DelayDelete(replyMsg)
	return nil
}
