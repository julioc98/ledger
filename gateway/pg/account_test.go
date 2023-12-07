package pg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julioc98/ledger/domain/entities"
	"github.com/julioc98/ledger/gateway/pg"
	"github.com/stretchr/testify/assert"
)

const (
	testDatabaseURL = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
)

func createTestDBConn() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), testDatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, nil
}

func setupTestDatabase(t *testing.T) (*pg.AccountPgxRepository, *pgx.Conn) {
	t.Helper()

	conn, err := createTestDBConn()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	repo := pg.NewAccountPgxRepository(conn)

	return repo, conn
}

func tearDownTestDatabase(conn *pgx.Conn) {
	conn.Exec(context.Background(), "TRUNCATE TABLE entries")
	conn.Close(context.Background())
}

func TestBalance(t *testing.T) {
	repo, conn := setupTestDatabase(t)
	defer tearDownTestDatabase(conn)

	t.Run("AccountWithBalance", func(t *testing.T) {
		// Insert test entries into the database
		_, err := conn.Exec(context.Background(),
			"INSERT INTO entries (id, account, direction, amount, created_at) VALUES (1, 'account1', 'debit', 100, NOW())",
		)
		assert.NoError(t, err)

		_, err = conn.Exec(context.Background(),
			"INSERT INTO entries (id, account, direction, amount, created_at) VALUES (2, 'account1', 'credit', 50, NOW())",
		)
		assert.NoError(t, err)

		// Retrieve and verify the balance
		balance, err := repo.Balance(context.Background(), "account1")
		assert.NoError(t, err)
		assert.Equal(t, int64(50), balance)
	})

	t.Run("AccountWithoutBalance", func(t *testing.T) {
		// Retrieve and verify the balance for an account without entries
		balance, err := repo.Balance(context.Background(), "account2")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), balance)
	})
}

func TestTransfer(t *testing.T) {
	repo, conn := setupTestDatabase(t)
	defer tearDownTestDatabase(conn)

	entries := entities.DoubleEntry{
		Credit: entities.Entry{Account: "account1", Direction: "credit", Amount: 100},
		Debit:  entities.Entry{Account: "account2", Direction: "debit", Amount: 100},
	}

	err := repo.Transfer(context.Background(), entries)
	assert.NoError(t, err)

	// Verify the state of the database after the transfer
	balance, err := repo.Balance(context.Background(), "account1")
	assert.NoError(t, err)
	assert.Equal(t, int64(-100), balance)

	balance, err = repo.Balance(context.Background(), "account2")
	assert.NoError(t, err)
	assert.Equal(t, int64(100), balance)
}

func TestTransfersHistory(t *testing.T) {
	repo, conn := setupTestDatabase(t)
	defer tearDownTestDatabase(conn)

	timeMock := time.Date(2023, 12, 5, 1, 56, 57, 324392000, time.UTC)
	account := "test_account"
	entries := []entities.Entry{
		{ID: 1, Account: account, Direction: "credit", Amount: 50, CreatedAt: &timeMock},
		{ID: 2, Account: account, Direction: "debit", Amount: 30, CreatedAt: &timeMock},
	}

	// Insert test entries into the database
	for _, entry := range entries {
		_, err := conn.Exec(context.Background(),
			"INSERT INTO entries (id, account, direction, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
			entry.ID, entry.Account, entry.Direction, entry.Amount, entry.CreatedAt,
		)
		assert.NoError(t, err)
	}

	// Retrieve and verify the history
	history, err := repo.TransfersHistory(context.Background(), account)
	assert.NoError(t, err)
	assert.Equal(t, entries, history)
}
