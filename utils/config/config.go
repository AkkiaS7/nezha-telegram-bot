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
}

type Bot struct {
	Token string
}

type AutoDelete struct {
	Enable bool
	Time   time.Duration
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
