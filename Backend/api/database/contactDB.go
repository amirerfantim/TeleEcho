package database

import (
	"TeleEcho/model"
	"github.com/sirupsen/logrus"
)

func CreateContact(userID, contactUserID uint, status model.Status) error {
	contactData := &model.Contact{UserID: userID, ContactUserID: contactUserID, Status: status}
	if err := DB.Create(contactData).Error; err != nil {
		logrus.Printf("Error creating contact:%s\n", err)
		return err
	}
	logrus.Printf("Contact created successfully\n")
	return nil
}
