package product

import (
	"altastore/delivery/common"
	"altastore/entities"
	"altastore/repository/product"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

		newProduct := entities.Product{
			Name:        newProductReq.Name,
			Price:       newProductReq.Price,
			Stock:       newProductReq.Stock,
			Description: newProductReq.Description,
			CategoryID:  newProductReq.CategoryID,
		}

		product, err := procon.Repo.Create(newProduct)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(product))
	}

}
func (procon ProductController) GetAllProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		perpage, _ := strconv.Atoi(c.QueryParam("perpage"))
		if page == 0 {
			page = 1
		}
		if perpage == 0 {
			perpage = 10
		}
		offset := (page - 1) * perpage
		product, _ := procon.Repo.GetAll(offset, perpage)

		if len(product) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		return c.JSON(http.StatusOK, common.PaginationResponse(page, perpage, product))
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
		id, _ := strconv.Atoi(c.Param("id"))

		_, err = procon.Repo.Delete(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Produk tidak ditemukan"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse("Berhasil menghapus produk"))
	}

}
func (procon ProductController) PutProductCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		PutProductReq := PutProductRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		err := c.Bind(&PutProductReq)
		if err != nil {
			return err
		}

		newProduct := entities.Product{
			Name:        PutProductReq.Name,
			Stock:       PutProductReq.Stock,
			Price:       PutProductReq.Price,
			Description: PutProductReq.Description,
			CategoryID:  PutProductReq.CategoryID,
		}

		result, err := procon.Repo.Update(newProduct, id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Ada kesalahan dalam update"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(result))
	}

}
