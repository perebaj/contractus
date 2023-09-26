package postgres

import "github.com/jmoiron/sqlx"

// TransactionStorage deal with the database layer for transactions.
type TransactionStorage struct {
	db *sqlx.DB
}

// NewTransactionStorage initialize a new transaction storage.
func NewTransactionStorage(db *sqlx.DB) *TransactionStorage {
	return &TransactionStorage{
		db: db,
	}
}

// SaveTransaction save a transaction in the database.
func (t TransactionStorage) SaveTransaction() error {
	return nil
}
