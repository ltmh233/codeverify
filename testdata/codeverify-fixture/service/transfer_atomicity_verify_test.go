package service

import (
	"errors"
	"testing"

	"example.com/codeverify-fixture/ledger"
)

// TestVerifyTransferAtomicityOnCreditFailure isolates the atomicity contract of
// Transfer: "On any returned error, neither account balance may change."
//
// The source account has sufficient funds, so the debit succeeds. The
// destination account is near its limit, so the credit fails with
// ErrBalanceLimit. A correct implementation must not leave the source debited.
func TestVerifyTransferAtomicityOnCreditFailure(t *testing.T) {
	src := mustAccount(t, 100, 100)
	dst := mustAccount(t, 95, 100)

	before := [2]int64{src.Balance(), dst.Balance()}

	err := Transfer(src, dst, 10)
	if !errors.Is(err, ledger.ErrBalanceLimit) {
		t.Fatalf("Transfer error = %v, want ErrBalanceLimit", err)
	}

	after := [2]int64{src.Balance(), dst.Balance()}
	if after != before {
		t.Fatalf("balances changed on failed transfer: before=%v after=%v", before, after)
	}
}
