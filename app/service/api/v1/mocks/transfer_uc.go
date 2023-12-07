// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/julioc98/ledger/app/service/api/v1"
	"github.com/julioc98/ledger/domain/account"
	"sync"
)

// Ensure, that TransferUseCaseMock does implement v1.TransferUseCase.
// If this is not the case, regenerate this file with moq.
var _ v1.TransferUseCase = &TransferUseCaseMock{}

// TransferUseCaseMock is a mock implementation of v1.TransferUseCase.
//
//	func TestSomethingThatUsesTransferUseCase(t *testing.T) {
//
//		// make and configure a mocked v1.TransferUseCase
//		mockedTransferUseCase := &TransferUseCaseMock{
//			TransferFunc: func(ctx context.Context, input account.TransferInput) error {
//				panic("mock out the Transfer method")
//			},
//		}
//
//		// use mockedTransferUseCase in code that requires v1.TransferUseCase
//		// and then make assertions.
//
//	}
type TransferUseCaseMock struct {
	// TransferFunc mocks the Transfer method.
	TransferFunc func(ctx context.Context, input account.TransferInput) error

	// calls tracks calls to the methods.
	calls struct {
		// Transfer holds details about calls to the Transfer method.
		Transfer []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input account.TransferInput
		}
	}
	lockTransfer sync.RWMutex
}

// Transfer calls TransferFunc.
func (mock *TransferUseCaseMock) Transfer(ctx context.Context, input account.TransferInput) error {
	callInfo := struct {
		Ctx   context.Context
		Input account.TransferInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockTransfer.Lock()
	mock.calls.Transfer = append(mock.calls.Transfer, callInfo)
	mock.lockTransfer.Unlock()
	if mock.TransferFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.TransferFunc(ctx, input)
}

// TransferCalls gets all the calls that were made to Transfer.
// Check the length with:
//
//	len(mockedTransferUseCase.TransferCalls())
func (mock *TransferUseCaseMock) TransferCalls() []struct {
	Ctx   context.Context
	Input account.TransferInput
} {
	var calls []struct {
		Ctx   context.Context
		Input account.TransferInput
	}
	mock.lockTransfer.RLock()
	calls = mock.calls.Transfer
	mock.lockTransfer.RUnlock()
	return calls
}