package models

type Category struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:200" json:"description"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
}
