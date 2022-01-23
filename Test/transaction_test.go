package Test

import (
	"altastore/configs"
	"altastore/constant"
	"altastore/delivery/common"
	tc "altastore/delivery/controllers/transaction"
	"altastore/delivery/middlewares"
	tr "altastore/repository/transaction"
	"altastore/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

var invoiceId string

func TestCheckout(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	transactionRepo := tr.NewTransactionRepository(db)
	transactionCon := tc.NewTransactionController(transactionRepo)

	e := echo.New()
	e.Validator = &tc.TransactionValidator{Validator: validator.New()}

	t.Run("Transaction Success", func(t *testing.T) {
		e.POST("/transactions/checkout", transactionCon.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		products := []tc.ProductTransaction{}

		products = append(products, tc.ProductTransaction{
			ProductID:         2,
			Quantity:          3,
		})

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"products": products,
				},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/checkout", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		invoiceId = response.Data.(map[string]interface{})["invoice_id"].(string)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, float64(3000000), response.Data.(map[string]interface{})["total_price"].(float64))
		assert.NotNil(t, response.Data.(map[string]interface{})["payment_url"])
		assert.Equal(t, "PENDING", response.Data.(map[string]interface{})["status"])
	})
	
	t.Run("Transaction Bad Request", func(t *testing.T) {
		e.POST("/transactions/checkout", transactionCon.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		reqBody, _ := json.Marshal(
				map[string]interface{}{},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/checkout", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Transaction Unauthorize", func(t *testing.T) {
		e.POST("/transactions/checkout", transactionCon.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		products := []tc.ProductTransaction{}

		products = append(products, tc.ProductTransaction{
			ProductID:         1,
			Quantity:          3,
		})

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"products": products,
				},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/checkout", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
}

func TestCallback(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	transactionRepo := tr.NewTransactionRepository(db)
	transactionCon := tc.NewTransactionController(transactionRepo)

	e := echo.New()
	e.Validator = &tc.TransactionValidator{Validator: validator.New()}

	t.Run("Callback Success", func(t *testing.T) {
		e.POST("/transactions/callback", transactionCon.Callback)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"external_id": invoiceId,
					"payment_method": "BANK_TRANSFER",
    				"payment_channel": "MANDIRI",
					"paid_at": "2022-01-19T18:13:01.246Z",
					"status": "PAID",
				},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/callback", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", fmt.Sprintf(constant.XENDIT_CALLBACK_TOKEN))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Callback Unautorize", func(t *testing.T) {
		e.POST("/transactions/callback", transactionCon.Callback)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"external_id": invoiceId,
					"payment_method": "BANK_TRANSFER",
    				"payment_channel": "MANDIRI",
					"paid_at": "2022-01-19T18:13:01.246Z",
					"status": "PAID",
				},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/callback", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", fmt.Sprintf(constant.XENDIT_CALLBACK_TOKEN + "adjd"))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotAcceptable, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Callback Not Found", func(t *testing.T) {
		e.POST("/transactions/callback", transactionCon.Callback)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"external_id": invoiceId + "DA",
					"payment_method": "BANK_TRANSFER",
    				"payment_channel": "MANDIRI",
					"paid_at": "2022-01-19T18:13:01.246Z",
					"status": "PAID",
				},
			)

		req := httptest.NewRequest(echo.POST, "/transactions/callback", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", fmt.Sprintf(constant.XENDIT_CALLBACK_TOKEN))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
}

func TestGetTransaction(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	transactionRepo := tr.NewTransactionRepository(db)
	transactionCon := tc.NewTransactionController(transactionRepo)

	e := echo.New()

	t.Run("Get Transaction Success By User", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Get Transaction Not Found By User", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/15", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Get Transaction Bad Request By User", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/hd", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Get Transaction Success By Admin", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenAdmin))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Get Transaction Not Found By Admin", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/15", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenAdmin))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Get Transaction Bad Request By Admin", func(t *testing.T) {
		e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions/zd", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenAdmin))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})

	t.Run("Get All Transaction Success By User", func(t *testing.T) {
		e.GET("/transactions", transactionCon.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Get All Transaction Success By Admin", func(t *testing.T) {
		e.GET("/transactions", transactionCon.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

		req := httptest.NewRequest(echo.GET, "/transactions", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
}