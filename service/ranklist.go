package service

import (
	"errors"
	"fmt"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/utils"
	"github.com/AkkiaS7/nezha-telegram-bot/utils/config"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

const (
	ErrRankOverflow    = "无当前排名数据"
	ErrNoMoreRank      = "后续无排名数据"
	ErrUnknownRankType = "未知排名类型"
)

var (
	rankLock            = sync.RWMutex{}
	maxRankList         int
	ServerCountRankList []*model.RankList
	OnlineCountRankList []*model.RankList
	MemTotalRankList    []*model.RankList
	MemUsedRankList     []*model.RankList
	DiskTotalRankList   []*model.RankList
	DiskUsedRankList    []*model.RankList
	Load15RankList      []*model.RankList
)

func rankListInit() {
	rankLock.Lock()
	defer rankLock.Unlock()
	maxRankList = config.Conf.MaxRank

	ServerCountRankList = model.GetServerCountTop(maxRankList)
	OnlineCountRankList = model.GetOnlineCountTop(maxRankList)
	MemTotalRankList = model.GetMemTotalTop(maxRankList)
	MemUsedRankList = model.GetMemUsedTop(maxRankList)
	DiskTotalRankList = model.GetDiskTotalTop(maxRankList)
	DiskUsedRankList = model.GetDiskUsedTop(maxRankList)
	Load15RankList = model.GetLoad15Top(maxRankList)
}

func GetRankByUserID(userID int64) (string, error) {
	var tmp *WebsocketMsg
	var err error
	UserMapLock.RLock()
	if user, ok := ValidUserMap[userID]; ok {
		tmp, err = GetWebsocketMsg(user.URL)
		if err != nil {
			UserMapLock.RUnlock()
			return "", err
		}
	} else if _, ok := InvalidUserMap[userID]; ok {
		UserMapLock.RUnlock()
		return "", errors.New("无法查询被禁用的账户，请私聊机器人重新设置地址")
	} else {
		UserMapLock.RUnlock()
		return "", gorm.ErrRecordNotFound
	}
	UserMapLock.RUnlock()
	rankList := model.RankList{}
	rankList.UserID = userID
	for _, v := range tmp.Servers {
		if time.UnixMilli(tmp.Now).Before(v.LastActive.Add(time.Second * 15)) {
			rankList.OnlineCount++
		}
		rankList.MemTotal += v.Host.MemTotal
		rankList.MemUsedTotal += v.State.MemUsed
		rankList.DiskTotal += v.Host.DiskTotal
		rankList.DiskUsedTotal += v.State.DiskUsed
		rankList.Load15Total += v.State.Load15
	}
	rankList.ServerCount = len(tmp.Servers)
	serverCountRank := "未上榜"
	onlineCountRank := "未上榜"
	memTotalRank := "未上榜"
	memUsedRank := "未上榜"
	diskTotalRank := "未上榜"
	diskUsedRank := "未上榜"
	load15Rank := "未上榜"

	rankLock.RLock()
	for i, v := range ServerCountRankList {
		if v.ServerCount < rankList.ServerCount || v.UserID == userID {
			serverCountRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(ServerCountRankList)) + "]"
			break
		}
	}
	if serverCountRank == "未上榜" && len(ServerCountRankList) < maxRankList {
		serverCountRank = "排名[" + strconv.Itoa(len(ServerCountRankList)+1) + "/" + strconv.Itoa(len(ServerCountRankList)+1) + "]"
	}

	for i, v := range OnlineCountRankList {
		if v.OnlineCount < rankList.OnlineCount || v.UserID == userID {
			onlineCountRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(OnlineCountRankList)) + "]"
			break
		}
	}
	if onlineCountRank == "未上榜" && len(OnlineCountRankList) < maxRankList {
		onlineCountRank = "排名[" + strconv.Itoa(len(OnlineCountRankList)+1) + "/" + strconv.Itoa(len(OnlineCountRankList)+1) + "]"
	}

	for i, v := range MemTotalRankList {
		if v.MemTotal <= rankList.MemTotal || v.UserID == userID {
			memTotalRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "]"
			break
		}
	}
	if memTotalRank == "未上榜" && len(MemTotalRankList) < maxRankList {
		memTotalRank = "排名[" + strconv.Itoa(len(MemTotalRankList)+1) + "/" + strconv.Itoa(len(MemTotalRankList)+1) + "]"
	}
	for i, v := range MemUsedRankList {
		if v.MemUsedTotal <= rankList.MemUsedTotal || v.UserID == userID {
			memUsedRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "]"
			break
		}
	}
	if memUsedRank == "未上榜" && len(MemUsedRankList) < maxRankList {
		memUsedRank = "排名[" + strconv.Itoa(len(MemUsedRankList)+1) + "/" + strconv.Itoa(len(MemTotalRankList)+1) + "]"
	}
	for i, v := range DiskTotalRankList {
		if v.DiskTotal <= rankList.DiskTotal || v.UserID == userID {
			diskTotalRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "]"
			break
		}
	}
	if diskTotalRank == "未上榜" && len(DiskTotalRankList) < maxRankList {
		diskTotalRank = "排名[" + strconv.Itoa(len(DiskTotalRankList)+1) + "/" + strconv.Itoa(len(MemTotalRankList)+1) + "]"
	}
	for i, v := range DiskUsedRankList {
		if v.DiskUsedTotal <= rankList.DiskUsedTotal || v.UserID == userID {
			diskUsedRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "]"
			break
		}
	}
	if diskUsedRank == "未上榜" && len(DiskUsedRankList) < maxRankList {
		diskUsedRank = "排名[" + strconv.Itoa(len(DiskUsedRankList)+1) + "/" + strconv.Itoa(len(MemTotalRankList)+1) + "]"
	}
	for i, v := range Load15RankList {
		if v.Load15Total <= rankList.Load15Total || v.UserID == userID {
			load15Rank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "]"
			break
		}
	}
	if load15Rank == "未上榜" && len(Load15RankList) < maxRankList {
		load15Rank = "排名[" + strconv.Itoa(len(Load15RankList)+1) + "/" + strconv.Itoa(len(MemTotalRankList)+1) + "]"
	}
	rankLock.RUnlock()
	if memUsedRank != "未上榜" || memTotalRank != "未上榜" || diskUsedRank != "未上榜" || diskTotalRank != "未上榜" || load15Rank != "未上榜" {
		rankList.Save()
		rankListInit()
	}
	str := fmt.Sprint(
		"服务器数量: ", rankList.ServerCount, " ", serverCountRank, "\n",
		"在线服务器数量: ", rankList.OnlineCount, " ", onlineCountRank, "\n",
		"内存总量: ", utils.AutoUnitConvert(rankList.MemTotal), " ", memTotalRank, "\n",
		"内存使用量: ", utils.AutoUnitConvert(rankList.MemUsedTotal), " ", memUsedRank, "\n",
		"磁盘总量: ", utils.AutoUnitConvert(rankList.DiskTotal), " ", diskTotalRank, "\n",
		"磁盘使用量: ", utils.AutoUnitConvert(rankList.DiskUsedTotal), " ", diskUsedRank, "\n",
		"总负载: ", fmt.Sprintf("%.2f", rankList.Load15Total), " ", load15Rank, "\n",
	)
	return str, nil
}

