package models

type Movie struct {
	Id    uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Title string `gorm:"size:255;not null;unique" json:"title"`
	Year  uint64 `gorm:"not null" json:"year"`
	Posts []Post `json:"posts"`
}
