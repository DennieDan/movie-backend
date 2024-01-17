package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Id        uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Title     string         `gorm:"size:255;not null;unique" json:"title"`
	Content   string         `gorm:"text;not null;" json:"content"`
	MovieID   uint64         `json:"movie_id"`
	TopicID   uint64         `json:"topic_id"`
	AuthorID  uint64         `gorm:"not null;" json:"author_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Movie     Movie          `json:"movie"`
	Topic     Topic          `json:"topic"`
	Author    User           `gorm:"foreignKey:AuthorID;references:Id" json:"author"`
	Comments  []Comment      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"comments"`
	Voters    []*User        `gorm:"many2many:post_votes;" json:"voters"`
	Savers    []*User        `gorm:"many2many:save_posts;" json:"savers"`
}
