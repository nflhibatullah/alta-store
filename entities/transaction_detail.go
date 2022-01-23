package entities

type TransactionDetail struct {
	TransactionID uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity int `gorm:"not null"`
	Product Product
}