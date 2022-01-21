package cart

import (
	"altastore/delivery/common"
	"altastore/delivery/middlewares"
	"altastore/entities"
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
	var cartRequest PostCartRequest

	if err := c.Bind(&cartRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, err := middlewares.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	cartData := entities.Cart{
		UserID:    uint(user.ID),
		ProductID: cartRequest.ProductID,
		Quantity:  cartRequest.Quantity,
	}

	_, err = cc.CartRepository.Create(cartData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.ErrorResponse(http.StatusInternalServerError, err.Error()))
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
		return c.JSON(http.StatusUnauthorized, common.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(carts))
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

	data := entities.Cart{
		UserID:    uint(user.ID),
		ProductID: uint(productId),
		Quantity:  cartRequest.Quantity,
	}

	_, err = cc.CartRepository.Update(data)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.ErrorResponse(http.StatusUnauthorized, err.Error()))
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