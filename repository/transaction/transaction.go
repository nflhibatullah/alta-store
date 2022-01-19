package transaction

import (
	"altastore/entities"

	"gorm.io/gorm"
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
	Create(transaction entities.Transaction) (entities.Transaction, error)
	Update(transactionId int, transaction entities.Transaction) (entities.Transaction, error)

	StoreItemProduct(transactionId int, item entities.TransactionDetail) (entities.TransactionDetail, error)
}

func (tr *TransactionRepository) GetAll() ([]entities.Transaction, error) {

	return nil, nil
}

func (tr *TransactionRepository) GetByUser(userId int) ([]entities.Transaction, error) {

	return nil, nil
}

func (tr *TransactionRepository) GetByTransaction(userId int, transactionId int) (entities.Transaction, error) {

	return entities.Transaction{}, nil
}

func (tr *TransactionRepository) Create(transaction entities.Transaction) (entities.Transaction, error) {
	if err := tr.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) Update(transactionId int, transaction entities.Transaction) (entities.Transaction, error) {
	// var t entities.Transaction

	// tr.db.

	return entities.Transaction{}, nil
}

func (tr *TransactionRepository) StoreItemProduct(transactionId int, item entities.TransactionDetail) (entities.TransactionDetail, error) {
	if err := tr.db.Create(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}