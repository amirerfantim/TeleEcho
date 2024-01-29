package database

import (
	"TeleEcho/configs"
	"TeleEcho/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var NotFoundUser = errors.New("user not found")
var NotFoundContact = errors.New("contact not found")
var UnauthorizedChatAccess = errors.New("unauthorized access to the chat")
var IncorrectPassword = errors.New("incorrect password")
var NotFoundGroup = errors.New("group not found")

func ConnectDB() error {
	var err error
	dsn := "host=" + configs.Config.DatabaseAddress + " user=" + configs.Config.DatabaseUser + " password=" + configs.Config.DatabasePassword +
		" dbname=" + configs.Config.DatabaseName + " port=" + configs.Config.DatabasePort
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Printf("Failed to connect to database:%s", err)
		return err
	}
	logrus.Printf("Connected to database successfully.\n")
	err = DB.AutoMigrate(&model.User{})
	err = DB.AutoMigrate(&model.Contact{})
	err = DB.AutoMigrate(&model.Group{})
	err = DB.AutoMigrate(&model.UserGroup{})
	err = DB.AutoMigrate(&model.DirectChat{})
	err = DB.AutoMigrate(&model.GroupChat{})
	err = DB.AutoMigrate(&model.Message{})
	if err != nil {
		logrus.Printf("Failed to migrate database:%s", err)
		return err
	}
	return nil

}
