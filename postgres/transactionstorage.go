package postgres

import (
	"github.com/birdie-ai/contractus"
	"github.com/jmoiron/sqlx"
)

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

// Save is responsible for saving a transaction into the database.
func (s TransactionStorage) Save(t *contractus.Transaction) error {
	_, err := s.db.Exec(`
		INSERT INTO transactions (type, date, product_description, product_price_cents, seller_name, seller_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, t.Type, t.Date, t.ProductDescription, t.ProductPriceCents, t.SellerName, t.SellerType)

	return err
}
