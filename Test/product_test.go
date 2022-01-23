package Test

import (
	"altastore/configs"
	"altastore/constant"
	"altastore/delivery/common"

	proController "altastore/delivery/controllers/product"
	"altastore/delivery/middlewares"
	proRepo "altastore/repository/product"
	"altastore/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProduct(t *testing.T) {

	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	proRepo := proRepo.NewProductRepo(db)
	proContoller := proController.NewProductControllers(proRepo)

	e := echo.New()

	t.Run(
		"Create Product Succes", func(t *testing.T) {
			e.POST("/product", proContoller.PostProductCtrl())
			e.Validator = &proController.ProductValidator{Validator: validator.New()}
			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name":        "Xiaomi",
					"price":       1000000,
					"stock":       10,
					"description": "hp kentang",
					"category_id": 1,
				},
			)

			req := httptest.NewRequest(echo.POST, "/product", bytes.NewBuffer(registerBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				log.Error(err)
				return
			}

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.Equal(t, "Xiaomi", response.Data.(map[string]interface{})["Name"])

		},
	)
	t.Run(
		"Create Product Failed (Bad Request)", func(t *testing.T) {
			e.POST("/product", proContoller.PostProductCtrl())
			e.Validator = &proController.ProductValidator{Validator: validator.New()}
			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Handphone",
				},
			)

			req := httptest.NewRequest(echo.POST, "/product", bytes.NewBuffer(registerBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseError

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.NotNil(t, response.Message)
		},
	)
}

func TestGetProduct(t *testing.T) {

	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	proRepo := proRepo.NewProductRepo(db)
	proController := proController.NewProductControllers(proRepo)
	e := echo.New()

	t.Run(
		"Get Product by id Success", func(t *testing.T) {
			e.GET(
				"/product/:id", proController.GetProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.GET, "/product/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			err := json.Unmarshal(rec.Body.Bytes(), &response)
			if err != nil {
				log.Error(err)
				return
			}
			assert.Equal(t, "Successful Operation", response.Message)
			assert.Equal(t, http.StatusOK, response.Code)

		},
	)

	t.Run(
		"Get Product by id Failed Not Found", func(t *testing.T) {

			e.GET(
				"/product/:id", proController.GetProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/product/10000", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			var response common.ResponseError
			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Get Product Failed Unauthorize", func(t *testing.T) {
			e.GET(
				"/product/:id", proController.GetProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/product/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrong"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestGetAllProduct(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	proRepo := proRepo.NewProductRepo(db)
	proController := proController.NewProductControllers(proRepo)
	e := echo.New()

	t.Run(
		"Get Product Success", func(t *testing.T) {
			e.GET(
				"/products/", proController.GetAllProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.GET, "/products/?page=1&perpage=10&search=Xiaomi", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)

		},
	)
	t.Run(
		"Get Product Success 2", func(t *testing.T) {
			e.GET(
				"/products/", proController.GetAllProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.GET, "/products/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)

		},
	)

	t.Run(
		"Get Product Failed Not Found", func(t *testing.T) {

			e.GET(
				"/products/", proController.GetAllProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/products/?page=1&perpage=10&search=apple", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			var response common.ResponseError
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Get Category Failed Unauthorize", func(t *testing.T) {
			e.GET("/products/", proController.GetAllProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/products/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestUpdateProduct(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	proRepo := proRepo.NewProductRepo(db)
	proController := proController.NewProductControllers(proRepo)
	e := echo.New()

	t.Run(
		"Update product Success", func(t *testing.T) {

			e.PUT(
				"/products/:id", proController.PutProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)
			e.Validator = &common.CustomValidator{Validator: validator.New()}

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Oppo",
				},
			)

			req := httptest.NewRequest(echo.PUT, "/products/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.Equal(t, "Oppo", response.Data.(map[string]interface{})["Name"])
		},
	)

	t.Run(
		"Update Product Failed Not Found", func(t *testing.T) {
			e.PUT(
				"/product/:id", proController.PutProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			dataBody, _ := json.Marshal(
				map[string]interface{}{"name": "Laptop"},
			)

			req := httptest.NewRequest(echo.PUT, "/products/200", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update Product Failed Bad Request", func(t *testing.T) {
			e.PUT(
				"/products/:id", proController.PutProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": 123,
				},
			)

			req := httptest.NewRequest(echo.PUT, "/products/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update Products Failed Unauthorize", func(t *testing.T) {
			e.PUT(
				"/products/:id", proController.PutProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Laptop",
				},
			)

			req := httptest.NewRequest(echo.PUT, "/products/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestDeleteProduct(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	proRepo := proRepo.NewProductRepo(db)
	proController := proController.NewProductControllers(proRepo)
	e := echo.New()
	t.Run(
		"Delete Product Fail not found", func(t *testing.T) {
			e.DELETE(
				"/products/:id", proController.DeleteProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/products/1000", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, "Produk tidak ditemukan", response.Message)
		},
	)

	t.Run(
		"Delete Product Success", func(t *testing.T) {
			e.DELETE(
				"/products/:id", proController.DeleteProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/products/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
		},
	)

	t.Run(
		"Delete Product Failed Unauthorize", func(t *testing.T) {
			e.DELETE(
				"/categories/:id", proController.DeleteProductCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/products/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}
