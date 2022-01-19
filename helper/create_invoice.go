package helper

import (
	"altastore/entities"
	"os"
	"time"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)


func CreateInvoice(transaction entities.Transaction) (entities.Transaction, error) {
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	data := invoice.CreateParams{
		ExternalID:  "invoice-" + time.Now().String(),
		Amount:      transaction.TotalPrice,
		PayerEmail:  "customer@customer.com",
		Description: "invoice  #1",
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		return transaction, err
	}

	transactionSuccess := entities.Transaction{
		UserID:        transaction.UserID,
		InvoiceID:     resp.ExternalID,
		PaymentMethod: resp.PaymentMethod,
		PaymentUrl:    resp.InvoiceURL,
		TotalPrice:    resp.Amount,
		Status:        resp.Status,
	}

	return transactionSuccess, nil
}