package account

import (
	"context"
	"testing"

	"github.com/julioc98/ledger/domain/entities"
	"github.com/stretchr/testify/assert"
)

type mockDepositRepository struct {
	transferFunc func(ctx context.Context, entries entities.DoubleEntry) error
}

func (m *mockDepositRepository) Transfer(ctx context.Context, entries entities.DoubleEntry) error {
	return m.transferFunc(ctx, entries)
}

func TestDepositUseCase_Deposit(t *testing.T) {
	mockRepo := &mockDepositRepository{
		transferFunc: func(ctx context.Context, entries entities.DoubleEntry) error {

			return nil
		},
	}

	uc := NewDepositUseCase(mockRepo)

	tests := []struct {
		name     string
		input    DepositInput
		expected error
	}{
		{
			name: "should deposit money into an account",
			input: DepositInput{
				Account: "123456789",
				Amount:  1000,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Deposit(context.Background(), tt.input)
			assert.ErrorIs(t, err, tt.expected, "unexpected error")
		})
	}
}
