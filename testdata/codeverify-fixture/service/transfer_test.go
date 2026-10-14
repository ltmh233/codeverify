package service

import (
	"errors"
	"testing"

	"example.com/codeverify-fixture/ledger"
)

func mustAccount(t *testing.T, balance, limit int64) *ledger.Account {
	t.Helper()
	a, err := ledger.NewAccount(balance, limit)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func TestTransferSuccess(t *testing.T) {
	src := mustAccount(t, 80, 100)
	dst := mustAccount(t, 10, 100)

	if err := Transfer(src, dst, 25); err != nil {
		t.Fatal(err)
	}
	if src.Balance() != 55 || dst.Balance() != 35 {
		t.Fatalf("balances = (%d, %d), want (55, 35)", src.Balance(), dst.Balance())
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	src := mustAccount(t, 10, 100)
	dst := mustAccount(t, 20, 100)

	if err := Transfer(src, dst, 15); !errors.Is(err, ledger.ErrInsufficientFunds) {
		t.Fatalf("Transfer error = %v, want ErrInsufficientFunds", err)
	}
	if src.Balance() != 10 || dst.Balance() != 20 {
		t.Fatalf("failed transfer changed balances to (%d, %d)", src.Balance(), dst.Balance())
	}
}

func TestTransferInvalidInputs(t *testing.T) {
	a := mustAccount(t, 50, 100)
	b := mustAccount(t, 30, 100)

	if err := Transfer(nil, b, 10); !errors.Is(err, ErrNilAccount) {
		t.Fatalf("Transfer(nil, b, 10) error = %v", err)
	}
	if err := Transfer(a, a, 10); !errors.Is(err, ErrSameAccount) {
		t.Fatalf("Transfer(a, a, 10) error = %v", err)
	}
	if err := Transfer(a, b, 0); !errors.Is(err, ledger.ErrInvalidAmount) {
		t.Fatalf("Transfer(a, b, 0) error = %v", err)
	}
	if a.Balance() != 50 || b.Balance() != 30 {
		t.Fatalf("invalid inputs changed balances to (%d, %d)", a.Balance(), b.Balance())
	}
}
