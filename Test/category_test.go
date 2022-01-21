package Test

import (
	"altastore/configs"
	"altastore/constant"
	"altastore/delivery/common"
	catController "altastore/delivery/controllers/category"
	userController "altastore/delivery/controllers/users"
	"altastore/delivery/middlewares"
	"altastore/entities"
	categoryRepo "altastore/repository/category"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"net/http"

	userRepo "altastore/repository/users"

	"altastore/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"testing"
)

var token string

func TestMain(m *testing.M) {
	config := configs.GetConfig()

	db := utils.InitDB(config)
	db.Migrator().DropTable(&entities.User{})
	db.Migrator().DropTable(&entities.Product{})
	db.Migrator().DropTable(&entities.Category{})
	utils.InitialMigrate(db)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()
	e.Validator = &common.CustomValidator{Validator: validator.New()}
	e.POST("/register", userContoller.PostUserCtrl())

	registerBody, _ := json.Marshal(
		map[string]interface{}{
			"name":     "Naufal",
			"email":    "naufal@gmail.com",
			"password": "123",
		},
	)

	req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	updateUser := entities.User{
		Role: "admin",
	}

	userRepo.Update(updateUser, 1)

	e.POST("/login", userContoller.Login())

	loginBody, _ := json.Marshal(
		map[string]interface{}{
			"email":    "naufal@gmail.com",
			"password": "123",
		},
	)

	reqLogin := httptest.NewRequest(echo.POST, "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	recLogin := httptest.NewRecorder()

	e.ServeHTTP(recLogin, reqLogin)

	var response common.ResponseSuccess

	json.Unmarshal(recLogin.Body.Bytes(), &response)

	token = response.Data.(string)

	fmt.Println("TEST MAIN JALAN NIH")

	m.Run()

}

func TestCreateCategory(t *testing.T) {

	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryContoller := catController.NewCategoryControllers(categoryRepo)

	e := echo.New()

	t.Run(
		"Create Category Succes", func(t *testing.T) {
			e.POST("/categories", categoryContoller.PostCategoryCtrl())
			e.Validator = &common.CustomValidator{Validator: validator.New()}
			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Handphone",
				},
			)

			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
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
			assert.Equal(t, "Handphone", response.Data.(map[string]interface{})["Name"])

		},
	)

	t.Run(
		"Create Cayegory Bad input", func(t *testing.T) {
			e.POST("/categories", categoryContoller.PostCategoryCtrl())

			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name": 123,
				},
			)

			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.NotNil(t, response.Message)
		},
	)

	t.Run(
		"Create Category Failed (Category Already Exist)", func(t *testing.T) {
			e.POST("/categories", categoryContoller.PostCategoryCtrl())

			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Handphone",
				},
			)

			req := httptest.NewRequest(echo.POST, "/categories", bytes.NewBuffer(registerBody))
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

func TestGetCategory(t *testing.T) {

	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryController := catController.NewCategoryControllers(categoryRepo)
	e := echo.New()

	t.Run(
		"Get Category by id Success", func(t *testing.T) {
			e.GET(
				"/categories/:id", categoryController.GetCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.GET, "/categories/1", nil)
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
		"Get Category by id Failed Not Found", func(t *testing.T) {

			e.GET(
				"/categories/:id", categoryController.GetCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/categories/10000", nil)
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
			e.GET(
				"/categories/:id", categoryController.GetCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/categories/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrong"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestGetAllCategory(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryController := catController.NewCategoryControllers(categoryRepo)
	e := echo.New()

	t.Run(
		"Get Category Success", func(t *testing.T) {
			e.GET(
				"/categories", categoryController.GetAllCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.GET, "/categories", nil)
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
		"Get Category Failed Not Found", func(t *testing.T) {
			categoryRepo.Delete(1)
			e.GET(
				"/categories/", categoryController.GetAllCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)

			req := httptest.NewRequest(echo.GET, "/categories/", nil)
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
			e.GET("/categories", categoryController.GetCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/categories", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestUpdateCategory(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryController := catController.NewCategoryControllers(categoryRepo)
	e := echo.New()

	t.Run(
		"Update Category Success", func(t *testing.T) {
			var cat entities.Category
			db.Unscoped().Where("id = 1").Find(&cat).Model(&cat).Update("deleted_at", nil)

			e.PUT(
				"/categories/:id", categoryController.PutCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)
			e.Validator = &common.CustomValidator{Validator: validator.New()}

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Laptop",
				},
			)

			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.Equal(t, "Laptop", response.Data.(map[string]interface{})["Name"])
		},
	)

	t.Run(
		"Update Category Failed Not Found", func(t *testing.T) {
			e.PUT(
				"/categories/:id", categoryController.PutCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			dataBody, _ := json.Marshal(
				map[string]interface{}{"name": "Laptop"},
			)

			req := httptest.NewRequest(echo.PUT, "/categories/5", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update Category Failed Bad Request", func(t *testing.T) {
			e.PUT(
				"/categories/:id", categoryController.PutCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)
			e.Validator = &common.CustomValidator{Validator: validator.New()}
			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": 123,
				},
			)

			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update Category Failed Unauthorize", func(t *testing.T) {
			e.PUT(
				"/categories/:id", categoryController.PutCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Laptop",
				},
			)
			e.Validator = &common.CustomValidator{Validator: validator.New()}

			req := httptest.NewRequest(echo.PUT, "/categories/1", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestDeleteCategory(t *testing.T) {
	config := configs.GetConfig()
	db := utils.InitDB(config)
	utils.InitialMigrate(db)
	categoryRepo := categoryRepo.NewCategoryRepo(db)
	categoryController := catController.NewCategoryControllers(categoryRepo)
	e := echo.New()
	t.Run(
		"Delete Category Fail not found", func(t *testing.T) {
			e.DELETE(
				"/categories/:id", categoryController.DeleteCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/categories/1000", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, "Kategori tidak ditemukan", response.Message)
		},
	)

	t.Run(
		"Delete Category Success", func(t *testing.T) {
			e.DELETE(
				"/categories/:id", categoryController.DeleteCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/categories/1", nil)
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
		"Delete Category Failed Unauthorize", func(t *testing.T) {
			e.DELETE(
				"/categories/:id", categoryController.DeleteCategoryCtrl(),
				middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
			)

			req := httptest.NewRequest(echo.DELETE, "/categories/1", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}
