package utils

import (
	"altastore/configs"
	"altastore/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(config *configs.AppConfig) *gorm.DB {

	// conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
	// 	config.Database.Username,
	// 	config.Database.Password,
	// 	config.Database.Host,
	// 	config.Database.Port,
	// 	config.Database.Name,
	// )

	conn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Name + "?parseTime=true&loc=Asia%2FJakarta&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func InitialMigrate(db *gorm.DB)  {
	db.AutoMigrate(&entities.Transaction{})
	db.AutoMigrate(&entities.TransactionDetail{})
}