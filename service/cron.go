package service

import (
	"log"
	"time"
)

func StartCronService() {
	log.Println("Cron Service Started")
	for range time.Tick(time.Hour) {
		go RecordAllStatus()
	}
}

func RecordAllStatus() {
	UserMapLock.RLock()
	defer UserMapLock.RUnlock()
	for _, user := range ValidUserMap {
		// 不再记录状态
		// go RecordRawStatus(user)
		go GetRankByUserID(user.UserID)
	}
}
