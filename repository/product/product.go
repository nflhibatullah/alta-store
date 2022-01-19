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

func (pr *ProductRepository) GetAll() ([]entities.Products, error) {
	products := []entities.Products{}
	pr.db.Find(&products)

	return products, nil
}

func (pr *ProductRepository) Get(productId int) ([]entities.Products, error) {
	product := []entities.Products{}
	pr.db.Where("id = ?", productId).Find(&product)

	return product, nil
}

func (pr *ProductRepository) Create(product entities.Products) (entities.Products, error) {
	pr.db.Save(&product)
	return product, nil
}

func (pr *ProductRepository) Delete(productId int) error {
	product := entities.Products{}
	pr.db.Find(&product, "id = ?", productId)
	pr.db.Delete(&product)
	return nil
}

func (pr *ProductRepository) Update(newProducts entities.Products, productId int) ([]entities.Products, error) {
	product := []entities.Products{}

	pr.db.Where("id = ?", productId).Find(&product).Save(
		map[string]interface{}{
			"name": newProducts.Name, "price": newProducts.Price, "description": newProducts.Description,
		},
	)

	return product, nil
}
