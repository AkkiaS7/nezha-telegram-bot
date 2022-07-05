package model

import "time"

type Status struct {
	ID             int64 `gorm:"primary_key,unique_index,auto_increment"`
	UserID         int64
	ServerID       int
	RecordTime     time.Time
	RecordType     int // 0-Raw, 1-Avg5min, 2-Avg15min, 3-Avg60Min 4-Avg24Hour
	IsOnline       bool
	MemTotal       int64
	SwapTotal      int
	DiskTotal      int64
	CPUUsed        float64
	MemUsed        int64
	SwapUsed       int
	DiskUsed       int64
	NetInTransfer  int64
	NetOutTransfer int64
	NetInSpeed     int
	NetOutSpeed    int
	Load1          float64
	Load5          float64
	Load15         float64
	TCPConnCount   int
	UDPConnCount   int
	ProcessCount   int
}

func DeleteAllStatusByUserID(userID int64) error {
	return DB.Delete(&Status{}, "user_id = ?", userID).Error
}

func (s *Status) Save() error {
	return DB.Save(s).Error
}
