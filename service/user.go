package service

import (
	"errors"
	"github.com/AkkiaS7/nezha-telegram-bot/model"
	"gorm.io/gorm"
	"log"
)

type UserMgr struct {
	ID    int64
	URL   string
	Token string
}

// SetURL 注册用户
func (um *UserMgr) SetURL() error {
	user, err := model.GetUserByID(um.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user != nil {
		return errors.New("更换监控地址将导致所有历史记录被移除\n如果您确定要更换 请使用 /override domain token")
	}
	_, err = GetBriefByWebsocket(um.URL)
	if err != nil {
		return err
	}

	newUser := model.User{
		UserID: um.ID,
		Valid:  true,
		URL:    um.URL,
		Token:  um.Token,
	}
	return newUser.Register()
}

func (um *UserMgr) Override() error {
	user, err := model.GetUserByID(um.ID)
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
	err = user.Register()
	if err != nil {
		return err
	}

	// 删除原有的数据
	err = model.DeleteAllStatusByUserID(um.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
