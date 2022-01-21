package entities

type Cart struct {
	UserID uint
	ProductID uint
	Quantity int
	Product Product
}