func GetRankList(rankType string, rank int) (string, error) {
	switch rankType {
	case "all":
		return GetAllRankList(rank)
	default:
		return GetSepRankList(rankType, rank)
	}
}

func GetAllRankList(rank int) (string, error) {
	if rank > len(ServerCountRankList) {
		return "", errors.New(ErrRankOverflow)
	}
	rankLock.RLock()
	defer rankLock.RUnlock()
	msg := "正在显示各指标排名第" + strconv.Itoa(rank) + "的数据:\n"
	msg += "服务器数量: " + strconv.Itoa(ServerCountRankList[rank-1].ServerCount) + "用户: " + GetATAbleStringByUserID(ServerCountRankList[rank-1].UserID) + "\n"
	msg += "在线服务器数量: " + strconv.Itoa(OnlineCountRankList[rank-1].OnlineCount) + "用户: " + GetATAbleStringByUserID(OnlineCountRankList[rank-1].UserID) + "\n"
	msg += "内存总量: " + utils.ParseForMarkdown(utils.AutoUnitConvert(MemTotalRankList[rank-1].MemTotal)) + "用户: " + GetATAbleStringByUserID(MemTotalRankList[rank-1].UserID) + "\n"
	msg += "内存使用量: " + utils.ParseForMarkdown(utils.AutoUnitConvert(MemUsedRankList[rank-1].MemUsedTotal)) + "用户: " + GetATAbleStringByUserID(MemUsedRankList[rank-1].UserID) + "\n"
	msg += "磁盘总量: " + utils.ParseForMarkdown(utils.AutoUnitConvert(DiskTotalRankList[rank-1].DiskTotal)) + "用户: " + GetATAbleStringByUserID(DiskTotalRankList[rank-1].UserID) + "\n"
	msg += "磁盘使用量: " + utils.ParseForMarkdown(utils.AutoUnitConvert(DiskUsedRankList[rank-1].DiskUsedTotal)) + "用户: " + GetATAbleStringByUserID(DiskUsedRankList[rank-1].UserID) + "\n"
	msg += "总负载: " + utils.ParseForMarkdown(fmt.Sprintf("%.2f", Load15RankList[rank-1].Load15Total)) + "用户: " + GetATAbleStringByUserID(Load15RankList[rank-1].UserID) + "\n"
	if rank+1 > len(ServerCountRankList) {
		return msg, errors.New(ErrNoMoreRank)
	}
	return msg, nil
}

