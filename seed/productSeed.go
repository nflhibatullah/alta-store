package seed

import (
	"altastore/entities"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
)

func ProductSeed(db *gorm.DB) {
	for i := 0; i < 100; i++ {
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
