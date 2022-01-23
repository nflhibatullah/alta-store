package category

import "altastore/entities"

type CreateCategoryResponseFormat struct {
	Message string              `json:"message"`
	Data    []entities.Category `json:"data"`
}

type GetCategoryResponseFormat struct {
	Message interface{}         `json:"message"`
	Data    []entities.Category `json:"data"`
}
type GetAllCategoryResponseFormat struct {
	Message string              `json:"message"`
	Data    []entities.Category `json:"data"`
}
type DeleteCategoryResponseFormat struct {
	Message string `json:"message"`
}

type PutCategoryResponseFormat struct {
	Message string              `json:"message"`
	Data    []entities.Category `json:"data"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
