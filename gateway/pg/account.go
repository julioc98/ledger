package pg

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/julioc98/ledger/domain/entities"
)

type AccountPgxRepository struct {
	conn *pgx.Conn
}

func NewAccountPgxRepository(conn *pgx.Conn) *AccountPgxRepository {
	return &AccountPgxRepository{conn}
}

func (r *AccountPgxRepository) Transfer(ctx context.Context, entries entities.DoubleEntry) error {
	queryTpl := `
		INSERT INTO entries (account, direction, amount)
		VALUES ($1, $2, $3)
	`

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, queryTpl, entries.Credit.Account, entries.Credit.Direction, entries.Credit.Amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, queryTpl, entries.Debit.Account, entries.Debit.Direction, entries.Debit.Amount)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *AccountPgxRepository) Balance(ctx context.Context, account string) (int64, error) {
	queryTpl := `
		SELECT
			SUM(CASE WHEN direction = 'debit' THEN amount ELSE 0 END) -
			SUM(CASE WHEN direction = 'credit' THEN amount ELSE 0 END) AS balance
		FROM entries
		WHERE account = $1;
	`

	var balance sql.NullInt64
	err := r.conn.QueryRow(ctx, queryTpl, account).Scan(&balance)
	if err != nil {
		return 0, err
	}

	if !balance.Valid {
		return 0, nil
	}

	return balance.Int64, nil
}

func (r *AccountPgxRepository) TransfersHistory(ctx context.Context, account string) ([]entities.Entry, error) {
	queryTpl := `
		SELECT id, account, direction, amount, created_at
		FROM entries
		WHERE account = $1
		ORDER BY created_at DESC;
	`

	rows, err := r.conn.Query(ctx, queryTpl, account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entities.Entry
	for rows.Next() {
		var e entities.Entry
		err := rows.Scan(&e.ID, &e.Account, &e.Direction, &e.Amount, &e.CreatedAt)
		if err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}
