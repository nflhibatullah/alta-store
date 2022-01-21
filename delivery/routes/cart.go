package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/cart"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCartPath(e *echo.Echo, cc *controller.CartController) {

	e.GET("/carts", cc.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)
	e.POST("/carts", cc.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)
	e.PUT("/carts/:productId", cc.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)
	e.DELETE("/carts/:productId", cc.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)

}