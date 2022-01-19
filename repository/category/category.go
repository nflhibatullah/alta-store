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

func (cr *CategoryRepository) Get(categoryId int) ([]entities.Category, error) {
	category := []entities.Category{}
	cr.db.Where("id = ?", categoryId).Find(&category)

	return category, nil
}

func (cr *CategoryRepository) Create(category entities.Category) (entities.Category, error) {
	cr.db.Save(&category)
	return category, nil
}

func (cr *CategoryRepository) Delete(categoryId int) (entities.Category, error) {
	category := entities.Category{}
	cr.db.Find(&category, "id = ?", categoryId)
	cr.db.Delete(&category)
	return category, nil
}

func (cr *CategoryRepository) Update(newCategory entities.Category, categoryId int) ([]entities.Category, error) {
	category := []entities.Category{}

	cr.db.Where("id = ?", categoryId).Find(&category).Save(
		map[string]interface{}{
			"name": newCategory.Name,
		},
	)

	return category, nil
}
