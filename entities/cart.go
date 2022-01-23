package entities

type Cart struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false"`
	ProductID uint `gorm:"primaryKey;autoIncrement:false"`
	Quantity int
	Product Product
}