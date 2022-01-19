package user

import (
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
	"altastore/repository/users"
	"net/http"
	"strconv"

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
		tmpUser := entities.User{}

		if err := c.Bind(&tmpUser); err != nil {
			return err
		}

		res, err := uscon.Repo.Login(tmpUser.Email, tmpUser.Password)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Data tidak ditemukan")
		}

		tokenJWT, _ := middlewares.CreateToken(int(res.ID))
		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "Berhasil login",
				"data":    tokenJWT,
			},
		)
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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		updateUserReq := PutUserRequestFormat{}
		if err := c.Bind(&updateUserReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		updateUser := entities.User{
			Name:     updateUserReq.Name,
			Password: updateUserReq.Password,
		}

		if _, err := uscon.Repo.Update(updateUser, id); err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}

// DELETE /users/:id
func (uscon UsersController) DeleteUserCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		deletedUser, _ := uscon.Repo.Delete(id)

		if deletedUser.ID != 0 {
			return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
		} else {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

	}

}
