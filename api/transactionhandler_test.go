package api

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/perebaj/contractus"
)

type mockTransactionStorage struct{}

func (m *mockTransactionStorage) SaveTransaction(_ context.Context, _ []contractus.Transaction) error {
	return nil
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

func TestConvert(t *testing.T) {
	content := `12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS
	12021-12-03T11:46:02-03:00DOMINANDO INVESTIMENTOS       0000050000MARIA CANDIDA
	`
	transac, err := convert(content)
	if err != nil {
		t.Fatal(err)
	}

	want := []contractus.Transaction{
		{
			Type:               1,
			Date:               time.Date(2022, 01, 15, 22, 20, 30, 0, time.UTC),
			ProductDescription: "CURSO DE BEM-ESTAR",
			ProductPriceCents:  12750,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
		{
			Type:               1,
			Date:               time.Date(2021, 12, 03, 14, 46, 02, 0, time.UTC),
			ProductDescription: "DOMINANDO INVESTIMENTOS",
			ProductPriceCents:  50000,
			SellerName:         "MARIA CANDIDA",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
	}

	if len(transac) == len(want) {
		assert(t, transac[0].Type, want[0].Type)
		assert(t, transac[0].Date.Format(time.RFC3339), want[0].Date.Format(time.RFC3339))
		assert(t, transac[0].ProductDescription, want[0].ProductDescription)
		assert(t, transac[0].ProductPriceCents, want[0].ProductPriceCents)
		assert(t, transac[0].SellerName, want[0].SellerName)
		assert(t, transac[0].SellerType, want[0].SellerType)
		assert(t, transac[0].Action, want[0].Action)

		assert(t, transac[1].Type, want[1].Type)
		assert(t, transac[1].Date.Format(time.RFC3339), want[1].Date.Format(time.RFC3339))
		assert(t, transac[1].ProductDescription, want[1].ProductDescription)
		assert(t, transac[1].ProductPriceCents, want[1].ProductPriceCents)
		assert(t, transac[1].SellerName, want[1].SellerName)
		assert(t, transac[1].SellerType, want[1].SellerType)
		assert(t, transac[1].Action, want[1].Action)
	} else {
		t.Fatalf("expected %d transactions, got %d", len(want), len(transac))
	}
}

func assert(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
