package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/perebaj/contractus"
)

type mockTransactionStorage struct{}

func (m *mockTransactionStorage) SaveTransaction(_ context.Context, _ []contractus.Transaction) error {
	return nil
}

func (m *mockTransactionStorage) Balance(_ context.Context, _ string, _ string, _ string) (*contractus.BalanceResponse, error) {
	return nil, nil
}

func (m *mockTransactionStorage) Transactions(_ context.Context, _ string) (contractus.TransactionResponse, error) {
	return contractus.TransactionResponse{}, nil
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

	token, authConfig := authTest(t)
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodPost, "/upload", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var response struct {
		Msg string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	assert(t, response.Msg, "file uploaded successfully")

}

func TestTransactionHandlerUpload_Unauthorized(t *testing.T) {
	_, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestTransactionHandlerBalanceProducer(t *testing.T) {
	token, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodGet, "/balance/producer?name=JOSE%20CARLOS", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d, Response Body %s", http.StatusOK, resp.Code, resp.Body.String())
	}
}

func TestTransactionHandlerBalanceProducer_Unauthorized(t *testing.T) {
	_, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodGet, "/balance/producer?name=JOSE%20CARLOS", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestTransactionHandlerBalanceAffiliate(t *testing.T) {
	token, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodGet, "/balance/affiliate?name=JOSE%20CARLOS", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestTransactionHandlerTransactions(t *testing.T) {
	token, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}

func TestTransactionHandlerTransactions_Unauthorized(t *testing.T) {
	_, authConfig := authTest(t)
	m := &mockTransactionStorage{}
	r := Router(authConfig, m)

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
	}
}

func TestSwaggerHandler(t *testing.T) {
	r := chi.NewRouter()
	RegisterSwaggerHandler(r)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/docs/api.yml", nil)
	resp = httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.Code)
	}
}

// authTest is a helper function to generate a JWT token and Auth struct for testing.
func authTest(t *testing.T) (token string, authConfig Auth) {
	t.Helper()
	secretKey := "secret" //Fake secret Key
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	_, token, _ = tokenAuth.Encode(map[string]interface{}{"email": "jj@example.com"})

	authConfig = Auth{
		JWTSecretKey: secretKey,
	}

	return token, authConfig
}
