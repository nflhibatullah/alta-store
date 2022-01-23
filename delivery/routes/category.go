package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/category"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterCategoryPath(e *echo.Echo, cc *controller.CategoryController)  {
	e.POST("/categories", cc.PostCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.GET("/categories", cc.GetAllCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.GET("/categories/:id", cc.GetCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.DELETE("/categories/:id", cc.DeleteCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
	e.PUT("/categories/:id", cc.PutCategoryCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRole)
}