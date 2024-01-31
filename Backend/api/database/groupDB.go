package database

import (
	"TeleEcho/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateGroup(adminUserID uint, name, description, profilePath string) error {
	groupData := &model.Group{AdminUserID: adminUserID, Name: name, Description: description, GroupProfilePicture: profilePath}
	if err := DB.Create(groupData).Error; err != nil {
		logrus.Printf("Error creating group:%s\n", err)
		return err
	}
	logrus.Printf("Group created successfully\n")

	userGroup := &model.UserGroup{
		UserID:  adminUserID,
		GroupID: groupData.ID,
	}
	if err := DB.Create(userGroup).Error; err != nil {
		logrus.Printf("Error creating user-group association: %s\n", err)
		return err
	}
	logrus.Printf("User-group association created successfully\n")
	return nil
}
func DoesGroupExist(adminUserID uint, groupName string) (bool, error) {
	var count int64
	err := DB.Model(&model.Group{}).Where("admin_user_id = ? AND name = ?", adminUserID, groupName).Count(&count).Error
	if err != nil {
		logrus.Printf("Error while checking user groups: %s", err)
		return false, err
	}
	return count == 0, nil
}

func GetUserGroups(userID uint) ([]model.Group, error) {
	var groups []model.Group

	err := DB.Table("user_groups").
		Select("groups.*").
		Joins("join groups on groups.id = user_groups.group_id").
		Where("user_groups.user_id = ?", userID).
		Scan(&groups).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Printf("No groups found for user ID: %d", userID)
			return nil, NotFoundGroup
		}
		logrus.Printf("Error retrieving joined groups for user ID %d: %s", userID, err)
		return nil, err
	}

	if len(groups) == 0 {
		return nil, NotFoundGroup
	}

	return groups, nil

}
func AddUserToGroup(userID, groupID uint) error {
	userGroup := model.UserGroup{
		UserID:  userID,
		GroupID: groupID,
	}

	result := DB.Create(&userGroup)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  userID,
			"groupID": groupID,
			"error":   result.Error,
		}).Error("Failed to add user to group")
		return result.Error
	}

	logrus.WithFields(logrus.Fields{
		"userID":  userID,
		"groupID": groupID,
	}).Info("User added to group successfully")
	return nil
}
func IsUserInGroup(userID, groupID uint) (bool, error) {
	var userGroup model.UserGroup
	result := DB.Where("user_id = ? AND group_id = ?", userID, groupID).First(&userGroup)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		logrus.WithFields(logrus.Fields{
			"userID":  userID,
			"groupID": groupID,
			"error":   result.Error,
		}).Error("Error occurred while checking user group membership")
		return false, result.Error
	}
	return true, nil
}
