package postgres

import (
	"github.com/birdie-ai/contractus"
	"github.com/jmoiron/sqlx"
)

// Storage deal with the database layer for transactions.
type Storage struct {
	db *sqlx.DB
}

// NewStorage initialize a new transaction storage.
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// SaveTransaction is responsible for saving a transaction into the database.
func (s Storage) SaveTransaction(t *contractus.Transaction) error {
	_, err := s.db.Exec(`
		INSERT INTO transactions (type, date, product_description, product_price_cents, seller_name, seller_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, t.Type, t.Date, t.ProductDescription, t.ProductPriceCents, t.SellerName, t.SellerType)

	return err
}

// TODO(JOJO): Have a way to paginate the transactions.
// Transactions is responsible for returning all the transactions from the database.
func (s Storage) Transactions() (contractus.TransactionResponse, error) {
	var transactions []contractus.Transaction

	err := s.db.Select(&transactions, `
		SELECT type, date, product_description, product_price_cents, seller_name, seller_type
		FROM transactions
	`)
	if err != nil {
		return contractus.TransactionResponse{}, err
	}

	return contractus.TransactionResponse{
		Transactions: transactions,
		Total:        len(transactions),
	}, nil
}
