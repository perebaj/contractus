package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/perebaj/contractus"
)

type transactionStorage interface {
	SaveTransaction(ctx context.Context, t []contractus.Transaction) error
}

type transactionHandler struct {
	storage transactionStorage
}

// RegisterHandler gather all the handlers for the API.
func RegisterHandler(r chi.Router, storage transactionStorage) {
	h := transactionHandler{
		storage: storage,
	}

	const balanceProducer = "/balance/producer"
	r.Method(http.MethodGet, balanceProducer, http.HandlerFunc(h.balance))

	const upload = "/upload"
	r.Method(http.MethodPost, upload, http.HandlerFunc(h.upload))
}

func (s transactionHandler) balance(w http.ResponseWriter, _ *http.Request) {
	t := struct {
		ProducerID string `json:"producer_id"`
	}{ProducerID: "123"}

	send(w, http.StatusOK, t)
}

func (s transactionHandler) upload(w http.ResponseWriter, r *http.Request) {
	content, err := parseFile(r)
	if err != nil {
		slog.Error("Failed to parse file", "error", err)
		sendErr(w, http.StatusBadRequest, err)
		return
	}

	transactions, err := convert(content)
	if err != nil {
		slog.Error("Failed to convert transactions", "error", err, "content", content)
		sendErr(w, http.StatusBadRequest, err)
		return
	}
	err = s.storage.SaveTransaction(r.Context(), transactions)
	if err != nil {
		slog.Error("Failed to save transactions", "error", err, "transactions", transactions)
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	send(w, http.StatusOK, nil)
}

// Transaction represents the raw transaction from the file.
type Transaction struct {
	Type               string `json:"type"`
	Date               string `json:"date"`
	ProductDescription string `json:"product_description"`
	ProductPriceCents  string `json:"product_price_cents"`
	SellerName         string `json:"seller_name"`
}

// Convert transform the raw transaction to the business Transaction structure.
// TODO(JOJO) Join errors in one, and return all the errors.
func (t *Transaction) Convert() (*contractus.Transaction, error) {
	typeInt, err := strconv.Atoi(t.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert type: %v", err)
	}
	productPriceCentsInt, err := strconv.Atoi(t.ProductPriceCents)
	if err != nil {
		return nil, fmt.Errorf("failed to convert product price cents: %v", err)
	}
	// To play around timezone offset the format should be set up as follows:
	// https://pkg.go.dev/time#pkg-constants
	dateTime, err := time.Parse("2006-01-02T15:04:05-07:00", t.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to convert date: %v", err)
	}

	prodDesc := strings.TrimSpace(t.ProductDescription)

	sellerName := strings.Replace(t.SellerName, "\n", "", -1)

	transac := &contractus.Transaction{
		Type:               typeInt,
		Date:               dateTime.UTC(),
		ProductDescription: prodDesc,
		ProductPriceCents:  int64(productPriceCentsInt),
		SellerName:         sellerName,
	}

	sellerType, err := transac.ConvertSellerType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert seller type: %v", err)
	}

	transacAction, err := transac.ConvertType()
	if err != nil {
		return nil, fmt.Errorf("failed to convert type: %v", err)
	}

	transac.SellerType = sellerType
	transac.Action = transacAction

	return transac, nil
}
