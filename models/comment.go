package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Id         uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserID     uint64 `json:"user_id"`
	PostID     uint64 `json:"post_id"`
	ResponseID *uint64
	Body       string         `gorm:"text;not null;" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Responses  []Comment      `gorm:"foreignKey:ResponseID;constraint:OnDelete:CASCADE;"`
	Voters     []*User        `gorm:"many2many:comment_votes;" json:"voters"`
}
