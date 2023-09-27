package api

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/perebaj/contractus"
)

type mockTransactionStorage struct{}

func (m *mockTransactionStorage) SaveTransaction(_ context.Context, _ []contractus.Transaction) error {
	return nil
}

func (m *mockTransactionStorage) Balance(_ context.Context, _ string, _ string) (*contractus.BalanceResponse, error) {
	return nil, nil
}

func TestTransactionHandlerUpload(t *testing.T) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}

	fileContent := `12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS
	12021-12-03T11:46:02-03:00DOMINANDO INVESTIMENTOS       0000050000MARIA CANDIDA`
	_, err = io.WriteString(fw, fileContent)
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
	m := &mockTransactionStorage{}
	r := chi.NewRouter()
	RegisterHandler(r, m)

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestTransactionHandlerBalanceProducer(t *testing.T) {
	m := &mockTransactionStorage{}
	r := chi.NewRouter()
	RegisterHandler(r, m)

	req := httptest.NewRequest(http.MethodGet, "/balance/producer?name=JOSE", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestTransactionHandlerBalanceAffiliate(t *testing.T) {
	m := &mockTransactionStorage{}
	r := chi.NewRouter()
	RegisterHandler(r, m)

	req := httptest.NewRequest(http.MethodGet, "/balance/affiliate?name=JOSE", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}
