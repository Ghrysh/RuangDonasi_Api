package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	Phone     *string   `gorm:"type:varchar(30);unique"`
	Password  *string   `gorm:"type:varchar(255)"`
	Role      string    `gorm:"type:varchar(20);default:'user'"` 
	GoogleID  *string   `gorm:"type:varchar(255);unique"`
	IsActive  bool		`gorm:"default:true" json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}