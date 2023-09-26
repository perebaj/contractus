// Package contractus implements the transaction struct and its methods.
package contractus

import (
	"fmt"
	"time"
)

// TransactionResponse have the fields that represent the API response.
type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
}

// Transaction have the fields that represent a single transaction.
type Transaction struct {
	Type               int       `json:"type" db:"type"`
	Date               time.Time `json:"date" db:"date"`
	ProductDescription string    `json:"product_description" db:"product_description"`
	ProductPriceCents  string    `json:"product_price_cents" db:"product_price_cents"`
	SellerName         string    `json:"seller_name" db:"seller_name"`
	SellerType         string    `json:"seller_type" db:"seller_type"`
}

func (t Transaction) typ() (string, error) {
	switch t.Type {
	case 1:
		return "venda produtor", nil
	case 2:
		return "venda afiliado", nil
	case 3:
		return "comissao paga", nil
	case 4:
		return "comissao recebida", nil
	default:
		return "", fmt.Errorf("invalid transaction type: %d", t.Type)
	}
}

func (t Transaction) sellerType() (string, error) {
	switch t.Type {
	case 1, 3:
		return "producer", nil
	case 2, 4:
		return "affiliate", nil
	default:
		return "", fmt.Errorf("invalid seller type: %d", t.Type)
	}
}
