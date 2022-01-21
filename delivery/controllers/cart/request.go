package cart

type PostCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity"`
}