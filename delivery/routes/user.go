package routes

import (
	"altastore/constant"
	user "altastore/delivery/controllers/users"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterUserPath(e *echo.Echo, userCtrl *user.UsersController) {

	e.POST("/register", userCtrl.PostUserCtrl())
	e.POST("/login", userCtrl.Login())
	e.GET("/users", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
}
