package main

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"TeleEcho/router"
	"github.com/sirupsen/logrus"
)

func main() {
	err := configs.ParseConfig()
	if err != nil {
		logrus.Printf("Can not read config")
	}
	err = database.ConnectDB()
	if err != nil {
		logrus.Printf("err:%s", err)
	}
	err = services.ConnectS3()
	if err != nil {
		logrus.Printf("err:%s", err)
	}
	e := router.New()
	err = e.Start(configs.Config.Address + ":" + configs.Config.Port)
	if err != nil {
		logrus.Printf("err:%s", err)
	}
}
