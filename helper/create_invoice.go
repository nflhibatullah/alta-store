package helper

import (
	"altastore/entities"
	"os"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(transaction entities.Transaction, email string) (entities.Transaction, error) {
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	items := []xendit.InvoiceItem{}

	for _, item := range transaction.TransactionDetails {
		items = append(items, xendit.InvoiceItem{
			Name: item.Product.Name,
			Quantity: item.Quantity,
			Price: float64(item.Product.Price),
		})
	}

	data := invoice.CreateParams{
		ExternalID:                     transaction.InvoiceID,
		Amount:                         transaction.TotalPrice,
		Description:                    "Invoice " + transaction.InvoiceID + " for " + email,
		PayerEmail:                     email,
		Items:                          items,
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		return transaction, err
	}

	transactionSuccess := entities.Transaction{
		PaymentUrl:         resp.InvoiceURL,
		Status:             resp.Status,
	}

	return transactionSuccess, nil
}