package transaction

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransactionRequest struct {
	Products []Item `json:"products" form:"products" validate:"required"`
}

type Item struct {
	ProductID uint `json:"product_id" form:"product_id"`
	Quantity int `json:"quantity" form:"product_id"`
}

type TransactionValidator struct {
	Validator *validator.Validate
}

func (cv *TransactionValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}