package models

import "time"

type Article struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Excerpt   string    `gorm:"type:text;not null"`
	Content   string    `gorm:"type:text"`
	Category  string    `gorm:"type:varchar(100);not null"`
	ReadTime  int       `gorm:"not null"`
	Color     string    `gorm:"type:varchar(20);default:'#2a1a0a'"`
	Accent    string    `gorm:"type:varchar(20);default:'#fb923c'"`
	Image     string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}