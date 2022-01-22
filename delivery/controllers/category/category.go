package category

import (
	"altastore/delivery/common"
	"altastore/entities"
	"altastore/repository/category"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CategoryController struct {
	Repo category.CategoryInterface
}

func NewCategoryControllers(catrep category.CategoryInterface) *CategoryController {
	return &CategoryController{Repo: catrep}
}

//CreateCategory

func (catcon CategoryController) PostCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		newCategoryReq := CreateCategoryRequestFormat{}
		c.Bind(&newCategoryReq)
		if err := c.Validate(newCategoryReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Kesalahan dalam input"))
		}

		newCategory := entities.Category{
			Name: newCategoryReq.Name,
		}

		category, err := catcon.Repo.Create(newCategory)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Kategori yang diinput sudah ada"))
		}

		data := CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(data))
	}

}
func (catcon CategoryController) GetAllCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		category, _ := catcon.Repo.GetAll()

		if len(category) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		data := []CategoryResponse{}
		for _, item := range category {
			data = append(
				data, CategoryResponse{
					ID:   item.ID,
					Name: item.Name,
				},
			)
		}

		return c.JSON(
			http.StatusOK, common.SuccessResponse(data),
		)
	}

}
func (catcon CategoryController) GetCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		category, err := catcon.Repo.Get(id)
		fmt.Println(err)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		data := CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(data))

	}

}
func (catcon CategoryController) DeleteCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		_, err := catcon.Repo.Delete(id)

		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Kategori tidak ditemukan"))
		}

		return c.JSON(
			http.StatusOK, common.SuccessResponse("Berhasil menghapus category"),
		)
	}

}
func (catcon CategoryController) PutCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		PutCategoryReq := PutCategoryRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		c.Bind(&PutCategoryReq)
		if err := c.Validate(&PutCategoryReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Kesalahan dalam input"))
		}
		newCategory := entities.Category{
			Name: PutCategoryReq.Name,
		}

		result, err := catcon.Repo.Update(newCategory, id)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Kategori not found"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(result))
	}

}
