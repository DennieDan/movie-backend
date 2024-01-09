package models

type CommentVotes struct {
	UserID    uint64 `gorm:"primaryKey"`
	CommentID uint64 `gorm:"primaryKey"`
	Score     int
}

type PostVotes struct {
	UserID uint64 `gorm:"primaryKey"`
	PostID uint64 `gorm:"primaryKey"`
	Score  int
}
