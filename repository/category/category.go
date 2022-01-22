package category

import (
	"altastore/entities"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) GetAll() ([]entities.Category, error) {
	categorys := []entities.Category{}
	cr.db.Find(&categorys)

	return categorys, nil
}

func (cr *CategoryRepository) Get(categoryId int) (entities.Category, error) {
	category := entities.Category{}
	err := cr.db.Where("id = ?", categoryId).First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (cr *CategoryRepository) Create(category entities.Category) (entities.Category, error) {
	err := cr.db.Save(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (cr *CategoryRepository) Delete(categoryId int) (entities.Category, error) {
	category := entities.Category{}
	err := cr.db.First(&category, "id = ?", categoryId).Delete(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

func (cr *CategoryRepository) Update(newCategory entities.Category, categoryId int) (entities.Category, error) {
	category := entities.Category{}

	err := cr.db.First(&category, "id=?", categoryId).Error
	if err != nil {
		return category, err
	}
	cr.db.Model(&category).Updates(newCategory)

	return category, nil
}
