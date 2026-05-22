package models

type AppSetting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"type:varchar(50);unique;not null"`
	Value string `gorm:"type:varchar(255);not null"`
}