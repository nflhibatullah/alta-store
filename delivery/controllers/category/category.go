package category

import (
	"altastore/delivery/common"
	"altastore/entities"
	"altastore/repository/category"
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
		if err := c.Bind(&newCategoryReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		newCategory := entities.Category{
			Name: newCategoryReq.Name,
		}

		_, err := catcon.Repo.Create(newCategory)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}
func (catcon CategoryController) GetAllCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {

		category, err := catcon.Repo.GetAll()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "success",
				"data":    category,
			},
		)
	}

}
func (catcon CategoryController) GetCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		category, _ := catcon.Repo.Get(id)

		if len(category) == 0 {
			return c.JSON(
				http.StatusNotFound, map[string]interface{}{
					"message": "Category not found",
				},
			)
		}

		return c.JSON(
			http.StatusOK, map[string]interface{}{
				"message": "succes",
				"data":    category,
			},
		)

	}

}
func (catcon CategoryController) DeleteCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		var err error
		id, err := strconv.Atoi(c.Param("id"))

		_, err = catcon.Repo.Delete(id)
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
func (catcon CategoryController) PutCategoryCtrl() echo.HandlerFunc {

	return func(c echo.Context) error {
		PutCategoryReq := PutCategoryRequestFormat{}
		id, _ := strconv.Atoi(c.Param("id"))
		err := c.Bind(&PutCategoryReq)

		newCategory := entities.Category{
			Name: PutCategoryReq.Name,
		}
		if id < 1 || err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		result, _ := catcon.Repo.Update(newCategory, id)
		if result.ID == 0 {
			return c.JSON(http.StatusInternalServerError, common.NewInternalServerErrorResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}

}
