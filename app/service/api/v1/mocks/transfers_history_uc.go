// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/julioc98/ledger/app/service/api/v1"
	"github.com/julioc98/ledger/domain/account"
	"sync"
)

// Ensure, that TransfersHistoryUseCaseMock does implement v1.TransfersHistoryUseCase.
// If this is not the case, regenerate this file with moq.
var _ v1.TransfersHistoryUseCase = &TransfersHistoryUseCaseMock{}

// TransfersHistoryUseCaseMock is a mock implementation of v1.TransfersHistoryUseCase.
//
//	func TestSomethingThatUsesTransfersHistoryUseCase(t *testing.T) {
//
//		// make and configure a mocked v1.TransfersHistoryUseCase
//		mockedTransfersHistoryUseCase := &TransfersHistoryUseCaseMock{
//			TransfersHistoryFunc: func(ctx context.Context, accountMoqParam string) (*account.TransfersHistoryOutput, error) {
//				panic("mock out the TransfersHistory method")
//			},
//		}
//
//		// use mockedTransfersHistoryUseCase in code that requires v1.TransfersHistoryUseCase
//		// and then make assertions.
//
//	}
type TransfersHistoryUseCaseMock struct {
	// TransfersHistoryFunc mocks the TransfersHistory method.
	TransfersHistoryFunc func(ctx context.Context, accountMoqParam string) (*account.TransfersHistoryOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// TransfersHistory holds details about calls to the TransfersHistory method.
		TransfersHistory []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// AccountMoqParam is the accountMoqParam argument value.
			AccountMoqParam string
		}
	}
	lockTransfersHistory sync.RWMutex
}

// TransfersHistory calls TransfersHistoryFunc.
func (mock *TransfersHistoryUseCaseMock) TransfersHistory(ctx context.Context, accountMoqParam string) (*account.TransfersHistoryOutput, error) {
	callInfo := struct {
		Ctx             context.Context
		AccountMoqParam string
	}{
		Ctx:             ctx,
		AccountMoqParam: accountMoqParam,
	}
	mock.lockTransfersHistory.Lock()
	mock.calls.TransfersHistory = append(mock.calls.TransfersHistory, callInfo)
	mock.lockTransfersHistory.Unlock()
	if mock.TransfersHistoryFunc == nil {
		var (
			transfersHistoryOutputOut *account.TransfersHistoryOutput
			errOut                    error
		)
		return transfersHistoryOutputOut, errOut
	}
	return mock.TransfersHistoryFunc(ctx, accountMoqParam)
}

// TransfersHistoryCalls gets all the calls that were made to TransfersHistory.
// Check the length with:
//
//	len(mockedTransfersHistoryUseCase.TransfersHistoryCalls())
func (mock *TransfersHistoryUseCaseMock) TransfersHistoryCalls() []struct {
	Ctx             context.Context
	AccountMoqParam string
} {
	var calls []struct {
		Ctx             context.Context
		AccountMoqParam string
	}
	mock.lockTransfersHistory.RLock()
	calls = mock.calls.TransfersHistory
	mock.lockTransfersHistory.RUnlock()
	return calls
}
