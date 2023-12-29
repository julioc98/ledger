package account

import (
	"context"

	"github.com/julioc98/ledger/domain"
	"github.com/julioc98/ledger/domain/entities"
)

type TransferRepository interface {
	Transfer(ctx context.Context, entries entities.DoubleEntry) error
}

type AccountRepository interface {
	Balance(ctx context.Context, account string) (int64, error)
}

// TransferUseCase is a use case for Transfering money into an account.
type TransferUseCase struct {
	repo  TransferRepository
	aRepo AccountRepository
}

// NewTransferUseCase creates a new TransferUseCase.
func NewTransferUseCase(repo TransferRepository, aRepo AccountRepository) *TransferUseCase {
	return &TransferUseCase{repo, aRepo}
}

type TransferInput struct {
	FromAccount string
	ToAccount   string
	Amount      int64
}

// Transfer Transfers money into an account.
func (uc *TransferUseCase) Transfer(ctx context.Context, input TransferInput) error {

	credit, err := entities.NewEntry(0, input.FromAccount, entities.Credit, input.Amount, nil)
	if err != nil {
		return err
	}

	debit, err := entities.NewEntry(0, input.ToAccount, entities.Debit, input.Amount, nil)
	if err != nil {
		return err
	}

	de, err := entities.NewDoubleEntry(*debit, *credit)
	if err != nil {
		return err
	}

	balance, err := uc.aRepo.Balance(ctx, input.FromAccount)
	if err != nil {
		return err
	}

	if balance < input.Amount {
		return domain.ErrInsufficientFunds
	}

	return uc.repo.Transfer(ctx, *de)
}
