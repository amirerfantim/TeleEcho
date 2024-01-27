package model

type Status struct {
	ProfilePictureHide bool
	PhoneNumberHide    bool
	IsBlocked          bool
}

type Contact struct {
	ID            uint   `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	UserID        uint   `json:"userID" gorm:"foreignKey;not null"`
	ContactUserID uint   `json:"contactUserID" gorm:"foreignKey;not null"`
	Status        Status `json:"status" gorm:"embedded"`
}

func NewContact(id, userID, contactUserID uint, status Status) *Contact {
	return &Contact{ID: id, UserID: userID, ContactUserID: contactUserID, Status: status}
}
