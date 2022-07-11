package config

import (
	"gopkg.in/ini.v1"
	"time"
)

var (
	Conf *Config
)

type Config struct {
	Bot
	AutoDelete
	RankList
}

type Bot struct {
	Token string
}

type AutoDelete struct {
	Enable bool
	Time   time.Duration
}

type RankList struct {
	Enable  bool
	MaxRank int
}

func Init() {
	cfg, err := ini.Load("data/config.ini")
	if err != nil {
		panic(err)
	}
	Conf = new(Config)
	if err := cfg.MapTo(Conf); err != nil {
		panic(err)
	}
	Conf.AutoDelete.Time *= time.Second
}
