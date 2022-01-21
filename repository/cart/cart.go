package cart

import (
	"altastore/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

type Cart interface {
	GetAll(userId int) ([]entities.Cart, error)
	Create(entities.Cart) (entities.Cart, error)
	Update(entities.Cart) (entities.Cart, error)
	Delete(userId int, productId int) (entities.Cart, error)
	CheckCart(entities.Cart) (entities.Cart, error)
}

func (cr *CartRepository) GetAll(userId int) ([]entities.Cart, error) {
	var carts []entities.Cart

	if err := cr.db.Preload("Product.Category").Preload(clause.Associations).Where("user_id = ?", userId).Find(&carts).Error; err != nil {
		return nil, err
	}

	return carts, nil
}

func (cr *CartRepository) CheckCart(c entities.Cart) (entities.Cart, error) {
	var cart entities.Cart

	if err := cr.db.Where("user_id = ? AND product_id = ?", c.UserID, c.ProductID).First(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (cr *CartRepository) Create(cart entities.Cart) (entities.Cart, error) {
	if err := cr.db.Create(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (cr *CartRepository) Update(cart entities.Cart) (entities.Cart, error) {
	var c entities.Cart

	if err := cr.db.Model(&c).Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&c).Updates(cart).Error; err != nil {
		return c, err
	}

	return c, nil
}

func (cr *CartRepository) Delete(userId int, productId int) (entities.Cart, error) {
	var cart entities.Cart

	err := cr.db.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&cart).Error
	if err != nil {
		return cart, err
	}
	return cart, nil
}