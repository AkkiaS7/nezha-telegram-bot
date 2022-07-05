package model

type User struct {
	ID          int64 `gorm:"primary_key"`
	UserID      int64 `gorm:"unique_index"` //Telegram ID
	Valid       bool
	ConnectType int // 0: websocket, 1: API
	URL         string
	Token       string
}

func (u *User) Register() error {
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
	return user.URL, nil
}

func GetAllValidUser() ([]*User, error) {
	var users []*User
	if err := DB.Where("valid = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
