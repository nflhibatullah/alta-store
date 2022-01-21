package transaction

import (
	"time"
)

type TransactionResponse struct {
	ID uint `json:"id"`
	InvoiceID string `json:"invoice_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentUrl string `json:"payment_url"`
	BankID string `json:"bank_id"`
	PaidAt time.Time `json:"pay_at"`
	TotalPrice float64 `json:"total_price"`
	Status string `json:"status"`
	Products []ProductTransaction`json:"products"`
}

type ProductTransaction struct {
	ProductID int `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	TotalProductPrice float64 `json:"total_product_price"`
	Quantity int `json:"quantity"`
	Category string `json:"category"`
}