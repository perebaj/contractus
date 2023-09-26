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

func RegisterHandler(r chi.Router, storage transactionStorage) {
	h := transactionHandler{
		storage: storage,
	}

	const balanceProducer = "/balance/producer"
	const balanceAffiliate = "/balance/affiliate"

	r.Method(http.MethodGet, balanceProducer, http.HandlerFunc(h.balance))
}

func (s transactionHandler) balance(w http.ResponseWriter, r *http.Request) {
	t := struct {
		ProducerID string `json:"producer_id"`
	}{ProducerID: "123"}

	send(w, http.StatusOK, t)
}
