package account

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockBalanceRepository struct{}

func (m *mockBalanceRepository) Balance(ctx context.Context, account string) (int64, error) {
	if account == "validAccount" {
		return 1000, nil
	}
	return 0, errors.New("account not found")
}

func TestBalanceUseCase_Balance(t *testing.T) {
	repo := &mockBalanceRepository{}
	useCase := NewBalanceUseCase(repo)

	// Test case 1: Valid account
	ctx := context.Background()
	account := "validAccount"
	expectedBalance := int64(1000)

	balance, err := useCase.Balance(ctx, account)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, expectedBalance, balance, "unexpected balance value")

	// Test case 2: Invalid account
	account = "invalidAccount"

	balance, err = useCase.Balance(ctx, account)
	assert.Error(t, err, "expected error")
	assert.Zero(t, balance, "expected balance 0")
}
