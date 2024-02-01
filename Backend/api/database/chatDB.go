package database

import (
	"TeleEcho/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateDirectChat(senderID, receiverID uint) (*model.DirectChat, error) {
	chat := model.DirectChat{SenderID: senderID, ReceiverID: receiverID}
	result := DB.Create(&chat)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"senderID":   senderID,
			"receiverID": receiverID,
			"error":      result.Error,
		}).Error("Failed to create direct chat")
		return nil, result.Error
	}

	chat = model.DirectChat{SenderID: receiverID, ReceiverID: senderID}
	result = DB.Create(&chat)
	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"senderID":   receiverID,
			"receiverID": senderID,
			"error":      result.Error,
		}).Error("Failed to create direct chat")
		return nil, result.Error
	}
	logrus.Printf("Direct chat between %d and %d created successfully\n", senderID, receiverID)
	return &chat, nil
}
func CreateGroupChat(groupID uint) (*model.GroupChat, error) {
	chat := model.GroupChat{GroupID: groupID}
	result := DB.Create(&chat)

	if result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"groupID": groupID,
			"error":   result.Error,
		}).Error("Failed to create group chat")
		return nil, result.Error
	}
	logrus.Printf("Group chat in group id %d created successfully\n", groupID)

	return &chat, nil
}

func DoesDirectChatExist(senderID, receiverID uint) (bool, error) {
	var chat model.DirectChat
	result := DB.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).Or("sender_id = ? AND receiver_id = ?", receiverID, senderID).First(&chat)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		logrus.WithFields(logrus.Fields{
			"senderID":   senderID,
			"receiverID": receiverID,
			"error":      result.Error,
		}).Error("Failed to check direct chat existence")
		return false, result.Error
	}

	return true, nil
}

func DoesGroupChatExist(groupID uint) (bool, error) {
	var chat model.GroupChat
	result := DB.Where("group_id = ?", groupID).First(&chat)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		logrus.WithFields(logrus.Fields{
			"groupID": groupID,
			"error":   result.Error,
		}).Error("Failed to check group chat existence")
		return false, result.Error
	}

	return true, nil
}
func DeleteDirectChatAndMessages(chatID uint) error {
	tx := DB.Begin()
	if err := tx.Where("chat_id = ? AND type = ?", chatID, model.TypeDirectChat).Delete(&model.Message{}).Error; err != nil {
		logrus.WithError(err).WithField("chatID", chatID).Error("Failed to delete direct chat messages")
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", chatID).Delete(&model.DirectChat{}).Error; err != nil {
		logrus.WithError(err).WithField("chatID", chatID).Error("Failed to delete direct chat")
		tx.Rollback()
		return err
	}

	tx.Commit()
	logrus.Println("Direct chat with id %d with all messages deleted successfully.", chatID)
	return nil
}

func DeleteGroupChatAndMessages(chatID uint) error {
	tx := DB.Begin()
	if err := tx.Where("chat_id = ? AND type = ?", chatID, model.TypeGroupChat).Delete(&model.Message{}).Error; err != nil {
		logrus.WithError(err).WithField("chatID", chatID).Error("Failed to delete group chat messages")
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", chatID).Delete(&model.GroupChat{}).Error; err != nil {
		logrus.WithError(err).WithField("chatID", chatID).Error("Failed to delete group chat")
		tx.Rollback()
		return err
	}

	tx.Commit()
	logrus.Println("Group chat with id %d with all messages deleted successfully.", chatID)
	return nil
}
