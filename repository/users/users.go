package users

import (
	"altastore/entities"
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAll() ([]entities.User, error) {

	users := []entities.User{}
	ur.db.Find(&users, "role=?", "user")
	return users, nil
}

func (ur *UserRepository) Get(userId int) (entities.User, error) {
	user := entities.User{}
	err := ur.db.Find(&user, userId).Error
	log.Error(err)
	fmt.Println(err)
	return user, nil
}

func (ur *UserRepository) Create(newUser entities.User) (entities.User, error) {
	err := ur.db.Save(&newUser).Error
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}
func (ur *UserRepository) Update(updateUser entities.User, userId int) (entities.User, error) {
	User := entities.User{}
	err := ur.db.First(&User, "id=?", userId).Error
	ur.db.Model(&User).Updates(updateUser)

	if err != nil {
		return User, err
	}

	return updateUser, nil
}

func (ur *UserRepository) Delete(userId int) error {
	user := entities.User{}
	ur.db.Find(&user, "id=?", userId)
	ur.db.Delete(&user)
	return nil
}

func (ur *UserRepository) Login(email string) (entities.User, error) {
	var user entities.User
	var err = ur.db.First(&user, "email = ?", email).Error

	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) GetDeleteData(id int) (entities.User, error) {
	var user entities.User
	err := ur.db.First(&user, "id = ?", id).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
