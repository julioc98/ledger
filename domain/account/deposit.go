package account

import (
	"context"

	"github.com/julioc98/ledger/domain/entities"
)

const (
	ExternalAccount = "external"
)

type DepositRepository interface {
	Transfer(ctx context.Context, entries entities.DoubleEntry) error
}

// DepositUseCase is a use case for depositing money into an account.
type DepositUseCase struct {
	repo DepositRepository
}

// NewDepositUseCase creates a new DepositUseCase.
func NewDepositUseCase(repo DepositRepository) *DepositUseCase {
	return &DepositUseCase{repo}
}

type DepositInput struct {
	Account string
	Amount  int64
}

// Deposit deposits money into an account.
func (uc *DepositUseCase) Deposit(ctx context.Context, input DepositInput) error {
	debit, err := entities.NewEntry(0, input.Account, entities.Debit, input.Amount, nil)
	if err != nil {
		return err
	}

	credit, err := entities.NewEntry(0, ExternalAccount, entities.Credit, input.Amount, nil)
	if err != nil {
		return err
	}

	de, err := entities.NewDoubleEntry(*debit, *credit)
	if err != nil {
		return err
	}

	return uc.repo.Transfer(ctx, *de)
}
