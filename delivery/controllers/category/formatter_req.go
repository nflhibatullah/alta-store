package category

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateCategoryRequestFormat struct {
	Name string `json:"name" form:"name"  validate:"required"`
}

type PutCategoryRequestFormat struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type CategoryValidator struct {
	Validator *validator.Validate
}

func (cv *CategoryValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
