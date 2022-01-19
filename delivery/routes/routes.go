package routes

import (
	"altastore/configs"
	user "altastore/delivery/controllers/users"
	"altastore/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, userCtrl *user.UsersController) {

	e.POST("/users/register", userCtrl.PostUserCtrl())
	e.POST("/users/login", userCtrl.Login())
	e.GET("/users", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
}
