package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/category"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCategoryPath(e *echo.Echo, cc *controller.CategoryController)  {
	e.POST("/category", cc.PostCategoryCtrl(), middlewares.CheckRole)
	e.GET("/category", cc.GetAllCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.GET("/category/:id", cc.GetCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.DELETE(
		"/category/:id", cc.DeleteCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole,
	)
	e.PUT("/category/:id", cc.PutCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
}