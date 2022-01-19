package user

import (
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
	"altastore/repository/users"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UsersController struct {
	Repo users.UserInterface
}

func NewUsersControllers(usrep users.UserInterface) *UsersController {
	return &UsersController{Repo: usrep}
}
func (uscon UsersController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var login entities.User
		c.Bind(&login)

		user, err := uscon.Repo.Login(login.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User tidak ditemukan")
		}

		hash, err := middlewares.Checkpwd(user.Password, login.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Ada kesalahan dalam kredensial")
		}

		var token string

		if hash {
			token, _ = middlewares.CreateToken(int(user.ID), user.Role)
		}

		return c.JSON(http.StatusOK, token)

	}
}

func (uscon UsersController) PostUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		newUserReq := RegisterUserRequestFormat{}

		if err := c.Bind(&newUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(newUserReq.Password), 14)
		newUser := entities.User{
			Name:     newUserReq.Name,
			Email:    newUserReq.Email,
			Password: string(hash),
		}

		_, err := uscon.Repo.Create(newUser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

func (uscon UsersController) GetAllUsersCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		user, err := uscon.Repo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		response := GetUsersResponseFormat{
			Message: "Successful Opration",
			Data:    user,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// GET /users/:id
func (uscon UsersController) GetUserCtrl() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := middlewares.ExtractTokenUser(c)
		
		user, err := uscon.Repo.Get(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "success",
				"data":    user,
			},
		)
	}

}

// PUT /users/:id
func (uscon UsersController) EditUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, _ := middlewares.ExtractTokenUser(c)

		updateUserReq := PutUserRequestFormat{}
		if err := c.Bind(&updateUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(updateUserReq.Password), 14)

		updateUser := entities.User{
			Name:     updateUserReq.Name,
			Email:    updateUserReq.Email,
			Password: string(hash),
		}

		if _, err := uscon.Repo.Update(updateUser, id); err != nil {
			log.Error(err)
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

// DELETE /users/:id
func (uscon UsersController) DeleteUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, _ := middlewares.ExtractTokenUser(c)
		deleteUserReq := DeleteRequestFormat{}
		if err := c.Bind(&deleteUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		user, err := uscon.Repo.GetDeleteData(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User tidak ditemukan")
		}

		_, err = middlewares.Checkpwd(user.Password, deleteUserReq.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Ada kesalahan dalam kredensial")
		}

		uscon.Repo.Delete(id)

		return c.JSON(http.StatusOK, "Berhasil menghapus user")

	}

}
