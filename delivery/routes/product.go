package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/product"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterProductPath(e *echo.Echo, pc *controller.ProductController) {
	e.POST("/products", pc.PostProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.GET(
		"/products", pc.GetAllProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)),
		middlewares.CheckRole,
	)
	e.GET("/products/:id", pc.GetProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.DELETE(
		"/products/:id", pc.DeleteProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole,
	)
	e.PUT("/products/:id", pc.PutProductCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
}
