package model

import (
	tele "gopkg.in/telebot.v3"
	"time"
)

type Message struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	tele.StoredMessage
}

func (m *Message) Save() {
	DB.Save(m)
}

func (m *Message) Delete() {
	DB.Delete(m)
}

func GetAllMessage() []Message {
	var messages []Message
	DB.Find(&messages)
	return messages
}

func GetAllMessageBefore(t time.Time) []Message {
	var messages []Message
	DB.Where("updated_at < ?", t).Find(&messages)
	return messages
}
