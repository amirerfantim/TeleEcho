package model

type User struct {
	ID             uint   `json:"id" gorm:"primary_key;auto_increment;<-:create"`
	Username       string `json:"username" gorm:"type:varchar(255);not null;unique"`
	Firstname      string `json:"firstname" gorm:"type:varchar(255);not null"`
	Lastname       string `json:"lastname" gorm:"type:varchar(255);default:''"`
	Phone          string `json:"phone" gorm:"type:varchar(255);not null;unique" `
	Password       string `json:"password" gorm:"type:varchar(255);not null"`
	ProfilePicture string `json:"profilePicture,omitempty" gorm:"type:varchar(255);unique;default:''"`
	Bio            string `json:"Bio,omitempty" gorm:"type:varchar(255);default:''"`
}

func NewUser(id uint, username, firstname, lastname, phone, password, profilePicture, bio string) *User {
	return &User{ID: id, Username: username, Firstname: firstname, Lastname: lastname, Phone: phone, Password: password, ProfilePicture: profilePicture, Bio: bio}
}
