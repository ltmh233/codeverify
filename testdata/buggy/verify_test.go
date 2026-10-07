package buggy

import "testing"

// Documents the stated contract: "amount must be positive".
// A negative amount should be rejected, not silently reverse the transfer.

func TestNegativeAmountCreatesMoney(t *testing.T) {
	src := &Account{Balance: 100}
	dst := &Account{Balance: 100}

	// amount <= 0 is invalid per the documented contract.
	err := Transfer(src, dst, -50)

	t.Logf("err=%v src=%d dst=%d total=%d", err, src.Balance, dst.Balance, src.Balance+dst.Balance)

	if err == nil {
		t.Errorf("expected error for negative amount, got nil (src=%d dst=%d)", src.Balance, dst.Balance)
	}
}

func TestZeroAmount(t *testing.T) {
	src := &Account{Balance: 100}
	dst := &Account{Balance: 100}
	err := Transfer(src, dst, 0)
	t.Logf("zero amount: err=%v src=%d dst=%d", err, src.Balance, dst.Balance)
}

func TestNegativeAmountStealsFromDst(t *testing.T) {
	src := &Account{Balance: 50}
	dst := &Account{Balance: 1000}

	err := Transfer(src, dst, -200)
	t.Logf("neg amount: err=%v src=%d dst=%d", err, src.Balance, dst.Balance)

	if err == nil {
		t.Errorf("negative amount accepted: src=%d dst=%d (no error)", src.Balance, dst.Balance)
	}
}

func TestNormalTransferStillWorks(t *testing.T) {
	src := &Account{Balance: 100}
	dst := &Account{Balance: 0}
	if err := Transfer(src, dst, 40); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Balance != 60 || dst.Balance != 40 {
		t.Errorf("unexpected balances src=%d dst=%d", src.Balance, dst.Balance)
	}
}

func TestInsufficientBalance(t *testing.T) {
	src := &Account{Balance: 10}
	dst := &Account{Balance: 0}
	if err := Transfer(src, dst, 100); err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}
