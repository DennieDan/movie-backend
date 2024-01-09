package models

import (
	"time"

	"gorm.io/gorm"
)

type Topic struct {
	Id        uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Name      string         `gorm:"size:255;not null;unique" json:"title"`
	Color     string         `gorm:"size:255;not null" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Posts     []Post         `json:"posts"`
}
