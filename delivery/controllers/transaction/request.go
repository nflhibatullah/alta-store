package transaction

type PostCustumerRequest struct {
	TotalPrice float64 `json:"total_price"`
	Products []uint `json:"products"`
}