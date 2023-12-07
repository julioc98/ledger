package account

import (
	"context"
	"testing"

	"github.com/julioc98/ledger/domain/entities"
	"github.com/stretchr/testify/assert"
)

type mockTransfersHistoryRepository struct{}

func (m *mockTransfersHistoryRepository) TransfersHistory(ctx context.Context, account string) ([]entities.Entry, error) {
	// Mock implementation of TransfersHistory repository method
	// Return some dummy data for testing
	entries := []entities.Entry{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
		{ID: 3, Amount: 300},
	}
	return entries, nil
}

func TestTransfersHistoryUseCase_TransfersHistory(t *testing.T) {
	repo := &mockTransfersHistoryRepository{}
	useCase := NewTransfersHistoryUseCase(repo)

	ctx := context.Background()
	account := "example-account"

	output, err := useCase.TransfersHistory(ctx, account)
	assert.NoError(t, err, "unexpected error")

	expectedEntries := []entities.Entry{
		{ID: 1, Amount: 100},
		{ID: 2, Amount: 200},
		{ID: 3, Amount: 300},
	}

	assert.Len(t, output.Entries, len(expectedEntries), "unexpected number of entries")

	for i, entry := range output.Entries {
		assert.Equal(t, expectedEntries[i], entry, "unexpected entry")
	}
}
