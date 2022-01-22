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

func (pr *ProductRepository) GetAll(offset, pageSize int, search string) ([]entities.Product, error) {
	products := []entities.Product{}

	pr.db.Preload("Category").Offset(offset).Limit(pageSize).Where("name LIKE ?", "%"+search+"%").Find(&products)

	return products, nil
}

func (pr *ProductRepository) Get(productId int) (entities.Product, error) {
	product := entities.Product{}
	err := pr.db.Preload("Category").Where("id = ?", productId).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (pr *ProductRepository) Create(product entities.Product) (entities.Product, error) {
	pr.db.Save(&product)
	return product, nil
}

func (pr *ProductRepository) Delete(productId int) (entities.Product, error) {
	product := entities.Product{}
	err := pr.db.First(&product, "id = ?", productId).Error
	if err != nil {
		return product, err
	}
	pr.db.Delete(&product)
	return product, nil
}

func (pr *ProductRepository) Update(newProduct entities.Product, productId int) (entities.Product, error) {
	product := entities.Product{}

	err := pr.db.First(&product, "id = ?", productId).Error
	if err != nil {
		return product, err
	}
	pr.db.Model(&product).Updates(newProduct)

	return product, nil
}
