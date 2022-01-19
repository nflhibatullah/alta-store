package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	//USERID AUTO GENERATE
	ID       uint
	Email    string
	Role     string `gorm:"default:user"`
	Name     string
	Password string
}
