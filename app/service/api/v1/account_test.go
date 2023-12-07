package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	v1 "github.com/julioc98/ledger/app/service/api/v1"
	"github.com/julioc98/ledger/app/service/api/v1/mocks"
	"github.com/julioc98/ledger/domain/account"
	"github.com/julioc98/ledger/domain/entities"
	"github.com/stretchr/testify/assert"
)

func NewTestHTTPRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	return req
}

func NewTestHTTPResponseWriter() http.ResponseWriter {
	return httptest.NewRecorder()
}

func TestCreateDepositHandler(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedCode  int
		mockError     error
		depositCalled bool
		invalidInput  bool
	}{
		{
			name:          "Valid request",
			requestBody:   `{"amount": 100}`,
			expectedCode:  http.StatusCreated,
			mockError:     nil,
			depositCalled: true,
		},
		{
			name:          "Internal server error",
			requestBody:   `{"amount": 100}`,
			expectedCode:  http.StatusInternalServerError,
			mockError:     fmt.Errorf("internal error"),
			depositCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepositUC := &mocks.DepositUseCaseMock{
				DepositFunc: func(ctx context.Context, input account.DepositInput) error {
					if tt.invalidInput {
						return fmt.Errorf("invalid input")
					}
					return tt.mockError
				},
			}

			req := NewTestHTTPRequest(http.MethodPost, "/deposit/example-account", strings.NewReader(tt.requestBody))
			w := NewTestHTTPResponseWriter()
			handler := v1.CreateDepositHandler(mockDepositUC)
			handler(w, req)

			assert.Equal(t, tt.expectedCode, w.(*httptest.ResponseRecorder).Code, "unexpected status code")

			if tt.depositCalled {
				calls := mockDepositUC.DepositCalls()
				assert.Len(t, calls, 1, "expected Deposit to be called once")

				if !tt.invalidInput {
					var depositInput account.DepositInput
					err := json.NewDecoder(strings.NewReader(tt.requestBody)).Decode(&depositInput)
					assert.NoError(t, err, "error decoding deposit input")

					assert.Equal(t, depositInput.Amount, calls[0].Input.Amount, "unexpected Deposit call arguments")
				} else {
					assert.Empty(t, calls[0].Input, "unexpected Deposit call arguments")
				}
			} else {
				assert.Empty(t, mockDepositUC.DepositCalls(), "expected Deposit not to be called")
			}
		})
	}
}

func TestCreateTransferHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedCode   int
		mockError      error
		transferCalled bool
		invalidInput   bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"to": "destination-account", "amount": 100}`,
			expectedCode:   http.StatusCreated,
			mockError:      nil,
			transferCalled: true,
		},
		{
			name:           "Internal server error",
			requestBody:    `{"to": "destination-account", "amount": 100}`,
			expectedCode:   http.StatusInternalServerError,
			mockError:      assert.AnError,
			transferCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransferUC := &mocks.TransferUseCaseMock{
				TransferFunc: func(ctx context.Context, input account.TransferInput) error {
					if tt.invalidInput {
						return assert.AnError
					}
					return tt.mockError
				},
			}
			req := NewTestHTTPRequest(http.MethodPost, "/transfer/source-account", strings.NewReader(tt.requestBody))
			accountParam := chi.URLParam(req, "account")
			req = req.WithContext(context.WithValue(req.Context(), "account", accountParam))

			w := NewTestHTTPResponseWriter()
			handler := v1.CreateTransferHandler(mockTransferUC)
			handler(w, req)

			assert.Equal(t, tt.expectedCode, w.(*httptest.ResponseRecorder).Code, "unexpected status code")

			if tt.transferCalled {
				calls := mockTransferUC.TransferCalls()
				assert.Len(t, calls, 1, "expected Transfer to be called once")

				if !tt.invalidInput {
					var transferRequest v1.CreateTransferRequest
					err := json.NewDecoder(strings.NewReader(tt.requestBody)).Decode(&transferRequest)
					assert.NoError(t, err, "error decoding transfer request")

					expectedInput := account.TransferInput{
						FromAccount: accountParam,
						ToAccount:   transferRequest.To,
						Amount:      transferRequest.Amount,
					}

					assert.Equal(t, expectedInput, calls[0].Input, "unexpected Transfer call arguments")

				}
			}
		})
	}
}

func TestGetBalanceHandler(t *testing.T) {
	tests := []struct {
		name          string
		account       string
		expectedCode  int
		expectedBody  string
		useCaseReturn int64
		useCaseError  error
	}{
		{
			name:          "Success",
			account:       "123",
			expectedCode:  http.StatusOK,
			expectedBody:  `{"balance":100}`,
			useCaseReturn: 100,
			useCaseError:  nil,
		},
		{
			name:          "Error from UseCase",
			account:       "456",
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `mocked error`,
			useCaseReturn: 0,
			useCaseError:  fmt.Errorf("mocked error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &mocks.BalanceUseCaseMock{
				BalanceFunc: func(ctx context.Context, account string) (int64, error) {
					return tt.useCaseReturn, tt.useCaseError
				},
			}

			handler := v1.GetBalanceHandler(mockUseCase)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("account", tt.account)
			req := httptest.NewRequest("GET", "/balance/"+tt.account, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			expectedBodyWithoutNewline := strings.ReplaceAll(tt.expectedBody, "\n", "")

			assert.Equal(t, tt.expectedBody, expectedBodyWithoutNewline)

			calls := mockUseCase.BalanceCalls()
			assert.Len(t, calls, 1)
			assert.Equal(t, tt.account, calls[0].Account)
		})
	}
}

func TestGetTransfersHistoryHandler(t *testing.T) {
	timeMock := time.Date(2023, 12, 5, 1, 56, 57, 324392000, time.UTC)
	tests := []struct {
		name           string
		account        string
		expectedCode   int
		expectedBody   string
		useCaseReturn  *account.TransfersHistoryOutput
		useCaseError   error
		transfersCalls int
	}{
		{
			name:           "Success",
			account:        "123",
			expectedCode:   http.StatusOK,
			expectedBody:   `{"entries":[{"ID":42,"Account":"jc","Direction":"debit","Amount":10000,"CreatedAt":"2023-12-05T01:56:57.324392Z"}]}`,
			useCaseReturn:  &account.TransfersHistoryOutput{Entries: []entities.Entry{{ID: 42, Account: "jc", Direction: "debit", Amount: 10000, CreatedAt: &timeMock}}},
			useCaseError:   nil,
			transfersCalls: 1,
		},
		{
			name:           "Error from UseCase",
			account:        "456",
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   `mocked error`,
			useCaseReturn:  nil,
			useCaseError:   fmt.Errorf("mocked error"),
			transfersCalls: 1,
		},
		{
			name:           "No Transfers",
			account:        "789",
			expectedCode:   http.StatusOK,
			expectedBody:   `{"entries":[]}`,
			useCaseReturn:  &account.TransfersHistoryOutput{Entries: []entities.Entry{}},
			useCaseError:   nil,
			transfersCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &mocks.TransfersHistoryUseCaseMock{
				TransfersHistoryFunc: func(ctx context.Context, account string) (*account.TransfersHistoryOutput, error) {
					return tt.useCaseReturn, tt.useCaseError
				},
			}

			handler := v1.GetTransfersHistoryHandler(mockUseCase)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("account", tt.account)
			req := httptest.NewRequest("GET", "/transfers/"+tt.account, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			// Remover '\n' da string expectedBody
			expectedBodyWithoutNewline := strings.ReplaceAll(tt.expectedBody, "\n", "")

			responseBodyWithoutNewline := strings.ReplaceAll(strings.TrimSpace(w.Body.String()), "\n", "")
			assert.Equal(t, expectedBodyWithoutNewline, responseBodyWithoutNewline)

			calls := mockUseCase.TransfersHistoryCalls()
			assert.Len(t, calls, tt.transfersCalls)
			if tt.transfersCalls > 0 {
				assert.Equal(t, tt.account, calls[0].AccountMoqParam)
			}
		})
	}
}
