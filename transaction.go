// Package contractus implements the transaction struct and its methods.
package contractus

import "fmt"

// Transaction have the fields that represent a single transaction.
type Transaction struct {
	Type               int    `json:"type"`
	Date               string `json:"date"`
	ProductDescription string `json:"product_description"`
	ProductPrice       string `json:"product_price"`
	SellerName         string `json:"seller_name"`
	SellerType         string `json:"seller_type"`
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
