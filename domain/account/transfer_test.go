package account

import (
	"context"
	"errors"
	"testing"

	"github.com/julioc98/ledger/domain"
	"github.com/julioc98/ledger/domain/entities"
	"github.com/stretchr/testify/assert"
)

type mockTransferRepo struct {
	transferFunc func(ctx context.Context, entries entities.DoubleEntry) error
}

func (m *mockTransferRepo) Transfer(ctx context.Context, entries entities.DoubleEntry) error {
	return m.transferFunc(ctx, entries)
}

type mockAccountRepo struct {
	balanceFunc func(ctx context.Context, account string) (int64, error)
}

func (m *mockAccountRepo) Balance(ctx context.Context, account string) (int64, error) {
	return m.balanceFunc(ctx, account)
}

func TestTransferUseCase_Transfer(t *testing.T) {
	// Create mock repositories
	mockTransferRepo := &mockTransferRepo{}
	mockAccountRepo := &mockAccountRepo{}

	// Create the TransferUseCase with the mock repositories
	transferUseCase := NewTransferUseCase(mockTransferRepo, mockAccountRepo)

	// Define test cases
	testCases := []struct {
		name           string
		fromAccount    string
		toAccount      string
		amount         int64
		balance        int64
		expectedError  error
		transferCalled bool
	}{
		{
			name:           "Sufficient funds",
			fromAccount:    "account1",
			toAccount:      "account2",
			amount:         100,
			balance:        200,
			expectedError:  nil,
			transferCalled: true,
		},
		{
			name:           "Insufficient funds",
			fromAccount:    "account1",
			toAccount:      "account2",
			amount:         200,
			balance:        100,
			expectedError:  domain.ErrInsufficientFunds,
			transferCalled: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock AccountRepository's Balance function
			mockAccountRepo.balanceFunc = func(ctx context.Context, account string) (int64, error) {
				if account == tc.fromAccount {
					return tc.balance, nil
				}
				return 0, errors.New("unexpected account")
			}

			// Set up the mock TransferRepository's Transfer function
			mockTransferRepo.transferFunc = func(ctx context.Context, entries entities.DoubleEntry) error {
				tc.transferCalled = true
				return nil
			}

			// Call the Transfer method of the TransferUseCase
			err := transferUseCase.Transfer(context.Background(), TransferInput{
				FromAccount: tc.fromAccount,
				ToAccount:   tc.toAccount,
				Amount:      tc.amount,
			})

			// Check if the error matches the expected error
			assert.Equal(t, tc.expectedError, err, "unexpected error")

		})
	}
}
