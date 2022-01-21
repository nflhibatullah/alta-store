package utils

import (
	"altastore/configs"
	"altastore/entities"
	"altastore/seed"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(config *configs.AppConfig) *gorm.DB {

	conn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Name + "?parseTime=true&loc=Asia%2FJakarta&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func InitialMigrate(db *gorm.DB) {
	db.Migrator().DropTable(&entities.TransactionDetail{})
	db.Migrator().DropTable(&entities.Transaction{})
	db.Migrator().DropTable(&entities.Cart{})
	db.Migrator().DropTable(&entities.Category{})
	db.Migrator().DropTable(&entities.Product{})
	db.Migrator().DropTable(&entities.User{})

	db.AutoMigrate(&entities.Cart{})
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Category{})
	db.AutoMigrate(&entities.Product{})
	db.AutoMigrate(&entities.Cart{})
	db.AutoMigrate(&entities.Transaction{})
	db.AutoMigrate(&entities.TransactionDetail{})

	if configs.Mode == "development" {
		seed.UserSeed(db)
		seed.CategorySeed(db)
		seed.ProductSeed(db)
	}

}
