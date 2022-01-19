package product

import (
	"altastore/delivery/common"
	"altastore/entities"
	"altastore/repository/product"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProductController struct {
	Repo product.ProductInterface
}

func NewProductControllers(prorep product.ProductInterface) *ProductController {
	return &ProductController{Repo: prorep}
}

//CreateProduct

func (procon ProductController) PostProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		newProductReq := CreateProductRequestFormat{}
		if err := c.Bind(&newProductReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		newProduct := entities.Products{
			Name:        newProductReq.Name,
			Price:       newProductReq.Price,
			Stock:       newProductReq.Stock,
			Description: newProductReq.Description,
		}

		_, err := procon.Repo.Create(newProduct)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}
func (procon ProductController) GetAllProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		product, err := procon.Repo.GetAll()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "success",
				"data":    product,
			},
		)
	}

}
func (procon ProductController) GetProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		product, _ := procon.Repo.Get(id)

		if len(product) == 0 {
			return c.JSON(
				http.StatusNotFound, map[string]interface{}{
					"message": "Product not found",
				},
			)
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "succes",
				"data":    product,
			},
		)

	}

}
func (procon ProductController) DeleteProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		var err error
		id, err := strconv.Atoi(c.Param("id"))

		_, err = procon.Repo.Delete(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "success",
			},
		)
	}

}
func (procon ProductController) PutProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		PutProductReq := PutProductRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		err := c.Bind(&PutProductReq)

		newProduct := entities.Products{
			Name:        PutProductReq.Name,
			Description: PutProductReq.Description,
		}
		if id < 1 || err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		_, err = procon.Repo.Update(newProduct, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}