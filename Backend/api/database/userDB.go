package database

import (
	"TeleEcho/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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
		logrus.Printf("Can not remove user with this id: %s", err)
		return err
	}
	return nil
}
func GetUserByUserID(userID uint) (*model.User, error) {
	var user model.User
	if err := DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFoundUser
		}
		logrus.Printf("Error while querying user: %s", err)
		return nil, err
	}
	return &user, nil
}
func GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFoundUser
		}
		logrus.Printf("Error while querying user by phone: %s", err)
		return nil, err
	}
	return &user, nil
}
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFoundUser
		}
		logrus.Printf("Error while querying user by username: %s", err)
		return nil, err
	}
	return &user, nil
}

func UpdateUserByUserID(userID uint, userUpdates model.User) error {
	var existingUser model.User
	if err := DB.First(&existingUser, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NotFoundUser
		}
		logrus.Printf("Error while querying user: %s", err)
		return err
	}

	if err := DB.Model(&existingUser).Updates(userUpdates).Error; err != nil {
		logrus.Printf("Error while updating user: %s", err)
		return err
	}
	logrus.Printf("User updated successfully")

	return nil
}
