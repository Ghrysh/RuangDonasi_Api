package models

import "time"

type News struct {
	ID        uint      `gorm:"primaryKey"`
	AdminID   uint      `gorm:"not null"`
	Admin     User      `gorm:"foreignKey:AdminID"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	ImageURL  string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}