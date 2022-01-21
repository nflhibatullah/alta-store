package transaction

import (
	"altastore/constant"
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
	"altastore/helper"
	repository "altastore/repository/transaction"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	// bind request data
	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, _ := middlewares.ExtractTokenUser(c)
	invoiceId := uuid.New().String()

	// create data to db
	transaction := entities.Transaction{
		UserID:        uint(user.ID),
		InvoiceID:     invoiceId,
		TotalPrice:    transactionRequest.TotalPrice,
	}

	transactionData, err := tc.TransactionRepository.Create(transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// store Detail Transaction Data 
	for _, item := range transactionRequest.Products {
		product := entities.TransactionDetail{
			TransactionID: transactionData.ID,
			ProductID:     uint(item.ProductID),
			Quantity:      item.Quantity,
		}
		_, err := tc.TransactionRepository.StoreItemProduct(int(transactionData.ID), product)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
		}
	}

	// get created & stored data from db
	transactionDb, err := tc.TransactionRepository.GetByTransaction(user.ID, int(transactionData.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// send data to payment gateway
	transactionSuccess, err := helper.CreateInvoice(transactionDb, user.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// update data to db
	transactionUpdate, err := tc.TransactionRepository.Update(int(transactionData.ID), transactionSuccess)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactionUpdate))
}

func (tc TransactionController) Callback(c echo.Context) error {

	req := c.Request()
	headers := req.Header

	xCallbackToken := headers.Get("X-Callback-Token")

	if xCallbackToken != constant.XENDIT_CALLBACK_TOKEN {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	var callbackRequest common.CallbackRequest
	if err := c.Bind(&callbackRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	transactionData, err := tc.TransactionRepository.GetByInvoiceId(callbackRequest.ExternalID) 
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	var data entities.Transaction
	data.PaidAt, _ = time.Parse(time.RFC3339, callbackRequest.PaidAt)
	data.PaymentMethod = callbackRequest.PaymentMethod
	data.BankID = callbackRequest.BankID
	data.Status = callbackRequest.Status

	_, err = tc.TransactionRepository.Update(int(transactionData.ID), data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) GetAll(c echo.Context) error {

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	transactions, err := tc.TransactionRepository.GetAll()

	if user.Role == "user" {
		transactions, err = tc.TransactionRepository.GetByUser(user.ID)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactions))
}

func (tc TransactionController) GetByTransaction(c echo.Context) error {
	transactionId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	transaction, err := tc.TransactionRepository.GetByTransaction(user.ID, transactionId)

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transaction))
}