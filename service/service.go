package service

import (
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var (
	C            *cache.Cache
	ValidUserMap map[int64]*model.User
)

func Init() {
	C = cache.New(10*time.Minute, 30*time.Minute)
	ValidUserMap = make(map[int64]*model.User)

	ValidUserMapInit()
	rankListInit()
	log.Println(ValidUserMap)
}

func Serve() {
	go StartCronService()
}

func ValidUserMapInit() {
	users, err := model.GetAllValidUser()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	for _, user := range users {
		ValidUserMap[user.UserID] = user
	}
}
