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
func DeleteUserContact(userID uint, contactUserID uint) error {
	var contact model.Contact
	err := DB.Where("contact_user_id = ? AND user_id = ?", contactUserID, userID).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Printf("No contact found with ContactUserID %d for user ID %d", contactUserID, userID)
			return NotFoundContact
		}
		logrus.Printf("Error retrieving contact for deletion: %s", err)
		return err
	}
	if err := DB.Delete(&contact).Error; err != nil {
		logrus.Printf("Error deleting contact with ContactUserID %d: %s", contactUserID, err)
		return err
	}
	logrus.Printf("Contact with ContactUserID %d for UserID %d deleted successfully", contactUserID, userID)
	return nil
}
func UpdateContactStatus(userID uint, contactUserID uint, newStatus model.Status) error {
	var contact model.Contact
	err := DB.Where("user_id = ? AND contact_user_id = ?", userID, contactUserID).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Printf("No contact found with UserID %d and ContactUserID %d", userID, contactUserID)
			return NotFoundContact
		}
		logrus.Printf("Error retrieving contact: %s", err)
		return err
	}
	contact.Status = newStatus
	if err := DB.Save(&contact).Error; err != nil {
		logrus.Printf("Error updating status for contact with UserID %d and ContactUserID %d: %s", userID, contactUserID, err)
		return err
	}
	logrus.Printf("Contact status updated successfully for UserID %d and ContactUserID %d", userID, contactUserID)
	return nil
}

type UserWithContactStatus struct {
	model.User
	ProfilePictureHide bool `json:"profilePictureHide"`
	PhoneNumberHide    bool `json:"phoneNumberHide"`
	IsBlocked          bool `json:"isBlocked"`
}

func GetUserContactsInfo(userID uint) ([]model.User, error) {
	var usersInfo []UserWithContactStatus
	var contactUsers []model.User
	err := DB.Table("contacts").
		Select("users.id, users.username, users.firstname, users.lastname, users.phone, users.profile_picture, users.bio, contacts.profile_picture_hide, contacts.phone_number_hide, contacts.is_blocked").
		Joins("left join users on users.id = contacts.contact_user_id").
		Where("contacts.user_id = ?", userID).
		Scan(&usersInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFoundContact
		}
		logrus.Printf("Error retrieving user contacts info: %s", err)
		return nil, err
	}
	for _, userInfo := range usersInfo {
		if userInfo.PhoneNumberHide {
			userInfo.User.Phone = "hide"
		}
		if userInfo.ProfilePictureHide {
			userInfo.User.ProfilePicture = "hide"
		}
		if userInfo.IsBlocked {
			userInfo.User.Phone = "hide"
			userInfo.User.Bio = "hide"
			userInfo.User.ProfilePicture = "hide"
		}
		contactUsers = append(contactUsers, userInfo.User)
	}

	return contactUsers, nil
}
