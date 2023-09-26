//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/birdie-ai/contractus"
	"github.com/birdie-ai/contractus/postgres"
	"github.com/jmoiron/sqlx"
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
	want := contractus.Transaction{
		Type:               1,
		Date:               time.Now().UTC(),
		ProductDescription: "Product description",
		ProductPriceCents:  "1000",
		SellerName:         "John Doe",
		SellerType:         "producer",
	}

	err := storage.SaveTransaction(&want)
	if err != nil {
		t.Fatalf("error saving transaction: %v", err)
	}

	var got contractus.Transaction
	err = db.Get(&got, "SELECT * FROM transactions LIMIT 1")
	if err != nil {
		t.Fatalf("error getting transaction: %v", err)
	}

	assert(t, got.Type, want.Type)
	assert(t, got.Date.Format(time.RFC3339), want.Date.Format(time.RFC3339))
	assert(t, got.ProductDescription, want.ProductDescription)
	assert(t, got.ProductPriceCents, want.ProductPriceCents)
	assert(t, got.SellerName, want.SellerName)
	assert(t, got.SellerType, want.SellerType)
}

func TestStorageTransactions(t *testing.T) {
	db := OpenDB(t)
	storage := postgres.NewStorage(db)

	want := contractus.Transaction{
		Type:               1,
		Date:               time.Now().UTC(),
		ProductDescription: "Product description",
		ProductPriceCents:  "1000",
		SellerName:         "John Doe",
		SellerType:         "producer",
	}

	want2 := contractus.Transaction{
		Type:               2,
		Date:               time.Now().UTC(),
		ProductDescription: "Product description 2",
		ProductPriceCents:  "2000",
		SellerName:         "John Doe 2",
		SellerType:         "affiliate",
	}

	err := storage.SaveTransaction(&want)
	if err != nil {
		t.Fatalf("error saving transaction 1: %v", err)
	}

	err = storage.SaveTransaction(&want2)
	if err != nil {
		t.Fatalf("error saving transaction 2: %v", err)
	}

	got, err := storage.Transactions()
	if err != nil {
		t.Fatalf("error getting transactions: %v", err)
	}

	// TODO(JOJO): Fix this undeterministic test.
	// For a while we don't have a consistent way to order the transactions, for this reason we choose to validate
	// the order based on the insertion order.
	// Reference: https://dba.stackexchange.com/questions/95822/does-postgres-preserve-insertion-order-of-records
	if got.Total == 2 {
		assert(t, got.Transactions[0].Type, want.Type)
		assert(t, got.Transactions[0].Date.Format(time.RFC3339), want.Date.Format(time.RFC3339))
		assert(t, got.Transactions[0].ProductDescription, want.ProductDescription)
		assert(t, got.Transactions[0].ProductPriceCents, want.ProductPriceCents)
		assert(t, got.Transactions[0].SellerName, want.SellerName)
		assert(t, got.Transactions[0].SellerType, want.SellerType)

		assert(t, got.Transactions[1].Type, want2.Type)
		assert(t, got.Transactions[1].Date.Format(time.RFC3339), want2.Date.Format(time.RFC3339))
		assert(t, got.Transactions[1].ProductDescription, want2.ProductDescription)
		assert(t, got.Transactions[1].ProductPriceCents, want2.ProductPriceCents)
		assert(t, got.Transactions[1].SellerName, want2.SellerName)
		assert(t, got.Transactions[1].SellerType, want2.SellerType)
	} else {
		t.Fatal("error getting transactions: expected 2 transactions")
	}
}

func assert(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
