// Package service implements operations that span multiple accounts.
package service

import (
	"errors"
	"fmt"

	"example.com/codeverify-fixture/ledger"
)

var (
	ErrNilAccount  = errors.New("source and destination accounts are required")
	ErrSameAccount = errors.New("source and destination must be different accounts")
)

// Transfer moves a positive amount from src to dst.
// It must either update both account balances or leave both unchanged.
// On any returned error, neither account balance may change.
// The caller must not use these accounts concurrently during the transfer.
func Transfer(src, dst *ledger.Account, amount int64) error {
	if src == nil || dst == nil {
		return ErrNilAccount
	}
	if src == dst {
		return ErrSameAccount
	}

	if err := src.Debit(amount); err != nil {
		return fmt.Errorf("debit: %w", err)
	}
	if err := dst.Credit(amount); err != nil {
		return fmt.Errorf("credit: %w", err)
	}
	return nil
}
