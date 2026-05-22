package models

import "time"

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null"`
	Slug string `gorm:"type:varchar(100);unique;not null"`
}

type Campaign struct {
	ID               uint      `gorm:"primaryKey"`
	CategoryID       uint      `gorm:"not null"`
	Category         Category  `gorm:"foreignKey:CategoryID"`
	Title            string    `gorm:"type:varchar(255);not null"`
	Location         string    `gorm:"type:varchar(255);not null;default:'Belum diatur'"`
	Description      string    `gorm:"type:text"`
	TargetAmount     float64   `gorm:"not null"`
	CurrentAmount    float64   `gorm:"default:0"`
	EndDate          time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	IsUrgent         bool      `gorm:"default:false"`
	Fundraiser       string    `gorm:"type:varchar(255);not null;default:'Admin'"`
	DonatorCount     int       `gorm:"default:0"`
	ImageURL         string    `gorm:"type:varchar(255)"`
	BeneficiaryCount int       `gorm:"default:0"`
	Status           string    `gorm:"type:varchar(20);default:'active'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}