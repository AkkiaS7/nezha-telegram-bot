package service

import (
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"github.com/patrickmn/go-cache"
	tele "gopkg.in/telebot.v3"
	"log"
	"sync"
	"time"
)

var (
	C              *cache.Cache
	UserMapLock    sync.RWMutex
	ValidUserMap   map[int64]*model.User
	InvalidUserMap map[int64]*model.User
	bot            *tele.Bot
)

func Init(b *tele.Bot) {
	bot = b
	C = cache.New(10*time.Minute, 30*time.Minute)
	ValidUserMap = make(map[int64]*model.User)
	InvalidUserMap = make(map[int64]*model.User)

	ValidUserMapInit()
	InvalidUserMapInit()
	rankListInit()
	UpdateUserInfo()
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
	UserMapLock.Lock()
	for _, user := range users {
		ValidUserMap[user.UserID] = user
	}
	UserMapLock.Unlock()
}

func InvalidUserMapInit() {
	users, err := model.GetAllInvalidUser()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	UserMapLock.Lock()
	for _, user := range users {
		InvalidUserMap[user.UserID] = user
	}
	UserMapLock.Unlock()
}
