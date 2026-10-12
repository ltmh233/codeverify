package buggy

import (
	"math"
	"testing"
)

// R2: 合法输入（amount>0 且 src.Balance>=amount）下，dst.Balance+amount 溢出 int64。
func TestTransferDstOverflow(t *testing.T) {
	src := &Account{Balance: math.MaxInt64}
	dst := &Account{Balance: 1}

	totalBefore := src.Balance + dst.Balance

	err := Transfer(src, dst, math.MaxInt64)
	if err != nil {
		t.Fatalf("unexpected error with sufficient balance: %v", err)
	}

	t.Logf("after transfer: src=%d dst=%d", src.Balance, dst.Balance)

	if dst.Balance < 0 {
		t.Fatalf("dst balance overflowed to negative: %d", dst.Balance)
	}

	// 总额守恒检查（用更大的比较避免再次溢出，这里用等价关系）
	if src.Balance < 0 || dst.Balance < 0 {
		t.Fatalf("negative balance detected src=%d dst=%d", src.Balance, dst.Balance)
	}
	if src.Balance+dst.Balance != totalBefore {
		t.Fatalf("total not conserved: before=%d after=%d", totalBefore, src.Balance+dst.Balance)
	}
}

// R1: amount 为负数时的行为。
func TestTransferNegativeAmount(t *testing.T) {
	src := &Account{Balance: 100}
	dst := &Account{Balance: 50}

	err := Transfer(src, dst, -30)
	t.Logf("err=%v src=%d dst=%d", err, src.Balance, dst.Balance)

	if err == nil {
		t.Fatalf("expected error for negative amount, got nil (src=%d dst=%d)", src.Balance, dst.Balance)
	}
}
