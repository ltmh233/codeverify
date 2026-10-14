// Package ledger provides a small in-memory account model.
package ledger

import "errors"

var (
	ErrInvalidAccount    = errors.New("invalid account configuration")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrBalanceLimit      = errors.New("account balance limit exceeded")
)

// Account is an in-memory account with an upper balance limit.
// Its methods are intended to be used with a non-nil receiver.
// Account is not safe for concurrent use.
type Account struct {
	balance int64
	limit   int64
}

// NewAccount creates an account with an initial balance in [0, limit].
func NewAccount(initialBalance, limit int64) (*Account, error) {
	if limit <= 0 || initialBalance < 0 || initialBalance > limit {
		return nil, ErrInvalidAccount
	}
	return &Account{balance: initialBalance, limit: limit}, nil
}

// Balance returns the current account balance.
func (a *Account) Balance() int64 {
	return a.balance
}

func validateAmount(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

// Debit removes a positive amount if the account has sufficient funds.
// On error, the account balance is unchanged.
func (a *Account) Debit(amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if a.balance < amount {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	return nil
}

// Credit adds a positive amount if the result does not exceed the balance limit.
// On error, the account balance is unchanged.
func (a *Account) Credit(amount int64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}
	if amount > a.limit-a.balance {
		return ErrBalanceLimit
	}
	a.balance += amount
	return nil
}
