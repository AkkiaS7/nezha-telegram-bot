package controller

import (
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func rankListInit() {
	bot.Handle("/rank", getRank, middleware.AutoDelete)
	AddCommand("/rank", "获取排名")
	bot.Handle("/ranklist", getRankList, middleware.AutoDelete)
	AddCommand("/ranklist", "/ranklist 1 获取排名列表")
}

func getRank(c tele.Context) error {
	userID := c.Sender().ID
	text, err := service.GetRankByUserID(userID)
	msg := ""
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg = "获取失败\n请先私聊bot并使用 /seturl 命令设置URL"
		} else {
			msg = "获取失败\n" + err.Error() + "\n请私聊bot并发送 /seturl 命令设置URL"
		}
	} else {
		msg = text
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

func getRankList(c tele.Context) error {
	text := c.Message().Payload
	rank, err := strconv.Atoi(text)
	log.Println(rank)
	msg := ""
	if err != nil {
		msg = "输入错误"
	} else {
		msg = service.GetRankList(rank)
	}
	opt := &tele.SendOptions{
		ParseMode: "MarkdownV2",
	}
	reply, err := bot.Reply(c.Message(), msg, opt, tele.NoPreview)
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
