package cart

type CartResponse struct {
	ProductID int `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	Category string `json:"category"`
	Status string `json:"status"`
}