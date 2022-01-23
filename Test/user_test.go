package Test

import (
	"altastore/configs"
	"altastore/constant"
	"altastore/delivery/common"
	userController "altastore/delivery/controllers/users"
	"altastore/delivery/middlewares"
	userRepo "altastore/repository/users"
	"altastore/utils"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	//"fmt"
	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {

	config := configs.GetConfig()

	db := utils.InitDB(config)
	utils.InitialMigrate(db)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()

	t.Run(
		"Register Success 1", func(t *testing.T) {
			e.POST("/register", userContoller.PostUserCtrl())
			e.Validator = &userController.UserValidator{Validator: validator.New()}
			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name":     "Arif",
					"email":    "arif@gmail.com",
					"password": "123",
				},
			)

			req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(registerBody))
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
			assert.Equal(t, "Arif", response.Data.(map[string]interface{})["name"])
			assert.Equal(t, "arif@gmail.com", response.Data.(map[string]interface{})["email"])

		},
	)
	t.Run(
		"Register Bad input", func(t *testing.T) {
			e.POST("/register", userContoller.PostUserCtrl())

			registerBody, _ := json.Marshal(
				map[string]interface{}{
					"name": "Naufal",

					"password": "123",
				},
			)

			req := httptest.NewRequest(echo.POST, "/register", bytes.NewBuffer(registerBody))
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
		"Register Failed (Email Already Exist)", func(t *testing.T) {
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

			var response common.ResponseError

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.NotNil(t, response.Message)
		},
	)
}