func GetSepRankList(rankType string, rankPage int) (string, error) {
	rankLock.RLock()
	defer rankLock.RUnlock()
	if (rankPage-1)*10 > len(ServerCountRankList) {
		return "", errors.New(ErrRankOverflow)
	}
	msg := ""
	switch rankType {
	case "serverCount":
		msg = "正在显示服务器总数排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(ServerCountRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(ServerCountRankList)) + "\\] " + strconv.Itoa(ServerCountRankList[i].ServerCount) + " 用户: " + GetATAbleStringByUserID(ServerCountRankList[i].UserID) + "\n"
		}

	case "onlineCount":
		msg = "正在显示在线服务器数量排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(OnlineCountRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(OnlineCountRankList)) + "\\] " + strconv.Itoa(OnlineCountRankList[i].OnlineCount) + " 用户: " + GetATAbleStringByUserID(OnlineCountRankList[i].UserID) + "\n"
		}
	case "ramTotal":
		msg = "正在显示内存总量排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(MemTotalRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRankList)) + "\\] " + utils.ParseForMarkdown(utils.AutoUnitConvert(MemTotalRankList[i].MemTotal)) + " 用户: " + GetATAbleStringByUserID(MemTotalRankList[i].UserID) + "\n"
		}
	case "ramUsed":
		msg = "正在显示内存使用量排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(MemUsedRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemUsedRankList)) + "\\] " + utils.ParseForMarkdown(utils.AutoUnitConvert(MemUsedRankList[i].MemUsedTotal)) + " 用户: " + GetATAbleStringByUserID(MemUsedRankList[i].UserID) + "\n"
		}
	case "diskTotal":
		msg = "正在显示磁盘总量排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(DiskTotalRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(DiskTotalRankList)) + "\\] " + utils.ParseForMarkdown(utils.AutoUnitConvert(DiskTotalRankList[i].DiskTotal)) + " 用户: " + GetATAbleStringByUserID(DiskTotalRankList[i].UserID) + "\n"
		}
	case "diskUsed":
		msg = "正在显示磁盘使用量排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(DiskUsedRankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(DiskUsedRankList)) + "\\] " + utils.ParseForMarkdown(utils.AutoUnitConvert(DiskUsedRankList[i].DiskUsedTotal)) + " 用户: " + GetATAbleStringByUserID(DiskUsedRankList[i].UserID) + "\n"
		}
	case "load15":
		msg = "正在显示15分钟负载排行榜第" + strconv.Itoa(rankPage) + "页的数据:\n"
		for i := (rankPage - 1) * 10; i < rankPage*10; i++ {
			if i >= len(Load15RankList) {
				break
			}
			msg += "排名\\[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(Load15RankList)) + "\\] " + fmt.Sprintf("%.2f", Load15RankList[i].Load15Total) + " 用户: " + GetATAbleStringByUserID(Load15RankList[i].UserID) + "\n"
		}
	default:
		return "", errors.New(ErrUnknownRankType)
	}
	if rankPage*10 >= len(ServerCountRankList) {
		return msg, errors.New(ErrNoMoreRank)
	}
	return msg, nil
}

func GetATAbleStringByUserID(userID int64) string {
	UserMapLock.RLock()
	user, ok := ValidUserMap[userID]
	if !ok {
		UserMapLock.RUnlock()
		return "无效用户"
		// 还需要重新校准排名列表
	}
	UserMapLock.RUnlock()
	name := user.FirstName
	if name == "" {
		name = user.LastName
	}
	if user.UserName != "" {
		return "[" + utils.ParseForMarkdown(name) + "](t.me/" + user.UserName + ")"
	}
	return "`" + utils.ParseForMarkdown(name) + "`"
}
