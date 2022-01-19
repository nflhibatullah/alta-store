package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/product"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterProductPath(e *echo.Echo, pc *controller.ProductController)  {
	e.POST("/product", pc.PostProductCtrl(), middlewares.CheckRole)
	e.GET("/product", pc.GetAllProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.GET("/product/:id", pc.GetProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.DELETE(
		"/product/:id", pc.DeleteProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole,
	)
	e.PUT("/product/:id", pc.PutProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
}