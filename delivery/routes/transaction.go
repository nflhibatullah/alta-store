package routes

import (
	"altastore/constant"
	controller "altastore/delivery/controllers/transaction"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTransactionPath(e *echo.Echo, uc *controller.TransactionController) {

	e.POST("/transactions", uc.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))

}