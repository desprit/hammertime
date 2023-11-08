package config

import (
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `env:"APP_ENVIRONMENT"`
	RootPath       string
	TG_CHAT_ID     string `env:"TG_CHAT_ID"`
	TG_BOT_TOKEN   string `env:"TG_BOT_TOKEN"`
	HAMMER_TOKEN_D string `env:"HAMMER_TOKEN_D"`
	HAMMER_TOKEN_M string `env:"HAMMER_TOKEN_M"`
	WEB_USER       string `env:"WEB_USER"`
	WEB_PASS       string `env:"WEB_PASS"`
}

func (c *Config) DbUri() string {
	if c.Env == "testing" {
		return "file:memdb1?mode=memory&_foreign_keys=on&_journal_mode=WAL&_timeout=5000"
	}
	if c.Env == "production" {
		return "/app/sqlite_data/db.sqlite3?_foreign_keys=on&_journal_mode=WAL&_timeout=5000"
	}
	return "file:sqlite_data/db.sqlite3?_foreign_keys=on&_journal_mode=WAL&_timeout=5000"
}

var instance Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&instance)
		if err != nil {
			log.Fatalln(err)
		}
		_, b, _, _ := runtime.Caller(0)
		instance.RootPath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
	})
	return &instance
}
