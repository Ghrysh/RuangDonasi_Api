package models

import "time"

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	CampaignID  uint      `gorm:"not null"`
	Campaign    Campaign  `gorm:"foreignKey:CampaignID"`
	Amount      float64   `gorm:"not null"`
	UniqueCode  int       `gorm:"not null"`
	TotalAmount float64   `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);default:'pending'"`
	PaymentURL  string    `gorm:"type:text"`                  
	CreatedAt   time.Time
	UpdatedAt   time.Time
}