package main

import (
	"TeleEcho/api/database"
	"TeleEcho/configs"
	"TeleEcho/router"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	err := configs.ParseConfig()
	if err != nil {
		logrus.Printf("Can not read config")
	}
	err = database.ConnectDB()
	if err != nil {
		fmt.Printf("err:%s", err)
	}
	e := router.New()
	err = e.Start(configs.Config.Address + ":" + configs.Config.Port)
	if err != nil {
		fmt.Printf("err:%s", err)
	}
}
