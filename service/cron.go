package service

import (
	"log"
	"time"
)

func StartCronService() {
	log.Println("Cron Service Started")
	for {
		select {
		case <-time.Tick(time.Hour):
			go RecordAllStatus()
		}
	}
}

func RecordAllStatus() {
	UserMapLock.RLock()
	defer UserMapLock.RUnlock()
	for _, user := range ValidUserMap {
		go RecordRawStatus(user)
		go GetRankByUserID(user.UserID)
	}
}
