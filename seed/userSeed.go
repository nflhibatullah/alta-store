package seed

import (
	"altastore/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSeed(db *gorm.DB) {
	password1, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 14)
	user1 := entities.User{
		Email:    "naufal@gmail.com",
		Role:     "admin",
		Name:     "Naufal",
		Password: string(password1),
	}
	user2 := entities.User{
		Email:    "furqon@gmail.com",
		Role:     "user",
		Name:     "furqon",
		Password: string(password1),
	}

	db.Create(&user1)
	db.Create(&user2)
}
