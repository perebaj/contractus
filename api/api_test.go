package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/perebaj/contractus"
)

// TODO(JOJO) assert the body
func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	sendErr(w, http.StatusBadRequest, Error{Code: "bad_request", Msg: "Bad request"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TODO(JOJO) assert the body
func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	send(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{"Hello"})
	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestConvert(t *testing.T) {
	content := `12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS
	12021-12-03T11:46:02-03:00DOMINANDO INVESTIMENTOS       0000050000MARIA CANDIDA
	`
	transac, err := convert(content, "jj@example.com")
	if err != nil {
		t.Fatal(err)
	}

	want := []contractus.Transaction{
		{
			Email:              "jj@example.com",
			Type:               1,
			Date:               time.Date(2022, 01, 15, 22, 20, 30, 0, time.UTC),
			ProductDescription: "CURSO DE BEM-ESTAR",
			ProductPriceCents:  12750,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
		{
			Email:              "jj@example.com",
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
		assert(t, transac[0].Email, want[0].Email)
		assert(t, transac[0].Type, want[0].Type)
		assert(t, transac[0].Date.Format(time.RFC3339), want[0].Date.Format(time.RFC3339))
		assert(t, transac[0].ProductDescription, want[0].ProductDescription)
		assert(t, transac[0].ProductPriceCents, want[0].ProductPriceCents)
		assert(t, transac[0].SellerName, want[0].SellerName)
		assert(t, transac[0].SellerType, want[0].SellerType)
		assert(t, transac[0].Action, want[0].Action)

		assert(t, transac[1].Email, want[1].Email)
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
