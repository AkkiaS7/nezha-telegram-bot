package service

import (
	"fmt"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/utils"
	"strconv"
	"sync"
	"time"
)

const maxRankList = 20

var (
	rankLock        = sync.RWMutex{}
	ServerCountRank = make([]*model.RankList, maxRankList)
	OnlineCountRank = make([]*model.RankList, maxRankList)
	MemTotalRank    = make([]*model.RankList, maxRankList)
	MemUsedRank     = make([]*model.RankList, maxRankList)
	DiskTotalRank   = make([]*model.RankList, maxRankList)
	DiskUsedRank    = make([]*model.RankList, maxRankList)
	Load15Rank      = make([]*model.RankList, maxRankList)
)

func rankListInit() {
	rankLock.Lock()
	defer rankLock.Unlock()
	ServerCountRank = model.GetServerCountTop(maxRankList)
	OnlineCountRank = model.GetOnlineCountTop(maxRankList)
	MemTotalRank = model.GetMemTotalTop(maxRankList)
	MemUsedRank = model.GetMemUsedTop(maxRankList)
	DiskTotalRank = model.GetDiskTotalTop(maxRankList)
	DiskUsedRank = model.GetDiskUsedTop(maxRankList)
	Load15Rank = model.GetLoad15Top(maxRankList)
}

func GetRankByUserID(userID int64) (string, error) {
	url, err := model.GetURLByID(userID)
	if err != nil {
		return "", err
	}
	tmp, err := GetWebsocketMsg(url)
	if err != nil {
		return "", err
	}
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

	for i, v := range ServerCountRank {
		if v.ServerCount < rankList.ServerCount || v.UserID == userID {
			serverCountRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(ServerCountRank)) + "]"
			break
		}
	}
	if serverCountRank == "未上榜" && len(ServerCountRank) < maxRankList {
		serverCountRank = "排名[" + strconv.Itoa(len(ServerCountRank)+1) + "/" + strconv.Itoa(len(ServerCountRank)+1) + "]"
	}

	for i, v := range OnlineCountRank {
		if v.OnlineCount < rankList.OnlineCount || v.UserID == userID {
			onlineCountRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(OnlineCountRank)) + "]"
			break
		}
	}
	if onlineCountRank == "未上榜" && len(OnlineCountRank) < maxRankList {
		onlineCountRank = "排名[" + strconv.Itoa(len(OnlineCountRank)+1) + "/" + strconv.Itoa(len(OnlineCountRank)+1) + "]"
	}

	for i, v := range MemTotalRank {
		if v.MemTotal <= rankList.MemTotal || v.UserID == userID {
			memTotalRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRank)) + "]"
			break
		}
	}
	if memTotalRank == "未上榜" && len(MemTotalRank) < maxRankList {
		memTotalRank = "排名[" + strconv.Itoa(len(MemTotalRank)+1) + "/" + strconv.Itoa(len(MemTotalRank)+1) + "]"
	}
	for i, v := range MemUsedRank {
		if v.MemUsedTotal <= rankList.MemUsedTotal || v.UserID == userID {
			memUsedRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRank)) + "]"
			break
		}
	}
	if memUsedRank == "未上榜" && len(MemUsedRank) < maxRankList {
		memUsedRank = "排名[" + strconv.Itoa(len(MemUsedRank)+1) + "/" + strconv.Itoa(len(MemTotalRank)+1) + "]"
	}
	for i, v := range DiskTotalRank {
		if v.DiskTotal <= rankList.DiskTotal || v.UserID == userID {
			diskTotalRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRank)) + "]"
			break
		}
	}
	if diskTotalRank == "未上榜" && len(DiskTotalRank) < maxRankList {
		diskTotalRank = "排名[" + strconv.Itoa(len(DiskTotalRank)+1) + "/" + strconv.Itoa(len(MemTotalRank)+1) + "]"
	}
	for i, v := range DiskUsedRank {
		if v.DiskUsedTotal <= rankList.DiskUsedTotal || v.UserID == userID {
			diskUsedRank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRank)) + "]"
			break
		}
	}
	if diskUsedRank == "未上榜" && len(DiskUsedRank) < maxRankList {
		diskUsedRank = "排名[" + strconv.Itoa(len(DiskUsedRank)+1) + "/" + strconv.Itoa(len(MemTotalRank)+1) + "]"
	}
	for i, v := range Load15Rank {
		if v.Load15Total <= rankList.Load15Total || v.UserID == userID {
			load15Rank = "排名[" + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(MemTotalRank)) + "]"
			break
		}
	}
	if load15Rank == "未上榜" && len(Load15Rank) < maxRankList {
		load15Rank = "排名[" + strconv.Itoa(len(Load15Rank)+1) + "/" + strconv.Itoa(len(MemTotalRank)+1) + "]"
	}
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
