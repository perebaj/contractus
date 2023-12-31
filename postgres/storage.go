package postgres

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/perebaj/contractus"
)

var (
	// ErrSellerNotFound is returned when the seller is not found in the database.
	ErrSellerNotFound = errors.New("seller not found")
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
func (s Storage) SaveTransaction(ctx context.Context, t []contractus.Transaction) error {
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO transactions (email, type, date, product_description, product_price_cents, seller_name, seller_type, action)
		VALUES (:email, :type, :date, :product_description, :product_price_cents, :seller_name, :seller_type, :action)
	`, t)

	return err
}

// Transactions is responsible for returning all the transactions from the database.
// TODO(JOJO): Have a way to paginate the transactions.
func (s Storage) Transactions(ctx context.Context, email string) (contractus.TransactionResponse, error) {
	var transactions []contractus.Transaction

	err := s.db.SelectContext(ctx, &transactions, `
		SELECT email, type, date, product_description, product_price_cents, seller_name, seller_type, action
		FROM transactions WHERE email = $1
	`, email)
	if err != nil {
		return contractus.TransactionResponse{}, err
	}

	return contractus.TransactionResponse{
		Transactions: transactions,
		Total:        len(transactions),
	}, nil
}

// Balance is responsible for return the balance of a seller.
func (s Storage) Balance(ctx context.Context, sellerType, sellerName, email string) (*contractus.BalanceResponse, error) {
	var transactions []contractus.Transaction

	err := s.db.SelectContext(ctx, &transactions, `
		SELECT type, date, product_description, product_price_cents, seller_name, seller_type
		FROM transactions
		WHERE seller_type = $1 AND seller_name = $2 AND email = $3
	`, sellerType, sellerName, email)

	if len(transactions) == 0 {
		return nil, ErrSellerNotFound
	}
	if err != nil {
		return nil, err
	}

	var balance int64
	for _, t := range transactions {
		if t.Type != 3 {
			balance += t.ProductPriceCents
		} else {
			balance -= t.ProductPriceCents
		}
	}

	return &contractus.BalanceResponse{
		Balance:    balance,
		SellerName: sellerName,
	}, nil
}
