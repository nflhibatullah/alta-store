package Test

import (
	"altastore/configs"
	"altastore/constant"
	"altastore/delivery/common"
	"altastore/delivery/controllers/cart"
	cartCon "altastore/delivery/controllers/cart"
	"altastore/delivery/middlewares"
	cartRepo "altastore/repository/cart"
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

func TestAddCart(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	cartRepo := cartRepo.NewCartRepository(db)
	cartController := cartCon.NewCartController(cartRepo)

	e := echo.New()

	t.Run("Cart Add Success", func(t *testing.T) {
		e.POST("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.POST, "/carts/1", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Cart Internal Server Error", func(t *testing.T) {
		e.POST("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.POST, "/carts/9", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Cart Unauthorize", func(t *testing.T) {
		e.POST("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.POST, "/carts/8jhj", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Cart Bad Request", func(t *testing.T) {
		e.POST("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.POST, "/carts/1", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NotNil(t, response.Message)
	})
}

func TestGetCart(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	cartRepo := cartRepo.NewCartRepository(db)
	cartController := cartCon.NewCartController(cartRepo)

	e := echo.New()

	t.Run("Cart Get Success", func(t *testing.T) {
		e.GET("/carts", cartController.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.GET, "/carts", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Cart Invalid Token", func(t *testing.T) {
		e.POST("/carts", cartController.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.POST, "/carts", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser + "sdjdsj"))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.NotNil(t, response.Message)
	})
}

func TestUpdateCart(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	cartRepo := cartRepo.NewCartRepository(db)
	cartController := cartCon.NewCartController(cartRepo)

	e := echo.New()
	e.Validator = &cart.CartValidator{Validator: validator.New()}

	t.Run("Cart Update Success", func(t *testing.T) {
		e.PUT("/carts/:productId", cartController.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"quantity": 8,
				},
			)

		req := httptest.NewRequest(echo.PUT, "/carts/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Cart Update Bad Request", func(t *testing.T) {
		e.PUT("/carts/:productId", cartController.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"quantity": "jh",
				},
			)

		req := httptest.NewRequest(echo.PUT, "/carts/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Cart Update Bad Request 2", func(t *testing.T) {
		e.PUT("/carts/:productId", cartController.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		reqBody, _ := json.Marshal(
				map[string]interface{}{},
			)

		req := httptest.NewRequest(echo.PUT, "/carts/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Cart Update Not Found", func(t *testing.T) {
		e.PUT("/carts/:productId", cartController.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		reqBody, _ := json.Marshal(
				map[string]interface{}{
					"quantity": 9,
				},
			)

		req := httptest.NewRequest(echo.PUT, "/carts/18", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
}

func TestDeleteCart(t *testing.T)  {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	cartRepo := cartRepo.NewCartRepository(db)
	cartController := cartCon.NewCartController(cartRepo)

	e := echo.New()

	t.Run("Cart Delete Success", func(t *testing.T) {
		e.DELETE("/carts/:productId", cartController.Delete, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.DELETE, "/carts/1", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Successful Operation", response.Message)
	})
	
	t.Run("Cart Internal Server Error", func(t *testing.T) {
		e.DELETE("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.DELETE, "/carts/9", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.NotNil(t, response.Message)
	})
	
	t.Run("Cart Bad Request", func(t *testing.T) {
		e.DELETE("/carts/:productId", cartController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

		req := httptest.NewRequest(echo.DELETE, "/carts/8jhj", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenUser))
		
		rec := httptest.NewRecorder()
		
		e.ServeHTTP(rec, req)

		var response common.ResponseSuccess
		json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.NotNil(t, response.Message)
	})
}
