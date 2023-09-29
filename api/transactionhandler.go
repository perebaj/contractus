package api

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-openapi/runtime/middleware"
	"github.com/perebaj/contractus"
	"github.com/perebaj/contractus/postgres"
)

type transactionStorage interface {
	SaveTransaction(ctx context.Context, t []contractus.Transaction) error
	Balance(ctx context.Context, sellerType, sellerName, email string) (*contractus.BalanceResponse, error)
	Transactions(ctx context.Context, email string) (contractus.TransactionResponse, error)
}

type transactionHandler struct {
	storage transactionStorage
}

// RegisterTransactionsHandler gather all the handlers for the API.
func RegisterTransactionsHandler(r chi.Router, storage transactionStorage) {
	h := transactionHandler{
		storage: storage,
	}

	const balanceProducer = "/balance/producer"
	r.Method(http.MethodGet, balanceProducer, http.HandlerFunc(h.producerBalance))

	const balanceAffiliate = "/balance/affiliate"
	r.Method(http.MethodGet, balanceAffiliate, http.HandlerFunc(h.affiliateBalance))

	const upload = "/upload"
	r.Method(http.MethodPost, upload, http.HandlerFunc(h.upload))

	const transactions = "/transactions"
	r.Method(http.MethodGet, transactions, http.HandlerFunc(h.transactions))

}

//go:embed docs/api.yml
var swagger embed.FS

// RegisterSwaggerHandler register the the swagger endpoint.
func RegisterSwaggerHandler(r chi.Router) {
	opts := middleware.SwaggerUIOpts{SpecURL: "docs/api.yml",
		Path:  "/docs",
		Title: "Contractus API",
	}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/docs/api.yml", http.FileServer(http.FS(swagger)))
}

func (s transactionHandler) producerBalance(w http.ResponseWriter, r *http.Request) {
	email, err := emailFromRequest(r)
	if err != nil {
		sendErr(w, http.StatusBadRequest, Error{"email_required", "email is required"})
		return
	}

	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		sendErr(w, http.StatusBadRequest, Error{"name_required", "name is required"})
		return
	}

	b, err := s.storage.Balance(r.Context(), "producer", name, email)
	if err != nil {
		if err == postgres.ErrSellerNotFound {
			sendErr(w, http.StatusNotFound, Error{"seller_not_found", "seller not found"})
			return
		}
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	send(w, http.StatusOK, b)
}

func (s transactionHandler) affiliateBalance(w http.ResponseWriter, r *http.Request) {
	email, err := emailFromRequest(r)
	if err != nil {
		sendErr(w, http.StatusBadRequest, Error{"email_required", "email is required"})
		return
	}
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		sendErr(w, http.StatusBadRequest, Error{"name_required", "name is required"})
		return
	}

	b, err := s.storage.Balance(r.Context(), "affiliate", name, email)
	if err != nil {
		if err == postgres.ErrSellerNotFound {
			sendErr(w, http.StatusNotFound, Error{"seller_not_found", "seller not found"})
			return
		}
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	send(w, http.StatusOK, b)
}

func (s transactionHandler) upload(w http.ResponseWriter, r *http.Request) {
	email, err := emailFromRequest(r)
	if err != nil {
		sendErr(w, http.StatusBadRequest, Error{"email_required", "email is required"})
		return
	}

	content, err := parseFile(r)
	if err != nil {
		slog.Error("Failed to parse file", "error", err)
		sendErr(w, http.StatusBadRequest, Error{"invalid_file", "invalid file"})
		return
	}

	transactions, err := convert(content, email)
	if err != nil {
		slog.Error("Failed to convert transactions", "error", err, "content", content)
		sendErr(w, http.StatusBadRequest, Error{"invalid_file", "invalid file"})
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

func (s transactionHandler) transactions(w http.ResponseWriter, r *http.Request) {
	email, err := emailFromRequest(r)
	if err != nil {
		sendErr(w, http.StatusBadRequest, Error{"email_required", "email is required"})
		return
	}

	t, err := s.storage.Transactions(r.Context(), email)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err)
		return
	}

	send(w, http.StatusOK, t)
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
func (t *Transaction) Convert(email string) (*contractus.Transaction, error) {
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
		Email:              email,
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

func emailFromRequest(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", fmt.Errorf("failed to get claims from context: %v", err)
	}
	return claims["email"].(string), nil
}
