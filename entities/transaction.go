package entities

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	InvoiceID string `gorm:"not null"`
	PaymentMethod string
	PaymentUrl string
	TotalPrice float64 `gorm:"not null"`
	Status string `gorm:"not null"`
	Transaction []TransactionDetail
}