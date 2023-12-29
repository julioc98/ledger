package account

import (
	"context"

	"github.com/julioc98/ledger/domain/entities"
)

type TransfersHistoryRepository interface {
	TransfersHistory(ctx context.Context, account string) ([]entities.Entry, error)
}

type TransfersHistoryUseCase struct {
	repo TransfersHistoryRepository
}

type TransfersHistoryOutput struct {
	Entries []entities.Entry `json:"entries"`
}

func NewTransfersHistoryUseCase(repo TransfersHistoryRepository) *TransfersHistoryUseCase {
	return &TransfersHistoryUseCase{repo}
}

func (uc *TransfersHistoryUseCase) TransfersHistory(ctx context.Context, account string) (*TransfersHistoryOutput, error) {
	entries, err := uc.repo.TransfersHistory(ctx, account)
	if err != nil {
		return nil, err
	}

	return &TransfersHistoryOutput{entries}, nil
}
