package model

type UserGroup struct {
	ID      uint `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	UserID  uint ` json:"userID" gorm:"foreignKey;not null"`
	GroupID uint `json:"groupID" gorm:"foreignKey;not null"`
}
