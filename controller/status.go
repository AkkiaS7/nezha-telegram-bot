package controller

import (
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	"github.com/AkkiaS7/nezha-telegram-bot/utils/config"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"strconv"
)

func statusInit() {
	bot.Handle("/b", getBrief, middleware.AutoDelete)
	AddCommand("/b", "获取服务器状态简述")
}

func getBrief(c tele.Context) error {
	userID := c.Sender().ID
	text, err := service.GetBriefByUserID(userID)
	msg := ""
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg = "获取失败\n请先私聊bot并使用 /seturl 命令进行注册"
		} else {
			msg = "获取失败\n" + err.Error()
		}
	} else {
		msg = text
	}
	reply, err := bot.Reply(c.Message(), msg)
	if err != nil {
		return err
	}
	if config.Conf.AutoDelete.Enable {
		replyMsg := &model.Message{}
		replyMsg.StoredMessage = tele.StoredMessage{
			MessageID: strconv.Itoa(reply.ID),
			ChatID:    c.Chat().ID,
		}
		replyMsg.Save()
		middleware.DelayDelete(replyMsg)
	}
	return nil
}
