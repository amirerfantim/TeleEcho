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
		logrus.WithFields(logrus.Fields{
			"adminUserID": adminUserID,
			"groupName":   groupName,
			"error":       err,
		}).Error("Error while checking for group existence")
		return false, err
	}
	return count > 0, nil
}
func DoesGroupExistByID(adminUserID, groupID uint) (bool, error) {
	var count int64
	err := DB.Model(&model.Group{}).Where("admin_user_id = ? AND id = ?", adminUserID, groupID).Count(&count).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"adminUserID": adminUserID,
			"id":          groupID,
			"error":       err,
		}).Error("Error while checking for group existence")
		return false, err
	}
	return count > 0, nil
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
func RemoveUserFromGroup(userID, groupID uint) error {
	result := DB.Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&model.UserGroup{})
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"userID":  userID,
			"groupID": groupID,
			"error":   result.Error,
		}).Error("Failed to remove user from group")
		return result.Error
	}
	if result.RowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"userID":  userID,
			"groupID": groupID,
		}).Warn("No user-group relation found to delete")
		return NotUserInGroup
	}

	logrus.WithFields(logrus.Fields{
		"userID":  userID,
		"groupID": groupID,
	}).Info("User removed from group successfully")
	return nil
}
func GetAllUsersInGroup(groupID uint) ([]model.User, error) {
	var users []model.User
	err := DB.Table("user_groups").
		Select("users.*").
		Joins("join users on users.id = user_groups.user_id").
		Where("user_groups.group_id = ?", groupID).
		Scan(&users).Error

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"groupID": groupID,
			"error":   err,
		}).Error("Failed to retrieve users from group")
		return nil, err
	}

	return users, nil
}

func GetGroupByID(id uint) (*model.Group, error) {
	var group model.Group

	if err := DB.First(&group, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Warnf("Group with ID %d not found", id)
			return nil, NotFoundGroup
		}

		logrus.WithError(err).Error("Failed to fetch group by ID")
		return nil, err
	}

	return &group, nil
}
func FindGroupChatByGroupID(groupID uint) (*model.GroupChat, error) {
	var groupChat model.GroupChat
	result := DB.Where("group_id = ?", groupID).First(&groupChat)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{
				"groupID": groupID,
			}).Info("No GroupChat found for the given GroupID")
			return nil, NotFoundChat
		}
		logrus.WithFields(logrus.Fields{
			"groupID": groupID,
		}).WithError(result.Error).Error("Failed to retrieve GroupChat")
		return nil, result.Error
	}

	return &groupChat, nil
}
