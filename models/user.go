package models

type User struct {
	Id         uint   `gorm:"primary_key;auto_increment" json:"id"`
	Username   string `gorm:"size:255;not null;unique" json:"username"`
	Email      string `gorm:"size:100;not null;unique" json:"email"`
	Password   []byte `json:"-"` // vi encrypted va khong muon show tren URL
	AvatarPath string `gorm:"size:255;null;" json:"avatar_path"`
}
