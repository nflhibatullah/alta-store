package routes

import (
	"altastore/configs"
	"altastore/delivery/controllers/category"
	"altastore/delivery/controllers/product"
	user "altastore/delivery/controllers/users"
	"altastore/delivery/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(
	e *echo.Echo, userCtrl *user.UsersController, productCtrl *product.ProductController,
	categoryCtrl *category.CategoryController,
) {

	e.POST("/users/register", userCtrl.PostUserCtrl())
	e.POST("/users/login", userCtrl.Login())
	e.GET("/users", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.GET("/users/:id", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.DELETE("/users", userCtrl.DeleteUserCtrl(), middleware.JWT([]byte(configs.SecretKey)))
	e.PUT("/users", userCtrl.EditUserCtrl(), middleware.JWT([]byte(configs.SecretKey)))

	e.POST("/product", productCtrl.PostProductCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.GET("/product", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.GET("/product/:id", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.DELETE(
		"/product/:id", userCtrl.DeleteUserCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole,
	)
	e.PUT("/product/:id", userCtrl.EditUserCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)

	e.POST("/category", userCtrl.PostUserCtrl(), middlewares.CheckRole)
	e.GET("/category", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.GET("/category/:id", userCtrl.GetAllUsersCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)
	e.DELETE(
		"/category/:id", userCtrl.DeleteUserCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole,
	)
	e.PUT("/category/:id", userCtrl.EditUserCtrl(), middleware.JWT([]byte(configs.SecretKey)), middlewares.CheckRole)

}
