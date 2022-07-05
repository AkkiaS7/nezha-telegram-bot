package service

import (
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"log"
	"time"
)

func RecordRawStatus(user *model.User) {
	msg, err := GetWebsocketMsg(user.URL)
	if err != nil {
		// TODO: 对于多次无法正常获取状态的用户需要被逐出
		return
	}
	recordTime := time.UnixMilli(msg.Now)
	for _, server := range msg.Servers {
		status := model.Status{}
		status.RecordType = 0
		status.UserID = user.ID
		status.ServerID = server.ID
		if recordTime.Sub(server.LastActive) > 30*time.Second {
			status.IsOnline = false
			status.RecordTime = recordTime

		} else {
			status.IsOnline = true
			status.RecordTime = server.LastActive

		}
		status.NetOutTransfer = server.State.NetOutTransfer
		status.NetInTransfer = server.State.NetInTransfer
		status.NetOutSpeed = server.State.NetOutSpeed
		status.NetInSpeed = server.State.NetInSpeed
		status.CPUUsed = server.State.CPU
		status.MemUsed = server.State.MemUsed
		status.DiskUsed = server.State.DiskUsed
		status.SwapUsed = server.State.SwapUsed
		status.MemTotal = server.Host.MemTotal
		status.SwapTotal = server.Host.SwapTotal
		status.DiskTotal = server.Host.DiskTotal
		status.Load1 = server.State.Load1
		status.Load5 = server.State.Load5
		status.Load15 = server.State.Load15
		status.TCPConnCount = server.State.TCPConnCount
		status.UDPConnCount = server.State.UDPConnCount
		status.ProcessCount = server.State.ProcessCount
		err := status.Save()
		if err != nil {
			log.Println(err)
		}
	}

}
