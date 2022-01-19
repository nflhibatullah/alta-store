package product

import "altastore/entities"

type ProductInterface interface {
	GetAll() ([]entities.Products, error)
	Get(toDoId int) ([]entities.Products, error)
	Create(todo entities.Products) (entities.Products, error)
	Delete(toDoId int) (entities.Products, error)
	Update(newProduct entities.Products, toDoId int) ([]entities.Products, error)
}
