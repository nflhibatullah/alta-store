package seed

import (
	"altastore/entities"
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

func ProductSeed(db *gorm.DB) {
	for i := 1; i <= 100; i++ {
		product := entities.Product{
			Name:        "Product " + fmt.Sprint(i),
			Price:       50000,
			Stock:       50,
			Description: "deskripsi product " + fmt.Sprint(i),
			CategoryID:  uint(rand.Intn(4-1) + 1),
		}
		db.Create(&product)
	}
}
