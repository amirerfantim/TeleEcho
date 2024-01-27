package model

import "time"

type GroupChat struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	GroupID   uint      `json:"groupID" gorm:"foreignKey;not null"`
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime"`
}
