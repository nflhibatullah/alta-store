package routes

import (
	controller "altastore/delivery/controllers/transaction"

	"github.com/labstack/echo/v4"
)

func RegisterTransactionPath(e *echo.Echo, uc controller.TransactionController) {

	e.POST("/transactions", uc.Create)

}