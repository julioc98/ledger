package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/julioc98/ledger/domain"
	"github.com/julioc98/ledger/domain/account"
)

type CreateDepositRequest struct {
	Amount int64 `json:"amount"`
}

//go:generate moq -stub -pkg mocks -out mocks/deposit_uc.go . DepositUseCase
type DepositUseCase interface {
	Deposit(ctx context.Context, input account.DepositInput) error
}

// CreateDepositHandler godoc
// @Summary      Create a deposit
// @Description  Create a deposit for the specified account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account path string true "Account ID"
// @Param request body CreateDepositRequest true "Deposit details"
// @Success      201 "Deposit created successfully"
// @Failure      400 "Invalid request payload"
// @Failure      500 "Internal Server Error"
// @Router       /api/v1/account/{account}/deposit [post]
func CreateDepositHandler(uc DepositUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountP := chi.URLParam(r, "account")
		var req CreateDepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		input := account.DepositInput{
			Account: accountP,
			Amount:  req.Amount,
		}
		if err := uc.Deposit(r.Context(), input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

type CreateTransferRequest struct {
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

//go:generate moq -stub -pkg mocks -out mocks/transfer_uc.go . TransferUseCase
type TransferUseCase interface {
	Transfer(ctx context.Context, input account.TransferInput) error
}

// CreateTransferHandler godoc
// @Summary      Create a transfer
// @Description  Create a transfer from the specified account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account path string true "Account ID"
// @Param request body CreateTransferRequest true "Transfer details"
// @Success      201 "Transfer created successfully"
// @Failure      400 "Invalid request payload"
// @Failure      402 "Insufficient funds"
// @Failure      500 "Internal Server Error"
// @Router       /api/v1/account/{account}/transfers [post]
func CreateTransferHandler(uc TransferUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountP := chi.URLParam(r, "account")
		var req CreateTransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		input := account.TransferInput{
			FromAccount: accountP,
			ToAccount:   req.To,
			Amount:      req.Amount,
		}
		if err := uc.Transfer(r.Context(), input); err != nil {
			if errors.Is(err, domain.ErrInsufficientFunds) {
				http.Error(w, err.Error(), http.StatusPaymentRequired)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

//go:generate moq -stub -pkg mocks -out mocks/balance_uc.go . BalanceUseCase
type BalanceUseCase interface {
	Balance(ctx context.Context, account string) (int64, error)
}

// GetBalanceHandler godoc
// @Summary      Get account balance
// @Description  Get the balance of the specified account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account path string true "Account ID"
// @Success      200 "Balance retrieved successfully"
// @Failure      500 "Internal Server Error"
// @Router       /api/v1/account/{account}/balance [get]
func GetBalanceHandler(uc BalanceUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := chi.URLParam(r, "account")
		balance, err := uc.Balance(r.Context(), account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int64{"balance": balance})
	}
}

//go:generate moq -stub -pkg mocks -out mocks/transfers_history_uc.go . TransfersHistoryUseCase
type TransfersHistoryUseCase interface {
	TransfersHistory(ctx context.Context, account string) (*account.TransfersHistoryOutput, error)
}

// GetTransfersHistoryHandler godoc
// @Summary      Get transfers history
// @Description  Get the transfer history of the specified account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account path string true "Account ID"
// @Success      200 "Transfers history retrieved successfully"
// @Failure      500 "Internal Server Error"
// @Router       /api/v1/account/{account}/transfers [get]
func GetTransfersHistoryHandler(uc TransfersHistoryUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account := chi.URLParam(r, "account")
		transfers, err := uc.TransfersHistory(r.Context(), account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transfers)
	}
}
