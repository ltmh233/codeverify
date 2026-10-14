package ledger

import (
	"errors"
	"testing"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount(40, 100)
	if err != nil {
		t.Fatal(err)
	}
	if got := account.Balance(); got != 40 {
		t.Fatalf("Balance() = %d, want 40", got)
	}

	if _, err := NewAccount(101, 100); !errors.Is(err, ErrInvalidAccount) {
		t.Fatalf("NewAccount(101, 100) error = %v, want ErrInvalidAccount", err)
	}
}

func TestDebit(t *testing.T) {
	account, _ := NewAccount(40, 100)
	if err := account.Debit(15); err != nil {
		t.Fatal(err)
	}
	if got := account.Balance(); got != 25 {
		t.Fatalf("Balance() = %d, want 25", got)
	}
	if err := account.Debit(30); !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("Debit(30) error = %v, want ErrInsufficientFunds", err)
	}
	if got := account.Balance(); got != 25 {
		t.Fatalf("failed Debit changed balance to %d", got)
	}
}

func TestCredit(t *testing.T) {
	account, _ := NewAccount(40, 100)
	if err := account.Credit(20); err != nil {
		t.Fatal(err)
	}
	if got := account.Balance(); got != 60 {
		t.Fatalf("Balance() = %d, want 60", got)
	}
	if err := account.Credit(50); !errors.Is(err, ErrBalanceLimit) {
		t.Fatalf("Credit(50) error = %v, want ErrBalanceLimit", err)
	}
	if got := account.Balance(); got != 60 {
		t.Fatalf("failed Credit changed balance to %d", got)
	}
}
