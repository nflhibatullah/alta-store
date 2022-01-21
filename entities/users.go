package entities

import (
	_ "github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//USERID AUTO GENERATE
	ID       uint
	Email    string `gorm:"uniqueIndex" `
	Role     string `gorm:"default:user"`
	Name     string `valid:"required"`
	Password string `valid:"required"`
	Transaction []Transaction
}
