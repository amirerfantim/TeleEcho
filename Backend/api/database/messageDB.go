package database

import (
	"TeleEcho/model"
	"github.com/sirupsen/logrus"
	"time"
)

func CreateMessage(chatID uint, messageType model.ChatType, content string) (*model.Message, error) {
	newMessage := model.Message{
		ChatID:    chatID,
		Type:      messageType,
		Content:   content,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	if err := DB.Create(&newMessage).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"chatID":  chatID,
			"type":    messageType,
			"content": content,
		}).WithError(err).Error("Failed to create message")
		return nil, err
	}
	logrus.Printf("Message with id %d created successfully", newMessage.ID)

	return &newMessage, nil
}
func GetMessagesByChatIDAndType(chatID uint, chatType model.ChatType) ([]model.Message, error) {
	var messages []model.Message
	err := DB.Where("chat_id = ? AND type = ?", chatID, chatType).Find(&messages).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatID": chatID,
			"type":   chatType,
		}).WithError(err).Error("Failed to retrieve messages")
		return nil, err
	}
	return messages, nil
}
