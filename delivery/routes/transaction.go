package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/transaction"
	"altastore/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTransactionPath(e *echo.Echo, uc *controller.TransactionController) {

	e.GET("/transactions", uc.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.GET("/transactions/:id", uc.GetByTransaction, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)
	e.POST("/transactions/checkout", uc.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)), middlewares.CheckRoleUser)
	e.POST("/transactions/callback", uc.Callback)

}