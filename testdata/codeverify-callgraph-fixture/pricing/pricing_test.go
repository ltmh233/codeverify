package pricing

import (
	"example.com/codeverify-callgraph-fixture/catalog"
	"example.com/codeverify-callgraph-fixture/promo"
	"example.com/codeverify-callgraph-fixture/shipping"
	"testing"
)

func newPricingEngine(t *testing.T) *Engine {
	t.Helper()
	c, err := catalog.NewCatalog([]catalog.Product{
		{SKU: "keyboard", PriceCents: 12000, Active: true},
		{SKU: "pen", PriceCents: 1500, Active: true},
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
	return e
}

func TestQuoteNoPromotion(t *testing.T) {
	e := newPricingEngine(t)
	q, err := e.Quote([]Line{{SKU: "pen", Quantity: 2}}, "", shipping.Domestic)
	if err != nil {
		t.Fatal(err)
	}
	if q.SubtotalCents != 3000 || q.DiscountCents != 0 || q.ShippingCents != 1200 || q.TotalCents != 4200 {
		t.Fatalf("quote = %+v", q)
	}
	if !q.ValidateTotals() {
		t.Fatal("inconsistent quote fields")
	}
}

func TestQuoteDiscountStillAboveThreshold(t *testing.T) {
	e := newPricingEngine(t)
	q, err := e.Quote([]Line{{SKU: "keyboard", Quantity: 2}}, "SAVE25", shipping.Domestic)
	if err != nil {
		t.Fatal(err)
	}
	if q.SubtotalCents != 24000 || q.DiscountCents != 6000 || q.ShippingCents != 0 || q.TotalCents != 18000 {
		t.Fatalf("quote = %+v", q)
	}
}
