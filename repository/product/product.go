package product

import (
	"altastore/entities"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) GetAll() ([]entities.Product, error) {
	products := []entities.Product{}
	pr.db.Preload("Category").Find(&products)

	return products, nil
}

func (pr *ProductRepository) Get(productId int) ([]entities.Product, error) {
	product := []entities.Product{}
	pr.db.Where("id = ?", productId).Find(&product)

	return product, nil
}

func (pr *ProductRepository) Create(product entities.Product) (entities.Product, error) {
	pr.db.Save(&product)
	return product, nil
}

func (pr *ProductRepository) Delete(productId int) (entities.Product, error) {
	product := entities.Product{}
	pr.db.Find(&product, "id = ?", productId)
	pr.db.Delete(&product)
	return product, nil
}

func (pr *ProductRepository) Update(newProduct entities.Product, productId int) (entities.Product, error) {
	product := entities.Product{}

	pr.db.Find(&product, "id=?", productId)
	pr.db.Model(&product).Updates(newProduct)

	return product, nil
}
