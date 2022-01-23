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
	"strings"
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
	var transactionRequest TransactionRequest

	// bind request data
	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	if err := c.Validate(&transactionRequest); err != nil {
      return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
    }

	user, _ := middlewares.ExtractTokenUser(c)
	invoiceId := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))

	// create data to db
	transaction := entities.Transaction{
		UserID: uint(user.ID),
		InvoiceID:     invoiceId,
	}

	transactionData, err := tc.TransactionRepository.Create(transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
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
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
	}

	// get data from db
	transactionDb, err := tc.TransactionRepository.GetByTransaction(user.ID, int(transactionData.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var totalPrice float64 = 0

	products := []ProductTransaction{}

	for _, p := range transactionDb.TransactionDetails {
		products = append(products, ProductTransaction{
			ProductID:   int(p.ProductID),
			ProductName: p.Product.Name,
			ProductPrice: float64(p.Product.Price),
			TotalProductPrice: float64(p.Product.Price) * float64(p.Quantity),
			Quantity:    p.Quantity,
			Category:    p.Product.Category.Name,
		})

		totalPrice += (float64(p.Product.Price) * float64(p.Quantity))

		// update product stock
		ok := tc.TransactionRepository.UpdateStockProduct(int(p.ProductID), p.Product.Stock - p.Quantity); 
		if !ok {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		// delete product in cart
		ok = tc.TransactionRepository.GetProductInCart(user.ID, int(p.ProductID))
		if ok {
			tc.TransactionRepository.DeleteProductInCart(user.ID, int(p.ProductID))
		}
	}

	// assign total price
	transactionDb.TotalPrice = totalPrice

	// send data to payment gateway
	transactionSuccess, err := helper.CreateInvoice(transactionDb, user.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	// update data to db
	transactionUpdate, err := tc.TransactionRepository.Update(int(transactionData.ID), transactionSuccess)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}
	
	// get data from db
	transactionGet, err := tc.TransactionRepository.GetByTransaction(user.ID, int(transactionUpdate.ID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	// formatting response
	data := TransactionResponse{
		ID:            transactionGet.ID,
		InvoiceID:     transactionGet.InvoiceID,
		PaymentMethod: transactionGet.PaymentMethod,
		PaymentUrl:    transactionGet.PaymentUrl,
		BankID:        transactionGet.BankID,
		PaidAt:        transactionGet.PaidAt,
		TotalPrice:    transactionGet.TotalPrice,
		Status: transactionGet.Status,
		Products:      products,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(data))
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

	if callbackRequest.Status != "PAID" {
		for _, p := range transactionData.TransactionDetails {
			// update product stock
			ok := tc.TransactionRepository.UpdateStockProduct(int(p.ProductID), p.Product.Stock + p.Quantity); 
			if !ok {
				return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
			}
		}
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
		return c.JSON(http.StatusUnauthorized, common.NewBadRequestResponse())
	}

	transactions, err := tc.TransactionRepository.GetAll()

	if user.Role == "user" {
		transactions, err = tc.TransactionRepository.GetByUser(user.ID)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	transactionDatas := []TransactionResponse{}

	for _, td := range transactions {
		products := []ProductTransaction{}

		for _, p := range td.TransactionDetails {
			products = append(products, ProductTransaction{
				ProductID:         int(p.ProductID),
				ProductName:       p.Product.Name,
				ProductPrice:      float64(p.Product.Price),
				TotalProductPrice: float64(p.Product.Price) * float64(p.Quantity),
				Quantity:          p.Quantity,
				Category:          p.Product.Category.Name,
			})
		}

		transactionDatas = append(transactionDatas, TransactionResponse{
			ID:            td.ID,
			InvoiceID:     td.InvoiceID,
			PaymentMethod: td.PaymentMethod,
			PaymentUrl:    td.PaymentUrl,
			BankID:        td.BankID,
			PaidAt:        td.PaidAt,
			TotalPrice:    td.TotalPrice,
			Status: td.Status,
			Products:      products,
		})
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactionDatas))
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

	if user.Role == "admin" {
		transaction, err = tc.TransactionRepository.GetByTransactionAdmin(transactionId)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	products := []ProductTransaction{}

	for _, p := range transaction.TransactionDetails {
		products = append(products, ProductTransaction{
			ProductID:         int(p.ProductID),
			ProductName:       p.Product.Name,
			ProductPrice:      float64(p.Product.Price),
			TotalProductPrice: float64(p.Product.Price) * float64(p.Quantity),
			Quantity:          p.Quantity,
			Category:          p.Product.Category.Name,
		})
	}

	data := TransactionResponse{
		ID:            transaction.ID,
		InvoiceID:     transaction.InvoiceID,
		PaymentMethod: transaction.PaymentMethod,
		PaymentUrl:    transaction.PaymentUrl,
		BankID:        transaction.BankID,
		PaidAt:        transaction.PaidAt,
		TotalPrice:    transaction.TotalPrice,
		Status: transaction.Status,
		Products:      products,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(data))
}