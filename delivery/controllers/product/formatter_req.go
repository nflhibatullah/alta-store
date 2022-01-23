package product

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateProductRequestFormat struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Price       int    `json:"price" form:"price" validate:"required"`
	Stock       int    `json:"stock" form:"stock" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	CategoryID  uint   `json:"category_id" form:"category_id" validate:"required"`
}

type PutProductRequestFormat struct {
	Name        string `json:"name" form:"name"`
	Price       int    `json:"price" form:"price"`
	Stock       int    `json:"stock" form:"stock"`
	Description string `json:"description" form:"description"`
	CategoryID  uint   `json:"category_id" form:"category_id"`
}

type ProductValidator struct {
	Validator *validator.Validate
}

func (cv *ProductValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
