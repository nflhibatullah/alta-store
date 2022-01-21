package transaction

type PostCustumerRequest struct {
	Products []Item `json:"products" validate:"required"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}