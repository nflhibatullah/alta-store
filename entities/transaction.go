package entities

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	InvoiceID string `gorm:"not null"`
	PaymentMethod string
	BankID string
	PaymentUrl string
	PaidAt time.Time `gorm:"default:null"`
	TotalPrice float64
	Status string `gorm:"not null;default:PENDING"`
	TransactionDetails []TransactionDetail
}