package entities

import "gorm.io/gorm"

type TransactionDetail struct {
	gorm.Model
	TransactionID uint `gorm:"not null"`
	ProductID uint `gorm:"not null"`
	Quantity uint `gorm:"not null"`
}