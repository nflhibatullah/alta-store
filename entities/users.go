package entities

import "gorm.io/gorm"
import _ "github.com/asaskevich/govalidator"

type User struct {
	gorm.Model
	//USERID AUTO GENERATE
	ID       uint
	Email    string `gorm:"uniqueIndex" `
	Role     string `gorm:"default:user"`
	Name     string `valid:"required"`
	Password string `valid:"required"`
}
