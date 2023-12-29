package entities

import (
	"time"

	"github.com/julioc98/ledger/domain"
)

// Direction is the direction of an entry.
type Direction string

const (
	// Debit is a debit entry.
	Debit Direction = "debit"
	// Credit is a credit entry.
	Credit Direction = "credit"
)

// Entry is a representation of a single entry in the ledger.
type Entry struct {
	ID        int64
	Account   string
	Direction Direction
	Amount    int64
	CreatedAt *time.Time
}

// NewEntry creates a new Entry with direction validation.
func NewEntry(id int64, account string, direction Direction, amount int64, createdAt *time.Time) (*Entry, error) {
	if direction != Debit && direction != Credit {
		return nil, domain.ErrIvalidDirection
	}

	return &Entry{
		ID:        id,
		Account:   account,
		Direction: direction,
		Amount:    amount,
		CreatedAt: createdAt,
	}, nil
}

// DoubleEntry is a representation of a double entry in the ledger.
type DoubleEntry struct {
	Debit  Entry
	Credit Entry
}

// NewDoubleEntry creates a new DoubleEntry with validation.
func NewDoubleEntry(debit, credit Entry) (*DoubleEntry, error) {
	// Validate that the amounts match.
	if debit.Amount != credit.Amount {
		return nil, domain.ErrAmountsNotMatch
	}

	return &DoubleEntry{
		Debit:  debit,
		Credit: credit,
	}, nil
}
