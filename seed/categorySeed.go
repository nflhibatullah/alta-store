package seed

import (
	"altastore/entities"
	"gorm.io/gorm"
)

func CategorySeed(db *gorm.DB) {

	category1 := entities.Category{
		Name: "category 1",
	}
	category2 := entities.Category{
		Name: "category 2",
	}
	category3 := entities.Category{
		Name: "category 3",
	}

	db.Create(&category1)
	db.Create(&category2)
	db.Create(&category3)
}
