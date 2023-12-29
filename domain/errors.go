package domain

import "errors"

var (
	// ErrInsufficientFunds is an error for insufficient funds.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidDirection is an error for invalid direction.
	ErrIvalidDirection = errors.New("invalid direction")

	// ErrAmountsNotMatch is an error for amounts not matching.
	ErrAmountsNotMatch = errors.New("debit and credit amounts must match")
)
