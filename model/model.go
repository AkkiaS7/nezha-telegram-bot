package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/db.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&User{}, &Status{}, &Message{}, &RankList{})
	if err != nil {
		panic(err)
	}
}
