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
var UnauthorizedChatAccess = errors.New("unauthorized access to the chat")
var IncorrectPassword = errors.New("incorrect password")

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

func CreateUser(username, firstname, lastname, phone, password, profilePicture, bio string) error {
	userData := &model.User{Username: username, Firstname: firstname, Lastname: lastname, Phone: phone,
		Password: password, ProfilePicture: profilePicture, Bio: bio}
	if err := DB.Create(userData).Error; err != nil {
		logrus.Printf("Error creating user:%s\n", err)
		return err
	}
	logrus.Printf("User created\n")
	return nil

}
func IsUsernameDuplicate(username string) bool {
	var count int64
	DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return false
	}
	return true
}
func IsPhoneDuplicate(phoneNumber string) bool {
	var count int64
	DB.Model(&model.User{}).Where("phone = ?", phoneNumber).Count(&count)
	if count > 0 {
		return false
	}
	return true
}
func CheckPassword(username string, hashedPassword string) (*model.User, error) {
	var user model.User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, NotFoundUser
		}
		return nil, result.Error
	}
	if hashedPassword != user.Password {
		return nil, IncorrectPassword
	}
	return &user, nil
}
func DeleteUserByUserID(userID uint) error {
	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		logrus.Printf("Error while querying basket: %s", err)
		return err
	}
	if err := DB.Delete(&user, userID).Error; err != nil {
		logrus.Printf("Can not remove basket with this id: %s", err)
		return err
	}
	return nil
}
