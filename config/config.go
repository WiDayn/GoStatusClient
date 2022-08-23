package config

import (
	"fmt"
	"github.com/go-ini/ini"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
)

type config struct {
	ClientId    string
	SecretKey   string `comment:"需要与服务器[SecretKey]字段一致"`
	Debug       bool
	DisplayName string
	Server      Server
	CT          PingAddr
	CU          PingAddr
	CM          PingAddr
}

type Server struct {
	IP   string
	Port int
}

type PingAddr struct {
	Addr string
}

var Config config

func Read() {
	if !isExist("./config.ini") {
		if _, err := os.Create("./config.ini"); err != nil {
			log.Println("[PANIC] Can not create config.ini")
			log.Fatalf(err.Error())
		}
	}

	if cfg, err := ini.Load("./config.ini"); err != nil {
		log.Println("[PANIC] Can not load config.ini")
		log.Fatalf(err.Error())
	} else {
		Config = config{
			ClientId:    uuid.NewV4().String(),
			DisplayName: "Server",
			CT:          PingAddr{Addr: "14.116.225.60"},
			CU:          PingAddr{Addr: "123.125.96.156"},
			CM:          PingAddr{Addr: "39.134.69.205"},
		}

		if err = ini.MapTo(&Config, "./config.ini"); err != nil {
			log.Println("[PANIC] config.ini error")
			log.Fatalf(err.Error())
		}

		if err := ini.ReflectFrom(cfg, &Config); err != nil {
			log.Println("[PANIC] config.ini error")
			log.Fatalf(err.Error())
		}

		if err := cfg.SaveTo("./config.ini"); err != nil {
			log.Println("[PANIC] Can not sava config.ini")
			log.Fatalf(err.Error())
		}
	}
	log.Println("[SUCCESS] Success load config.ini")
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
