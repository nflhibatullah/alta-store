package transaction

import (
	"altastore/delivery/common"
	"altastore/entities"
	"altastore/helper"
	repository "altastore/repository/transaction"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionRepository repository.Transaction
}

func NewTransactionController(transaction repository.Transaction) *TransactionController {
	return &TransactionController{TransactionRepository: transaction}
}

func (tc TransactionController) Create(c echo.Context) error {
	var transactionRequest PostCustumerRequest

	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userId := 1
	invoiceId := "INV-" + time.Now().String() + fmt.Sprint(userId)

	transaction := entities.Transaction{
		UserID:        uint(userId),
		InvoiceID:     invoiceId,
		TotalPrice:    transactionRequest.TotalPrice,
		Status:        "PENDING",
	}

	transactionData, err := tc.TransactionRepository.Create(transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	for _, item := range transactionRequest.Products {
		product := entities.TransactionDetail{
			TransactionID: transactionData.ID,
			ProductID:     uint(item),
			Quantity:      2,
		}
		_, err := tc.TransactionRepository.StoreItemProduct(int(transactionData.ID), product)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
		}
	}

	transactionSuccess, err := helper.CreateInvoice(transactionData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactionSuccess))
}