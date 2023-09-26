package postgres

import "github.com/jmoiron/sqlx"

type TransactionStorage struct {
	db *sqlx.DB
}

func NewTransactionStorage(db *sqlx.DB) *TransactionStorage {
	return &TransactionStorage{
		db: db,
	}
}

func (t TransactionStorage) SaveTransaction() error {
	return nil
}
