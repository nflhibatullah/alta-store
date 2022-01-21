package product

import "altastore/entities"

type ProductInterface interface {
	GetAll(offset, pageSize int) ([]entities.Product, error)
	Get(productId int) ([]entities.Product, error)
	Create(product entities.Product) (entities.Product, error)
	Delete(productId int) (entities.Product, error)
	Update(newProduct entities.Product, productId int) (entities.Product, error)
}
