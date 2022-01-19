package product

import "altastore/entities"

type ProductInterface interface {
	GetAll() ([]entities.Product, error)
	Get(toDoId int) ([]entities.Product, error)
	Create(todo entities.Product) (entities.Product, error)
	Delete(toDoId int) (entities.Product, error)
	Update(newProduct entities.Product, toDoId int) ([]entities.Product, error)
}
