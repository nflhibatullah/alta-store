package product

import "altastore/entities"

type CreateProductResponseFormat struct {
	Message string             `json:"message"`
	Data    []entities.Product `json:"data"`
}

type GetProductResponseFormat struct {
	Message interface{}        `json:"message"`
	Data    []entities.Product `json:"data"`
}
type GetAllProductResponseFormat struct {
	Message string             `json:"message"`
	Data    []entities.Product `json:"data"`
}
type DeleteProductResponseFormat struct {
	Message string `json:"message"`
}

type PutProductResponseFormat struct {
	Message string             `json:"message"`
	Data    []entities.Product `json:"data"`
}
