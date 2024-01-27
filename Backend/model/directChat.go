package model

import "time"

type DirectChat struct {
	ID         uint      `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	SenderID   uint      `json:"userID" gorm:"foreignKey;not null"`
	ReceiverID uint      `json:"receiverID" gorm:"foreignKey;not null"`
	CreatedAt  time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime"`
}
