package model

type Group struct {
	ID                  uint   `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	AdminUserID         uint   `json:"adminUserID" gorm:"foreignKey;not null"`
	Name                string ` json:"name,omitempty" gorm:"type:varchar(255);not null"`
	Description         string `json:"description,omitempty" gorm:"type:varchar(255);default=''"`
	GroupProfilePicture string `json:"groupProfilePicture,omitempty" gorm:"type:varchar(255);unique;default:''"`
}

func NewGroup(id uint, name, description string) *Group {
	return &Group{ID: id, Name: name, Description: description}
}
