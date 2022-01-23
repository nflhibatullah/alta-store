package transaction

import (
	"altastore/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type Transaction interface {
	GetAll() ([]entities.Transaction, error)
	GetByUser(userId int) ([]entities.Transaction, error)
	GetByTransaction(userId int, transactionId int) (entities.Transaction, error)
	GetByTransactionAdmin(transactionId int) (entities.Transaction, error)
	GetByInvoiceId(invoiceId string) (entities.Transaction, error)
	Create(transaction entities.Transaction) (entities.Transaction, error)
	Update(transactionId int, transaction entities.Transaction) (entities.Transaction, error)

	StoreItemProduct(transactionId int, item entities.TransactionDetail) (entities.TransactionDetail, error)

	GetProductInCart(userId int, productId int) bool
	UpdateStockProduct(productId int, stock int) bool
	DeleteProductInCart(userId int, productId int) bool
}

func (tr *TransactionRepository) GetAll() ([]entities.Transaction, error) {

	var transactions []entities.Transaction

	if err := tr.db.Preload("TransactionDetails.Product.Category").Preload(clause.Associations).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) GetByUser(userId int) ([]entities.Transaction, error) {

	var transactions []entities.Transaction

	if err := tr.db.Preload("TransactionDetails.Product.Category").Preload(clause.Associations).Where(
		"user_id = ?", userId,
	).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) GetByTransaction(userId int, transactionId int) (entities.Transaction, error) {
	var transaction entities.Transaction

	err := tr.db.Preload("TransactionDetails.Product.Category").Preload(clause.Associations).Where(
		"user_id = ?", userId,
	).First(&transaction, transactionId).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) GetByTransactionAdmin(transactionId int) (entities.Transaction, error) {
	var transaction entities.Transaction

	err := tr.db.Preload("TransactionDetails.Product.Category").Preload(clause.Associations).First(
		&transaction, transactionId,
	).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) GetByInvoiceId(invoiceId string) (entities.Transaction, error) {
	var transaction entities.Transaction

	err := tr.db.Preload("TransactionDetails.Product.Category").Where(
		"invoice_id = ?", invoiceId,
	).First(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) Create(transaction entities.Transaction) (entities.Transaction, error) {
	if err := tr.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) Update(transactionId int, transaction entities.Transaction) (
	entities.Transaction, error,
) {
	var t entities.Transaction

	tr.db.First(&t, transactionId)

	if err := tr.db.Model(&t).Updates(transaction).Error; err != nil {
		return t, err
	}

	return t, nil
}

func (tr *TransactionRepository) StoreItemProduct(
	transactionId int, item entities.TransactionDetail,
) (entities.TransactionDetail, error) {
	if err := tr.db.Create(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func (tr *TransactionRepository) UpdateStockProduct(productId int, stock int) bool {
	var p entities.Product

	if err := tr.db.Model(&p).First(&p, productId).Update("Stock", stock).Error; err != nil {
		return false
	}

	return true
}

func (tr *TransactionRepository) GetProductInCart(userId int, productId int) bool {
	var cart entities.Cart

	if err := tr.db.Where("user_id = ? AND product_id = ?", userId, productId).First(&cart).Error; err != nil {
		return false
	}

	return true
}

func (tr *TransactionRepository) DeleteProductInCart(userId int, productId int) bool {
	var cart entities.Cart

	if err := tr.db.Where("user_id = ? AND product_id = ?", userId, productId).Delete(&cart).Error; err != nil {
		return false
	}

	return true
}
