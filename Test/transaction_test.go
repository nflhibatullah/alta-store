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
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		invoiceId = response.Data.(map[string]interface{})["invoice_id"].(string)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, float64(750000), response.Data.(map[string]interface{})["total_price"])
		assert.NotNil(t, response.Data.(map[string]interface{})["payment_url"])
		assert.Equal(t, "PENDING", response.Data.(map[string]interface{})["status"])
	})
	
	t.Run("Transaction Bad Request", func(t *testing.T) {
		e.POST("/transactions/checkout", transactionCon.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		products := []tc.ProductTransaction{}

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

func TestGetTransaction(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	transactionRepo := tr.NewTransactionRepository(db)
	transactionCon := tc.NewTransactionController(transactionRepo)

	e := echo.New()

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
	
	// t.Run("Get By ID Transaction Success By User", func(t *testing.T) {
	// 	e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

	// 	req := httptest.NewRequest(echo.GET, "/transactions/1", nil)
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
	// 	rec := httptest.NewRecorder()
		
	// 	e.ServeHTTP(rec, req)

	// 	var response common.ResponseSuccess
	// 	json.Unmarshal(rec.Body.Bytes(), &response)

	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, "Successful Operation", response.Message)
	// })
	
	// t.Run("Get By ID Transaction Success By Admin", func(t *testing.T) {
	// 	e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

	// 	req := httptest.NewRequest(echo.GET, "/transactions/1", nil)
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		
	// 	rec := httptest.NewRecorder()
		
	// 	e.ServeHTTP(rec, req)

	// 	var response common.ResponseSuccess
	// 	json.Unmarshal(rec.Body.Bytes(), &response)

	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, "Successful Operation", response.Message)
	// })
	
	// t.Run("Get By ID Transaction Not Found By User", func(t *testing.T) {
	// 	e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

	// 	req := httptest.NewRequest(echo.GET, "/transactions/15", nil)
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
	// 	rec := httptest.NewRecorder()
		
	// 	e.ServeHTTP(rec, req)

	// 	var response common.ResponseSuccess
	// 	json.Unmarshal(rec.Body.Bytes(), &response)

	// 	assert.Equal(t, http.StatusNotFound, rec.Code)
	// 	assert.NotNil(t, response.Message)
	// })
	
	// t.Run("Get By ID Transaction Bad Request", func(t *testing.T) {
	// 	e.GET("/transactions/:id", transactionCon.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

	// 	req := httptest.NewRequest(echo.GET, "/transactions/1sj5", nil)
	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
	// 	rec := httptest.NewRecorder()
		
	// 	e.ServeHTTP(rec, req)

	// 	var response common.ResponseSuccess
	// 	json.Unmarshal(rec.Body.Bytes(), &response)

	// 	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 	assert.NotNil(t, response.Message)
	// })
	
}