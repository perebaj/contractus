//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/perebaj/contractus"
	"github.com/perebaj/contractus/postgres"
)

// OpenDB create a new database for testing and return a connection to it.
func OpenDB(t *testing.T) *sqlx.DB {
	t.Helper()

	cfg := postgres.Config{
		URL:             os.Getenv("CONTRACTUS_POSTGRES_URL"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		t.Fatalf("error connecting to Postgres: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		t.Fatalf("error pinging postgres: %v", err)
	}

	// create a new database with random suffix
	postgresURL, err := url.Parse(cfg.URL)
	if err != nil {
		t.Fatalf("error parsing Postgres connection URL: %v", err)
	}
	database := strings.TrimLeft(postgresURL.Path, "/")

	randSuffix := fmt.Sprintf("%x", time.Now().UnixNano())

	database = fmt.Sprintf("%s-%x", database, randSuffix)
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, database))
	if err != nil {
		t.Fatalf("error creating database for test: %v", err)
	}

	postgresURL.Path = "/" + database
	cfg.URL = postgresURL.String()
	testDB, err := postgres.OpenDB(cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// after run the tests, drop the database
	t.Cleanup(func() {
		testDB.Close()
		defer db.Close()
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s" WITH (FORCE);`, database))
		if err != nil {
			t.Fatalf("error dropping database for test: %v", err)
		}
	})

	return testDB
}

func TestStorageSaveTransaction(t *testing.T) {
	db := OpenDB(t)
	storage := postgres.NewStorage(db)
	ctx := context.Background()
	want := []contractus.Transaction{
		{
			Type:               1,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description",
			ProductPriceCents:  1000,
			SellerName:         "John Doe",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
	}

	err := storage.SaveTransaction(ctx, want)
	if err != nil {
		t.Fatalf("error saving transaction: %v", err)
	}

	var got contractus.Transaction
	err = db.Get(&got, "SELECT * FROM transactions LIMIT 1")
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}

	assert(t, got.Type, want[0].Type)
	assert(t, got.Date.Format(time.RFC3339), want[0].Date.Format(time.RFC3339))
	assert(t, got.ProductDescription, want[0].ProductDescription)
	assert(t, got.ProductPriceCents, want[0].ProductPriceCents)
	assert(t, got.SellerName, want[0].SellerName)
	assert(t, got.SellerType, want[0].SellerType)
}

func TestStorageTransactions(t *testing.T) {
	db := OpenDB(t)
	storage := postgres.NewStorage(db)
	ctx := context.Background()

	want := []contractus.Transaction{
		{
			Type:               1,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description",
			ProductPriceCents:  1000,
			SellerName:         "John Doe",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
		{
			Type:               2,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description 2",
			ProductPriceCents:  2000,
			SellerName:         "John Doe 2",
			SellerType:         "affiliate",
			Action:             "venda afiliado",
		},
	}

	err := storage.SaveTransaction(ctx, want)
	if err != nil {
		t.Fatalf("error saving transaction 1: %v", err)
	}

	got, err := storage.Transactions(ctx)
	if err != nil {
		t.Fatalf("error getting transactions: %v", err)
	}

	// TODO(JOJO): Fix this undeterministic test.
	// For a while we don't have a consistent way to order the transactions, for this reason we choose to validate
	// the order based on the insertion order.
	// Reference: https://dba.stackexchange.com/questions/95822/does-postgres-preserve-insertion-order-of-records
	if got.Total == 2 {
		assert(t, got.Transactions[0].Type, want[0].Type)
		assert(t, got.Transactions[0].Date.Format(time.RFC3339), want[0].Date.Format(time.RFC3339))
		assert(t, got.Transactions[0].ProductDescription, want[0].ProductDescription)
		assert(t, got.Transactions[0].ProductPriceCents, want[0].ProductPriceCents)
		assert(t, got.Transactions[0].SellerName, want[0].SellerName)
		assert(t, got.Transactions[0].SellerType, want[0].SellerType)
		assert(t, got.Transactions[0].Action, want[0].Action)

		assert(t, got.Transactions[1].Type, want[1].Type)
		assert(t, got.Transactions[1].Date.Format(time.RFC3339), want[1].Date.Format(time.RFC3339))
		assert(t, got.Transactions[1].ProductDescription, want[1].ProductDescription)
		assert(t, got.Transactions[1].ProductPriceCents, want[1].ProductPriceCents)
		assert(t, got.Transactions[1].SellerName, want[1].SellerName)
		assert(t, got.Transactions[1].SellerType, want[1].SellerType)
		assert(t, got.Transactions[1].Action, want[1].Action)
	} else {
		t.Fatal("error getting transactions: expected 2 transactions")
	}
}

func TestStorageBalance(t *testing.T) {
	db := OpenDB(t)
	storage := postgres.NewStorage(db)
	ctx := context.Background()

	transactions1 := []contractus.Transaction{
		{
			Type:               1,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description",
			ProductPriceCents:  12750,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
		{
			Type:               3,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description 2",
			ProductPriceCents:  4500,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "comissao paga",
		},
		{
			Type:               1,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description 3",
			ProductPriceCents:  12750,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
	}

	err := storage.SaveTransaction(ctx, transactions1)
	if err != nil {
		t.Fatalf("error saving transactions 1: %v", err)
	}

	got, err := storage.Balance(ctx, "producer", "JOSE CARLOS")
	if err != nil {
		t.Fatalf("error getting balance: %v", err)
	}

	assert(t, got.Balance, int64(21000))
	assert(t, got.SellerName, "JOSE CARLOS")

	transactions2 := []contractus.Transaction{
		{
			Type:               2,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description",
			ProductPriceCents:  155000,
			SellerName:         "CARLOS BATISTA",
			SellerType:         "affiliate",
			Action:             "venda afiliado",
		},
		{
			Type:               4,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description 2",
			ProductPriceCents:  50000,
			SellerName:         "CARLOS BATISTA",
			SellerType:         "affiliate",
			Action:             "comissao recebida",
		},
	}

	err = storage.SaveTransaction(ctx, transactions2)
	if err != nil {
		t.Fatalf("error saving transactions 2: %v", err)
	}

	got, err = storage.Balance(ctx, "affiliate", "CARLOS BATISTA")
	if err != nil {
		t.Fatalf("error getting balance: %v", err)
	}

	assert(t, got.Balance, int64(205000))
	assert(t, got.SellerName, "CARLOS BATISTA")
}

func TestStorageBalance_NotFound(t *testing.T) {
	db := OpenDB(t)
	storage := postgres.NewStorage(db)
	ctx := context.Background()

	transactions1 := []contractus.Transaction{
		{
			Type:               1,
			Date:               time.Now().UTC(),
			ProductDescription: "Product description",
			ProductPriceCents:  12750,
			SellerName:         "JOSE CARLOS",
			SellerType:         "producer",
			Action:             "venda produtor",
		},
	}

	err := storage.SaveTransaction(ctx, transactions1)
	if err != nil {
		t.Fatalf("error saving transactions 1: %v", err)
	}

	got, err := storage.Balance(ctx, "producer", "INVALID NAME")
	if !errors.Is(err, postgres.ErrSellerNotFound) {
		t.Fatalf("error getting balance: %v", err)
	}

	assert(t, got, (*contractus.BalanceResponse)(nil))
}

func assert(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
