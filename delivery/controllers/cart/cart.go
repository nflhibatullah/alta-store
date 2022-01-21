package cart

import (
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
	"altastore/helper"
	repository "altastore/repository/cart"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartController struct {
	CartRepository repository.Cart 
}

func NewCartController(cart repository.Cart) *CartController {
	return &CartController{CartRepository: cart}
}

func (cc CartController) Create(c echo.Context) error {
	const DEFAULT_ADD = 1

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.NewBadRequestResponse())
	}

	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}
	
	cartData := entities.Cart{
		UserID:    uint(user.ID),
		ProductID: uint(productId),
		Quantity:  DEFAULT_ADD,
	}

	data, err := cc.CartRepository.CheckCart(cartData)
	if err != nil {
		_, err = cc.CartRepository.Create(cartData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	} else {
		cartData.Quantity += data.Quantity
		_, err = cc.CartRepository.Update(cartData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (cc CartController) GetAll(c echo.Context) error {
	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	carts, err := cc.CartRepository.GetAll(user.ID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.NewBadRequestResponse())
	}

	data := []CartResponse{}

	for _, dt := range carts {
		data = append(data, CartResponse{
			ProductID:   int(dt.ProductID),
			ProductName: dt.Product.Name,
			Quantity:    helper.CheckAvailableQuantity(dt.Quantity, dt.Product.Stock),
			Price:       float64(dt.Product.Price),
			TotalPrice:  float64(dt.Product.Price) * float64(dt.Quantity),
			Category:    dt.Product.Category.Name,
			Status: helper.CheckProductStatus(dt.Product.Stock),
		})
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(data))
}

func (cc CartController) Update(c echo.Context) error {
	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	var cartRequest UpdateCartRequest

	if err := c.Bind(&cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if err := c.Validate(&cartRequest); err != nil {
      return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
    }

	data := entities.Cart{
		UserID:    uint(user.ID),
		ProductID: uint(productId),
		Quantity:  cartRequest.Quantity,
	}

	_, err = cc.CartRepository.Update(data)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.ErrorResponse(http.StatusNotFound, err.Error()))
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (cc CartController) Delete(c echo.Context) error {
	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err = cc.CartRepository.Delete(user.ID, productId)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}