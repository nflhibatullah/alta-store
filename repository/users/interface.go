package users

import "altastore/entities"

type UserInterface interface {
	Login(name, password string) (entities.User, error)
	GetAll() ([]entities.User, error)
	Get(userId int) (entities.User, error)
	Create(newUser entities.User) (entities.User, error)
	Update(updateUser entities.User, userId int) (entities.User, error)
	Delete(userId int) (entities.User, error)
}
