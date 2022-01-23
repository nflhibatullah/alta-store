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

type ProductResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
