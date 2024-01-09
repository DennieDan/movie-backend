package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	Id        uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Title     string         `gorm:"size:255;not null;unique" json:"title"`
	Year      uint64         `gorm:"not null" json:"year"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Posts     []Post         `json:"posts"`
}
