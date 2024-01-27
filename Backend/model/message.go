package model

import "time"

type Type string

const (
	TypeDirectChat Type = "DirectChat"
	TypeGroupChat  Type = "GroupChat"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	ChatID    uint      `json:"chatID" gorm:"foreignKey;not null"`
	Type      Type      ` json:"type" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:varchar(1000);not null"`
	IsRead    bool      ` json:"isRead,omitempty" gorm:"type:bool;default:false"`
	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime"`
}

func NewMessage(id, chatId uint, chatType Type, content string, isRead bool, createdAT time.Time) *Message {
	return &Message{ID: id, ChatID: chatId, Type: chatType, Content: content, IsRead: isRead, CreatedAt: createdAT}
}
