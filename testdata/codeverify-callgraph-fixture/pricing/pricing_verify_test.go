package pricing

import (
	"testing"

	"example.com/codeverify-callgraph-fixture/catalog"
	"example.com/codeverify-callgraph-fixture/promo"
	"example.com/codeverify-callgraph-fixture/shipping"
)

// 契约：免运费资格必须基于折扣后的 payable 金额。
// 场景：subtotal=12000 >= 门槛 10000，但 25% 折扣后 payable=9000 < 门槛，
// 因此必须收取国内运费 1200。
func TestVerifyFreeShippingUsesPayable(t *testing.T) {
	c, err := catalog.NewCatalog([]catalog.Product{
		{SKU: "keyboard", PriceCents: 12000, Active: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	p, err := promo.NewEngine([]promo.Rule{{Code: "SAVE25", Percent: 25}})
	if err != nil {
		t.Fatal(err)
	}
	s, err := shipping.NewCalculator(1200, 2500, 10000)
	if err != nil {
		t.Fatal(err)
	}
	e, err := NewEngine(c, p, s)
	if err != nil {
		t.Fatal(err)
	}

	q, err := e.Quote([]Line{{SKU: "keyboard", Quantity: 1}}, "SAVE25", shipping.Domestic)
	if err != nil {
		t.Fatal(err)
	}

	if q.SubtotalCents != 12000 || q.DiscountCents != 3000 {
		t.Fatalf("unexpected subtotal/discount: %+v", q)
	}
	if q.PayableMerchandiseCents() != 9000 {
		t.Fatalf("payable = %d, want 9000", q.PayableMerchandiseCents())
	}
	if q.ShippingCents != 1200 {
		t.Fatalf("ShippingCents = %d, want 1200 (payable 9000 < threshold 10000)", q.ShippingCents)
	}
	if q.TotalCents != 10200 {
		t.Fatalf("TotalCents = %d, want 10200", q.TotalCents)
	}
	if !q.ValidateTotals() {
		t.Fatalf("inconsistent quote: %+v", q)
	}
}
