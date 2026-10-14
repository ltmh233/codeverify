package checkout

import (
	"errors"
	"testing"

	"example.com/codeverify-callgraph-fixture/catalog"
	"example.com/codeverify-callgraph-fixture/promo"
	"example.com/codeverify-callgraph-fixture/shipping"
)

// 验证 catalog/promo/shipping 的错误分支与边界语义，排查 R1 之外的潜在缺陷。
func TestVerifyOtherPackageBoundaries(t *testing.T) {
	// catalog：未激活商品应返回 ErrUnavailable
	c, err := catalog.NewCatalog([]catalog.Product{
		{SKU: "pen", PriceCents: 1500, Active: true},
		{SKU: "old", PriceCents: 1500, Active: false},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Lookup("old"); !errors.Is(err, catalog.ErrUnavailable) {
		t.Fatalf("inactive lookup err = %v, want ErrUnavailable", err)
	}

	// promo：cap 生效、未知码、负 subtotal
	p, err := promo.NewEngine([]promo.Rule{{Code: "S", Percent: 50, MaxDiscountCents: 100}})
	if err != nil {
		t.Fatal(err)
	}
	if d, err := p.Discount("S", 300); err != nil || d != 100 {
		t.Fatalf("Discount cap = %d, %v, want 100", d, err)
	}
	if _, err := p.Discount("X", 100); !errors.Is(err, promo.ErrUnknownCode) {
		t.Fatalf("unknown code err = %v, want ErrUnknownCode", err)
	}
	if _, err := p.Discount("", -1); !errors.Is(err, promo.ErrBadSubtotal) {
		t.Fatalf("negative subtotal err = %v, want ErrBadSubtotal", err)
	}

	// shipping：负金额、未知区域、非法配置
	s, err := shipping.NewCalculator(1200, 2500, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Fee(-1, shipping.Domestic); !errors.Is(err, shipping.ErrInvalidAmount) {
		t.Fatalf("negative fee err = %v, want ErrInvalidAmount", err)
	}
	if _, err := s.Fee(100, shipping.Zone("mars")); !errors.Is(err, shipping.ErrUnknownZone) {
		t.Fatalf("unknown zone err = %v, want ErrUnknownZone", err)
	}
	if _, err := shipping.NewCalculator(1200, 2500, 0); !errors.Is(err, shipping.ErrInvalidRates) {
		t.Fatalf("invalid rates err = %v, want ErrInvalidRates", err)
	}
}
