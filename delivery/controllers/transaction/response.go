package transaction

import "altastore/entities"

type TransactionResponse struct {
	ID uint `json:"id"`
	InvoiceID string `json:"invoice_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentUrl string `json:"payment_url"`
	TotalPrice float64 `json:"total_price"`
	Products []entities.TransactionDetail `json:"products"`
}