func TestLoginUser(t *testing.T) {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()

	t.Run(
		"Login Success", func(t *testing.T) {
			e.POST("/login", userContoller.Login())

			loginBody, _ := json.Marshal(
				map[string]interface{}{
					"email":    "naufal@gmail.com",
					"password": "123",
				},
			)

			req := httptest.NewRequest(echo.POST, "/login", bytes.NewBuffer(loginBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			token = response.Data.(string)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.NotNil(t, response.Data)

		},
	)
	
	t.Run(
		"Login Success user", func(t *testing.T) {
			e.POST("/login", userContoller.Login())

			loginBody, _ := json.Marshal(
				map[string]interface{}{
					"email":    "arif@gmail.com",
					"password": "123",
				},
			)

			req := httptest.NewRequest(echo.POST, "/login", bytes.NewBuffer(loginBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			tokenUser = response.Data.(string)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.NotNil(t, response.Data)

		},
	)

	t.Run(
		"Login Failed (User Not Found)", func(t *testing.T) {
			e.POST("/login", userContoller.Login())

			loginBody, _ := json.Marshal(
				map[string]interface{}{
					"email":    "furqonzt98332@gmail.com",
					"password": "1234qwesdrdad",
				},
			)

			req := httptest.NewRequest(echo.POST, "/login", bytes.NewBuffer(loginBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.NotNil(t, response.Message)
		},
	)

	t.Run(
		"Login Failed (Wrong Password)", func(t *testing.T) {
			e.POST("/login", userContoller.Login())

			loginBody, _ := json.Marshal(
				map[string]interface{}{
					"email":    "naufal@gmail.com",
					"password": "12345",
				},
			)

			req := httptest.NewRequest(echo.POST, "/login", bytes.NewBuffer(loginBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.NotNil(t, response.Message)
		},
	)
}

func TestGetUser(t *testing.T) {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()

	t.Run(
		"Get User by id Success", func(t *testing.T) {
			e.GET("/users/profile", userContoller.GetUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/users/profile", nil)
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
		"Get User by id Failed Not Found", func(t *testing.T) {

			e.GET(
				"/users/", userContoller.GetUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)
			wrongToken, _ := middlewares.CreateToken(4, "admin", "fufu@gmail.com")

			req := httptest.NewRequest(echo.GET, "/users/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", wrongToken))
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
		"Get User Failed Unauthorize", func(t *testing.T) {
			e.GET("/users/profile", userContoller.GetUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/users/profile", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestGetAllUser(t *testing.T) {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()

	t.Run(
		"Get User Success", func(t *testing.T) {
			e.GET("/users", userContoller.GetAllUsersCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/users", nil)
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
		"Get User Failed Not Found", func(t *testing.T) {
			userRepo.Delete(2)
			userRepo.Delete(3)
			e.GET(
				"/users/", userContoller.GetAllUsersCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
				middlewares.CheckRole,
			)
			wrongToken, _ := middlewares.CreateToken(8, "admin", "")

			req := httptest.NewRequest(echo.GET, "/users/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", wrongToken))
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
		"Get User Failed Unauthorize", func(t *testing.T) {
			e.GET("/users/profile", userContoller.GetAllUsersCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.GET, "/users/profile", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestUpdateUser(t *testing.T) {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()

	t.Run(
		"Update User Success", func(t *testing.T) {
			e.PUT("/users/update", userContoller.EditUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
			e.Validator = &userController.UserValidator{Validator: validator.New()}
			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name":     "Naufal Aammar Hibatullah",
					"email":    "naufalaammar@gmail.com",
					"password": "12345",
				},
			)

			req := httptest.NewRequest(echo.PUT, "/users/update", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)
			fmt.Println(response)
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "Successful Operation", response.Message)
			assert.Equal(t, "Naufal Aammar Hibatullah", response.Data.(map[string]interface{})["name"])
			assert.Equal(t, "naufalaammar@gmail.com", response.Data.(map[string]interface{})["email"])
		},
	)

	t.Run(
		"Update User Failed Not Found", func(t *testing.T) {
			e.PUT("/users/update", userContoller.EditUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name":     "Naufal Aammmar Hibatullah",
					"email":    "naufalaammar@gmail.com",
					"password": "12345",
				},
			)

			wrongToken, _ := middlewares.CreateToken(5, "user", "fufu@gmail.com")

			req := httptest.NewRequest(echo.PUT, "/users/update", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", wrongToken))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update User Failed Bad Request", func(t *testing.T) {
			e.PUT("/users/update", userContoller.EditUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name":     "Naufal Aammar Hibatullah",
					"email":    "naufalaammar@gmail.com",
					"password": 1234,
				},
			)

			req := httptest.NewRequest(echo.PUT, "/users/update", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Update User Failed Unauthorize", func(t *testing.T) {
			e.PUT("/users/update", userContoller.EditUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"name":     "Naufal Aammar Hibatullah",
					"email":    "naufalaammar@gmail.com",
					"password": "12345",
				},
			)

			req := httptest.NewRequest(echo.PUT, "/users/update", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}

func TestDeleteUser(t *testing.T) {
	config := configs.GetConfig()

	db := utils.InitDB(config)

	userRepo := userRepo.NewUsersRepo(db)
	userContoller := userController.NewUsersControllers(userRepo)

	e := echo.New()
	t.Run(
		"Delete User Fail wrong password", func(t *testing.T) {
			e.DELETE("/users/delete", userContoller.DeleteUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"password": "12345aaaaa",
				},
			)

			req := httptest.NewRequest(echo.DELETE, "/users/delete", bytes.NewBuffer(dataBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var response common.ResponseSuccess

			json.Unmarshal(rec.Body.Bytes(), &response)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "Kesalahan pada kredensial", response.Message)
		},
	)

	t.Run(
		"Delete User Success", func(t *testing.T) {
			e.DELETE("/users/delete", userContoller.DeleteUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
			dataBody, _ := json.Marshal(
				map[string]interface{}{
					"password": "12345",
				},
			)

			req := httptest.NewRequest(echo.DELETE, "/users/delete", bytes.NewBuffer(dataBody))
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
		"Delete User Failed Not Found", func(t *testing.T) {
			e.DELETE("/users/delete", userContoller.DeleteUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			wrongToken, _ := middlewares.CreateToken(17, "user", "fufu@gmail.com")

			req := httptest.NewRequest(echo.DELETE, "/users/delete", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", wrongToken))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)

	t.Run(
		"Delete User Failed Unauthorize", func(t *testing.T) {
			e.DELETE("/users/delete", userContoller.DeleteUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

			req := httptest.NewRequest(echo.DELETE, "/users/delete", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token+"wrongtoken"))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.NotNil(t, rec.Body)
		},
	)
}
