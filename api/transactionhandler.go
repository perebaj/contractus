package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type transactionStorage interface {
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
}

func (s transactionHandler) balance(w http.ResponseWriter, _ *http.Request) {
	t := struct {
		ProducerID string `json:"producer_id"`
	}{ProducerID: "123"}

	send(w, http.StatusOK, t)
}
