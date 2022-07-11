package controller

import (
	"errors"
	"github.com/AkkiaS7/nezha-telegram-bot/middleware"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/service"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func rankListInit() {
	bot.Handle("/rank", getRank, middleware.AutoDelete)
	AddCommand("/rank", "获取排名")
	bot.Handle("/ranklist", getRankListMenu, middleware.AutoDelete)
	AddCommand("/ranklist", "/ranklist 获取排名列表")

	bot.Handle("\fbtnRank", btnRank, middleware.AutoDelete)
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

func getRankListMenuBtn() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btnAllRank := menu.Data("各指标总榜", "btnRank", "all|1")
	btnServerCountRank := menu.Data("服务器总数", "btnRank", "serverCount|1")
	btnServerOnlineRank := menu.Data("在线服务器数", "btnRank", "onlineCount|1")
	btnRAMTotalRank := menu.Data("内存总量", "btnRank", "ramTotal|1")
	btnRAMUsedRank := menu.Data("内存使用量", "btnRank", "ramUsed|1")
	btnDiskTotalRank := menu.Data("磁盘总量", "btnRank", "diskTotal|1")
	btnDiskUsedRank := menu.Data("磁盘使用量", "btnRank", "diskUsed|1")
	btnLoadRank := menu.Data("平均负载", "btnRank", "load15|1")
	menu.Inline(
		menu.Row(btnAllRank),
		menu.Row(btnServerCountRank, btnServerOnlineRank),
		menu.Row(btnRAMTotalRank, btnRAMUsedRank),
		menu.Row(btnDiskTotalRank, btnDiskUsedRank),
		menu.Row(btnLoadRank),
	)
	return menu
}

func getRankListMenu(c tele.Context) error {
	menu := getRankListMenuBtn()
	reply, err := bot.Reply(c.Message(), "请选择排名类型", menu)
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

func btnRank(c tele.Context) error {
	strs := strings.Split(c.Data(), "|")
	if len(strs) != 2 {
		return errors.New("错误的输入参数")
	}
	rankType := strs[0]
	rank, err := strconv.Atoi(strs[1])
	if err != nil {
		log.Println(err)
		return err
	}
	if rankType == "backToMenu" {
		menu := getRankListMenuBtn()
		err := c.Edit("请选择排名类型", menu)
		return err
	}

	msg, msgErr := service.GetRankList(rankType, rank)
	if msgErr != nil && msgErr.Error() == service.ErrUnknownRankType {
		log.Println(msgErr)
		return msgErr
	}

	menu := &tele.ReplyMarkup{}
	btnBack := menu.Data("回到菜单", "btnRank", "backToMenu|1")
	btnNext := menu.Data("下一页", "btnRank", rankType+"|"+strconv.Itoa(rank+1))
	btnPrev := menu.Data("上一页", "btnRank", rankType+"|"+strconv.Itoa(rank-1))
	if msgErr == nil {
		if rank != 1 {
			menu.Inline(menu.Row(btnNext, btnPrev), menu.Row(btnBack))
		} else {
			menu.Inline(menu.Row(btnNext), menu.Row(btnBack))
		}
	} else if msgErr.Error() == service.ErrNoMoreRank {
		if rank != 1 {
			menu.Inline(menu.Row(btnPrev), menu.Row(btnBack))
		} else {
			menu.Inline(menu.Row(btnBack))
		}
	}
	err = c.Edit(msg, menu, tele.ModeMarkdownV2, tele.NoPreview)
	msgModel := &model.Message{
		StoredMessage: tele.StoredMessage{
			MessageID: strconv.Itoa(c.Message().ID),
			ChatID:    c.Chat().ID,
		},
	}
	middleware.DelayDelete(msgModel)
	return err
}
