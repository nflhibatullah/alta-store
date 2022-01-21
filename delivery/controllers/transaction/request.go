package transaction

type PostCustumerRequest struct {
	Products []Item `json:"products"`
}

type Item struct {
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
}