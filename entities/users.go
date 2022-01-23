package entities

import (
	_ "github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	//USERID AUTO GENERATE
	ID       uint
	Email    string `gorm:"unique" `
	Role     string `gorm:"default:user"`
	Name     string 
	Password string 
	Transaction []Transaction
}
