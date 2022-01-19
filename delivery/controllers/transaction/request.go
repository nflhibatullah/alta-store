package transaction

type PostCustumerRequest struct {
	TotalPrice float64 `json:"total_price"`
	Products []Item `json:"products"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}