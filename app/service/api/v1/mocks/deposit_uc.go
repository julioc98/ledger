// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/julioc98/ledger/app/service/api/v1"
	"github.com/julioc98/ledger/domain/account"
	"sync"
)

// Ensure, that DepositUseCaseMock does implement v1.DepositUseCase.
// If this is not the case, regenerate this file with moq.
var _ v1.DepositUseCase = &DepositUseCaseMock{}

// DepositUseCaseMock is a mock implementation of v1.DepositUseCase.
//
//	func TestSomethingThatUsesDepositUseCase(t *testing.T) {
//
//		// make and configure a mocked v1.DepositUseCase
//		mockedDepositUseCase := &DepositUseCaseMock{
//			DepositFunc: func(ctx context.Context, input account.DepositInput) error {
//				panic("mock out the Deposit method")
//			},
//		}
//
//		// use mockedDepositUseCase in code that requires v1.DepositUseCase
//		// and then make assertions.
//
//	}
type DepositUseCaseMock struct {
	// DepositFunc mocks the Deposit method.
	DepositFunc func(ctx context.Context, input account.DepositInput) error

	// calls tracks calls to the methods.
	calls struct {
		// Deposit holds details about calls to the Deposit method.
		Deposit []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input account.DepositInput
		}
	}
	lockDeposit sync.RWMutex
}

// Deposit calls DepositFunc.
func (mock *DepositUseCaseMock) Deposit(ctx context.Context, input account.DepositInput) error {
	callInfo := struct {
		Ctx   context.Context
		Input account.DepositInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockDeposit.Lock()
	mock.calls.Deposit = append(mock.calls.Deposit, callInfo)
	mock.lockDeposit.Unlock()
	if mock.DepositFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.DepositFunc(ctx, input)
}

// DepositCalls gets all the calls that were made to Deposit.
// Check the length with:
//
//	len(mockedDepositUseCase.DepositCalls())
func (mock *DepositUseCaseMock) DepositCalls() []struct {
	Ctx   context.Context
	Input account.DepositInput
} {
	var calls []struct {
		Ctx   context.Context
		Input account.DepositInput
	}
	mock.lockDeposit.RLock()
	calls = mock.calls.Deposit
	mock.lockDeposit.RUnlock()
	return calls
}
