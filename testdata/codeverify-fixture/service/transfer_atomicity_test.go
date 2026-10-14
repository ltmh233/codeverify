package service

import (
	"errors"
	"testing"

	"example.com/codeverify-fixture/ledger"
)

// TestTransferAtomicityOnCreditFailure checks the documented contract of
// Transfer: "On any returned error, neither account balance may change."
//
// src has enough funds for the debit, but dst is already close to its limit
// so the credit must fail with ErrBalanceLimit. In that case the transfer must
// leave both balances unchanged.
func TestTransferAtomicityOnCreditFailure(t *testing.T) {
	src := mustAccount(t, 80, 100)
	dst := mustAccount(t, 90, 100)

	err := Transfer(src, dst, 25)
	if !errors.Is(err, ledger.ErrBalanceLimit) {
		t.Fatalf("Transfer error = %v, want ErrBalanceLimit", err)
	}

	if src.Balance() != 80 {
		t.Fatalf("src.Balance() = %d after failed transfer, want 80 (unchanged)", src.Balance())
	}
	if dst.Balance() != 90 {
		t.Fatalf("dst.Balance() = %d after failed transfer, want 90 (unchanged)", dst.Balance())
	}
}
