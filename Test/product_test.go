package Test

//
//import (
//	"altastore/configs"
//	"altastore/constant"
//	"altastore/delivery/common"
//
//	"altastore/delivery/middlewares"
//	"altastore/entities"
//	proRepo "altastore/repository/category"
//	proRepo "altastore/repository/product"
//	"altastore/utils"
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/go-playground/validator/v10"
//	"github.com/labstack/echo/v4"
//	"github.com/labstack/echo/v4/middleware"
//	"github.com/labstack/gommon/log"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestCreateProduct(t *testing.T) {
//
//	config := configs.GetConfig()
//
//	db := utils.InitDB(config)
//	utils.InitialMigrate(db)
//
//	proRepo := proRepo.NewCategoryRepo(db)
//	categoryContoller := catController.NewCategoryControllers(proRepo)
//
//	e := echo.New()
//
//	t.Run(
//		"Create Category Succes", func(t *testing.T) {
//			e.POST("/categories", categoryContoller.PostCategoryCtrl())
//			e.Validator = &common.CustomValidator{Validator: validator.New()}
//			registerBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": "Handphone",
//				},
//			)
//
//			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//			err := json.Unmarshal(rec.Body.Bytes(), &response)
//			if err != nil {
//				log.Error(err)
//				return
//			}
//
//			assert.Equal(t, http.StatusOK, response.Code)
//			assert.Equal(t, "Successful Operation", response.Message)
//			assert.Equal(t, "Handphone", response.Data.(map[string]interface{})["Name"])
//
//		},
//	)
//
//	t.Run(
//		"Create Cayegory Bad input", func(t *testing.T) {
//			e.POST("/categories", categoryContoller.PostCategoryCtrl())
//
//			registerBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": 123,
//				},
//			)
//
//			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusBadRequest, response.Code)
//			assert.NotNil(t, response.Message)
//		},
//	)
//
//	t.Run(
//		"Create Category Failed (Category Already Exist)", func(t *testing.T) {
//			e.POST("/categories", categoryContoller.PostCategoryCtrl())
//
//			registerBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": "Handphone",
//				},
//			)
//
//			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseError
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusBadRequest, response.Code)
//			assert.NotNil(t, response.Message)
//		},
//	)
//}
//
//func TestGetCategory(t *testing.T) {
//
//	config := configs.GetConfig()
//	db := utils.InitDB(config)
//	utils.InitialMigrate(db)
//	proRepo := proRepo.NewCategoryRepo(db)
//	proController := catController.NewCategoryControllers(proRepo)
//	e := echo.New()
//
//	t.Run(
//		"Get Category by id Success", func(t *testing.T) {
//			e.GET(
//				"/categories/:id", proController.GetCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			req := httptest.NewRequest(echo.GET, "/categories/1", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			err := json.Unmarshal(rec.Body.Bytes(), &response)
//			if err != nil {
//				log.Error(err)
//				return
//			}
//			assert.Equal(t, "Successful Operation", response.Message)
//			assert.Equal(t, http.StatusOK, response.Code)
//
//		},
//	)
//
//	t.Run(
//		"Get Category by id Failed Not Found", func(t *testing.T) {
//
//			e.GET(
//				"/categories/:id", proController.GetCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//				middlewares.CheckRole,
//			)
//
//			req := httptest.NewRequest(echo.GET, "/categories/10000", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//			var response common.ResponseError
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusNotFound, response.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//
//	t.Run(
//		"Get Category Failed Unauthorize", func(t *testing.T) {
//			e.GET(
//				"/categories/:id", proController.GetCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//				middlewares.CheckRole,
//			)
//
//			req := httptest.NewRequest(echo.GET, "/categories/1", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrong"))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusUnauthorized, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//}
//
//func TestGetAllCategory(t *testing.T) {
//	config := configs.GetConfig()
//	db := utils.InitDB(config)
//	utils.InitialMigrate(db)
//	proRepo := proRepo.NewCategoryRepo(db)
//	proController := catController.NewCategoryControllers(proRepo)
//	e := echo.New()
//
//	t.Run(
//		"Get Category Success", func(t *testing.T) {
//			e.GET(
//				"/categories", proController.GetAllCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			req := httptest.NewRequest(echo.GET, "/categories", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//			assert.Equal(t, http.StatusOK, response.Code)
//			assert.Equal(t, "Successful Operation", response.Message)
//
//		},
//	)
//
//	t.Run(
//		"Get Category Failed Not Found", func(t *testing.T) {
//			proRepo.Delete(1)
//			e.GET(
//				"/categories/", proController.GetAllCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//				middlewares.CheckRole,
//			)
//
//			req := httptest.NewRequest(echo.GET, "/categories/", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//			var response common.ResponseError
//			json.Unmarshal(rec.Body.Bytes(), &response)
//			assert.Equal(t, http.StatusNotFound, response.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//
//	t.Run(
//		"Get Category Failed Unauthorize", func(t *testing.T) {
//			e.GET("/categories", proController.GetCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
//
//			req := httptest.NewRequest(echo.GET, "/categories", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusUnauthorized, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//}
//
//func TestUpdateCategory(t *testing.T) {
//	config := configs.GetConfig()
//	db := utils.InitDB(config)
//	utils.InitialMigrate(db)
//	proRepo := proRepo.NewCategoryRepo(db)
//	proController := catController.NewCategoryControllers(proRepo)
//	e := echo.New()
//
//	t.Run(
//		"Update Category Success", func(t *testing.T) {
//			var cat entities.Category
//			db.Unscoped().Where("id = 1").Find(&cat).Model(&cat).Update("deleted_at", nil)
//
//			e.PUT(
//				"/categories/:id", proController.PutCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//			e.Validator = &common.CustomValidator{Validator: validator.New()}
//
//			dataBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": "Laptop",
//				},
//			)
//
//			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusOK, response.Code)
//			assert.Equal(t, "Successful Operation", response.Message)
//			assert.Equal(t, "Laptop", response.Data.(map[string]interface{})["Name"])
//		},
//	)
//
//	t.Run(
//		"Update Category Failed Not Found", func(t *testing.T) {
//			e.PUT(
//				"/categories/:id", proController.PutCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			dataBody, _ := json.Marshal(
//				map[string]interface{}{"name": "Laptop"},
//			)
//
//			req := httptest.NewRequest(echo.PUT, "/categories/5", bytes.NewBuffer(dataBody))
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusNotFound, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//
//	t.Run(
//		"Update Category Failed Bad Request", func(t *testing.T) {
//			e.PUT(
//				"/categories/:id", proController.PutCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//			e.Validator = &common.CustomValidator{Validator: validator.New()}
//			dataBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": 123,
//				},
//			)
//
//			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusBadRequest, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//
//	t.Run(
//		"Update Category Failed Unauthorize", func(t *testing.T) {
//			e.PUT(
//				"/categories/:id", proController.PutCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			dataBody, _ := json.Marshal(
//				map[string]interface{}{
//					"name": "Laptop",
//				},
//			)
//			e.Validator = &common.CustomValidator{Validator: validator.New()}
//
//			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusUnauthorized, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//}
//
//func TestDeleteCategory(t *testing.T) {
//	config := configs.GetConfig()
//	db := utils.InitDB(config)
//	utils.InitialMigrate(db)
//	proRepo := proRepo.NewCategoryRepo(db)
//	proController := catController.NewCategoryControllers(proRepo)
//	e := echo.New()
//	t.Run(
//		"Delete Category Fail not found", func(t *testing.T) {
//			e.DELETE(
//				"/categories/:id", proController.DeleteCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			req := httptest.NewRequest(echo.DELETE, "/categories/1000", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusNotFound, response.Code)
//			assert.Equal(t, "Kategori tidak ditemukan", response.Message)
//		},
//	)
//
//	t.Run(
//		"Delete Category Success", func(t *testing.T) {
//			e.DELETE(
//				"/categories/:id", proController.DeleteCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			req := httptest.NewRequest(echo.DELETE, "/categories/1", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			var response common.ResponseSuccess
//
//			json.Unmarshal(rec.Body.Bytes(), &response)
//
//			assert.Equal(t, http.StatusOK, response.Code)
//			assert.Equal(t, "Successful Operation", response.Message)
//		},
//	)
//
//	t.Run(
//		"Delete Category Failed Unauthorize", func(t *testing.T) {
//			e.DELETE(
//				"/categories/:id", proController.DeleteCategoryCtrl(),
//				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
//			)
//
//			req := httptest.NewRequest(echo.DELETE, "/categories/1", nil)
//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
//			req.Header.Set("Content-Type", "application/json")
//			rec := httptest.NewRecorder()
//
//			e.ServeHTTP(rec, req)
//
//			assert.Equal(t, http.StatusUnauthorized, rec.Code)
//			assert.NotNil(t, rec.Body)
//		},
//	)
//}
