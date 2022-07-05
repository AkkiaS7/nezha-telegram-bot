package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/AkkiaS7/nezha-telegram-bot/utils"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

func GetBriefByUserID(userID int64) (string, error) {
	url, err := model.GetURLByID(userID)
	if err != nil {
		return "", err
	}
	return GetBriefByWebsocket(url)
}

func GetBriefByWebsocket(url string) (string, error) {
	tmp, err := GetWebsocketMsg(url)
	if err != nil {
		return "", err
	}
	type brief struct {
		Online         int
		Offline        int
		RamOver80      int
		CPUOver80      int
		DiskOver80     int
		NetInTransfer  int64
		NetOutTransfer int64
		NetInSpeed     int
		NetOutSpeed    int
	}

	b := &brief{}
	for _, v := range tmp.Servers {
		b.NetOutSpeed += v.State.NetOutSpeed
		b.NetInSpeed += v.State.NetInSpeed
		b.NetOutTransfer += v.State.NetOutTransfer
		b.NetInTransfer += v.State.NetInTransfer
		if v.LastActive.Unix() > time.Now().Unix()-30 {
			b.Online++
		} else {
			b.Offline++
		}
		if float64(v.State.MemUsed)/float64(v.Host.MemTotal) > 0.8 {
			b.RamOver80++
		}
		if v.State.CPU > 80 {
			b.CPUOver80++
		}
		if float64(v.State.DiskUsed)/float64(v.Host.DiskTotal) > 0.8 {
			b.DiskOver80++
		}
	}
	str := fmt.Sprint("在线: ", b.Online, ", 离线: ", b.Offline, "\n",
		"内存使用率超过80%: ", b.RamOver80, "\n",
		"CPU使用率超过80%: ", b.CPUOver80, "\n",
		"磁盘使用率超过80%: ", b.DiskOver80, "\n",
		"下行流量: ", utils.AutoUnitConvert(b.NetInTransfer), ", 上行流量: ", utils.AutoUnitConvert(b.NetOutTransfer), "\n",
		"下行带宽: ", utils.AutoBandwidthConvert(int64(b.NetInSpeed)), "， 上行带宽: ", utils.AutoBandwidthConvert(int64(b.NetOutSpeed)), "\n")

	return str, nil
}

func GetWebsocketMsg(url string) (*WebsocketMsg, error) {
	dialer := websocket.Dialer{}
	url = strings.Replace(url, "http", "ws", 1)
	if url[len(url)-1] == '/' {
		url += "ws"
	} else {
		url += "/ws"
	}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msgT, msg, err := conn.ReadMessage()
	if msgT != websocket.TextMessage {
		return nil, errors.New("msg type error")
	}
	tmp := &WebsocketMsg{}
	err = json.Unmarshal(msg, tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil

}

type WebsocketMsg struct {
	Now     int64     `json:"now"`
	Servers []Servers `json:"servers"`
}
type Host struct {
	Platform        string   `json:"Platform"`
	PlatformVersion string   `json:"PlatformVersion"`
	CPU             []string `json:"CPU"`
	MemTotal        int64    `json:"MemTotal"`
	DiskTotal       int64    `json:"DiskTotal"`
	SwapTotal       int      `json:"SwapTotal"`
	Arch            string   `json:"Arch"`
	Virtualization  string   `json:"Virtualization"`
	BootTime        int      `json:"BootTime"`
	CountryCode     string   `json:"CountryCode"`
	Version         string   `json:"Version"`
}
type State struct {
	CPU            float64 `json:"CPU"`
	MemUsed        int64   `json:"MemUsed"`
	SwapUsed       int     `json:"SwapUsed"`
	DiskUsed       int64   `json:"DiskUsed"`
	NetInTransfer  int64   `json:"NetInTransfer"`
	NetOutTransfer int64   `json:"NetOutTransfer"`
	NetInSpeed     int     `json:"NetInSpeed"`
	NetOutSpeed    int     `json:"NetOutSpeed"`
	Uptime         int     `json:"Uptime"`
	Load1          float64 `json:"Load1"`
	Load5          float64 `json:"Load5"`
	Load15         float64 `json:"Load15"`
	TCPConnCount   int     `json:"TcpConnCount"`
	UDPConnCount   int     `json:"UdpConnCount"`
	ProcessCount   int     `json:"ProcessCount"`
}
type Servers struct {
	ID           int         `json:"ID"`
	CreatedAt    time.Time   `json:"CreatedAt"`
	UpdatedAt    time.Time   `json:"UpdatedAt"`
	DeletedAt    interface{} `json:"DeletedAt"`
	Name         string      `json:"Name"`
	Tag          string      `json:"Tag"`
	DisplayIndex int         `json:"DisplayIndex"`
	Host         Host        `json:"Host"`
	State        State       `json:"State"`
	LastActive   time.Time   `json:"LastActive"`
}
