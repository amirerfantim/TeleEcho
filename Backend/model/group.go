package model

type Group struct {
	ID          uint   `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	Name        string ` json:"name,omitempty" gorm:"type:varchar(255);not null"`
	Description string `json:"description,omitempty" gorm:"type:varchar(255);default=''"`
}

func NewGroup(id uint, name, description string) *Group {
	return &Group{ID: id, Name: name, Description: description}
}
