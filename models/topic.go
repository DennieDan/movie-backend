package models

type Topic struct {
	Id    uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name  string `gorm:"size:255;not null;unique" json:"title"`
	Color string `gorm:"size:255;not null" json:"color"`
	// Posts []Post `json:"posts"`
}
