package models

import "time"

type Post struct {
	Id        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	MovieID   uint32    `gorm:"" json:"movie_id"`
	TopicID   uint32    `gorm:"" json:"topic_id"`
	AuthorID  uint32    `gorm:"not null;" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
