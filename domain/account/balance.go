package account

import "context"

type BalanceRepository interface {
	Balance(ctx context.Context, account string) (int64, error)
}

type BalanceUseCase struct {
	repo BalanceRepository
}

func NewBalanceUseCase(repo BalanceRepository) *BalanceUseCase {
	return &BalanceUseCase{repo}
}

func (uc *BalanceUseCase) Balance(ctx context.Context, account string) (int64, error) {
	return uc.repo.Balance(ctx, account)
}
