package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id           uint           `gorm:"primary_key;auto_increment" json:"id"`
	Username     string         `gorm:"size:255;not null;unique" json:"username"`
	Email        string         `gorm:"size:255;not null;unique" json:"email"`
	Password     []byte         `json:"-"` // vi encrypted va khong muon show tren URL
	AvatarPath   string         `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Posts        []Post         `gorm:"foreignKey:AuthorID;" json:"posts"`
	Comments     []Comment      `gorm:"foreignKey:UserID;" json:"comments"`
	VoteComments []*Comment     `gorm:"many2many:comment_votes;" json:"voted_comments"`
	VotePosts    []*Post        `gorm:"many2many:post_votes;" json:"voted_posts"`
	SavedPosts   []*Post        `gorm:"many2many:save_posts;" json:"saved_posts"`
}

// encrypt the password
func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

// compare the input password with the stored password
func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
