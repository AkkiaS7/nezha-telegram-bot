package service

import (
	"errors"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"log"
)

type UserMgr struct {
	UserID    int64
	UserName  string
	FirstName string
	LastName  string
	URL       string
	Token     string
}

// SetURL 注册用户
func (um *UserMgr) SetURL() error {
	user, err := model.GetUserByID(um.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		return errors.New("更换监控地址将导致所有历史记录被移除\n如果您确定要更换 请使用 /override domain 需要包含http/https标识")
	}
	_, err = GetBriefByWebsocket(um.URL)
	if err != nil {
		return err
	}

	newUser := model.User{
		UserID:    um.UserID,
		UserName:  um.UserName,
		FirstName: um.FirstName,
		LastName:  um.LastName,
		Valid:     true,
		URL:       um.URL,
		Token:     um.Token,
	}
	err = newUser.Save()
	if err != nil {
		return err
	}
	UserMapLock.Lock()
	ValidUserMap[um.UserID] = &newUser
	UserMapLock.Unlock()
	return nil
}

func (um *UserMgr) Override() error {
	user, err := model.GetUserByID(um.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("请先使用 /reg 命令进行注册")
		}
		return err
	}
	_, err = GetBriefByWebsocket(um.URL)
	if err != nil {
		return err
	}

	// 修改用户
	user.URL = um.URL
	user.Token = um.Token
	user.Valid = true
	err = user.Save()
	if err != nil {
		return err
	}
	UserMapLock.Lock()
	delete(InvalidUserMap, user.UserID)
	ValidUserMap[user.UserID] = user
	UserMapLock.Unlock()
	// 删除原有的数据
	err = model.DeleteAllStatusByUserID(um.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func SuspendUser(user *model.User) {
	UserMapLock.Lock()
	defer UserMapLock.Unlock()
	delete(ValidUserMap, user.UserID)
	InvalidUserMap[user.UserID] = user
	user.Suspend()
}

func UpdateUserInfo() {
	users, err := model.GetUserWithMissInfo()
	if err != nil {
		log.Println(err)
		return
	}
	UserMapLock.Lock()
	defer UserMapLock.Unlock()
	for _, user := range users {
		chat, err := bot.ChatByID(user.UserID)
		if err != nil {
			if err == tele.ErrChatNotFound {
				SuspendUser(user)
				continue
			}
			continue
		}
		user.UserName = chat.Username
		user.FirstName = chat.FirstName
		user.LastName = chat.LastName
		if err := user.Save(); err != nil {
			log.Println(err)
			continue
		}
		if _, ok := ValidUserMap[user.UserID]; ok {
			ValidUserMap[user.UserID] = user
		} else if _, ok := InvalidUserMap[user.UserID]; ok {
			InvalidUserMap[user.UserID] = user
		}
	}
}
