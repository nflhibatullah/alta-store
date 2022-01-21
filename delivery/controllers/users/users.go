package user

import (
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
	"altastore/repository/users"

	"net/http"

	"golang.org/x/crypto/bcrypt"

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
			return echo.NewHTTPError(http.StatusNotFound, common.ErrorResponse(404, "Users not found"))
		}

		hash, err := middlewares.Checkpwd(user.Password, login.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Kesalahan pada password"))
		}

		var token string

		if hash {
			token, _ = middlewares.CreateToken(int(user.ID), user.Role, user.Email)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(token))

	}
}

func (uscon UsersController) PostUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		newUserReq := RegisterUserRequestFormat{}
		c.Bind(&newUserReq)

		if err := c.Validate(newUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Harap isi data dengan baik"))
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(newUserReq.Password), 14)
		newUser := entities.User{
			Name:     newUserReq.Name,
			Email:    newUserReq.Email,
			Password: string(hash),
		}

		_, err := uscon.Repo.Create(newUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "email telah terdaftar"))
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

func (uscon UsersController) GetAllUsersCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		user, _ := uscon.Repo.GetAll()
		if len(user) == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User tidak ditemukan"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(user))
	}
}

// GET /users/:id
func (uscon UsersController) GetUserCtrl() echo.HandlerFunc {
	return func(c echo.Context) error {
		userToken, err := middlewares.ExtractTokenUser(c) 
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		user, err := uscon.Repo.Get(userToken.ID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User tidak ditemukan"))
			
		}
		
		return c.JSON(http.StatusOK, common.SuccessResponse(user))

	}
}

// PUT /users/:id
func (uscon UsersController) EditUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		user, _ := middlewares.ExtractTokenUser(c)

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

		if _, err := uscon.Repo.Update(updateUser, user.ID); err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}
		
	return c.JSON(http.StatusOK, common.SuccessResponse(user))
	}
}

// DELETE /users/:id
func (uscon UsersController) DeleteUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		userjwt, _ := middlewares.ExtractTokenUser(c)
		deleteUserReq := DeleteRequestFormat{}
		c.Bind(&deleteUserReq)

		user, err := uscon.Repo.GetDeleteData(userjwt.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, common.ErrorResponse(404, "User tidak ditemukan"))
		}

		_, err = middlewares.Checkpwd(user.Password, deleteUserReq.Password)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Kesalahan pada kredensial"))
		}

		uscon.Repo.Delete(userjwt.ID)

		return c.JSON(http.StatusOK, common.SuccessResponse("Berhasil menghapus user"))

	}

}
