package database

import (
	"TeleEcho/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
func GetUserContacts(userID uint) ([]model.Contact, error) {
	var contacts []model.Contact
	err := DB.Where("user_id = ?", userID).Find(&contacts).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Printf("No contacts found for user ID %d", userID)
			return nil, NotFoundContact
		}
		logrus.Printf("Error retrieving contacts for user ID %d: %s", userID, err)
		return nil, err
	}
	return contacts, nil
}
