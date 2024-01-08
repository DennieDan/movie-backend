package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id         uint           `gorm:"primary_key;auto_increment" json:"id"`
	Username   string         `gorm:"size:255;not null;unique" json:"username"`
	Email      string         `gorm:"size:255;not null;unique" json:"email"`
	Password   []byte         `json:"-"` // vi encrypted va khong muon show tren URL
	AvatarPath string         `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// encrypt the password
func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}
