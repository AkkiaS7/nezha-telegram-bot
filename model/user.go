package model

import "errors"

type User struct {
	ID          int64 `gorm:"primary_key"`
	UserID      int64 `gorm:"unique_index"` //Telegram ID
	UserName    string
	FirstName   string
	LastName    string
	Valid       bool
	ConnectType int // 0: websocket, 1: API
	URL         string
	Token       string
}

func (u *User) Save() error {
	return DB.Save(u).Error
}

func GetUserByID(id int64) (*User, error) {
	var user User
	if err := DB.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetURLByID(id int64) (string, error) {
	var user User
	if err := DB.Where("user_id = ?", id).First(&user).Error; err != nil {
		return "", err
	}
	if !user.Valid {
		return "", errors.New("用户已被禁用")
	}
	return user.URL, nil
}

func GetAllValidUser() ([]*User, error) {
	var users []*User
	if err := DB.Where("valid = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllInvalidUser() ([]*User, error) {
	var users []*User
	if err := DB.Where("valid = ?", false).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserWithMissInfo() ([]*User, error) {
	var users []*User
	if err := DB.Where("user_name is NULL OR first_name is NULL").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) Suspend() {
	u.Valid = false
	u.Save()
	DeleteAllStatusByUserID(u.UserID)
	DeleteRankByUserID(u.UserID)
}
