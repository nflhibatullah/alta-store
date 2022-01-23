package users

import "altastore/entities"

type UserInterface interface {
	Login(email string) (entities.User, error)
	GetAll() ([]entities.User, error)
	Get(userId int) (entities.User, error)
	Create(newUser entities.User) (entities.User, error)
	Update(updateUser entities.User, userId int) (entities.User, error)
	Delete(userId int) error
	GetDeleteData(id int) (entities.User, error)
}
