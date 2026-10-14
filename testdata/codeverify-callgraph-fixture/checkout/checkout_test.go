package checkout

import (
	"errors"
	"testing"

	"example.com/codeverify-callgraph-fixture/catalog"
	"example.com/codeverify-callgraph-fixture/pricing"
	"example.com/codeverify-callgraph-fixture/promo"
	"example.com/codeverify-callgraph-fixture/shipping"
)

func newCheckoutService(t *testing.T) *Service {
	t.Helper()
	c, err := catalog.NewCatalog([]catalog.Product{
		{SKU: "keyboard", PriceCents: 12000, Active: true},
		{SKU: "notebook", PriceCents: 6000, Active: true},
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
	e, err := pricing.NewEngine(c, p, s)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(e)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestCheckoutQuoteNoPromotion(t *testing.T) {
	s := newCheckoutService(t)
	q, err := s.Quote(Request{Lines: []pricing.Line{{SKU: "notebook", Quantity: 1}}, Zone: shipping.Domestic})
	if err != nil {
		t.Fatal(err)
	}
	if q.TotalCents != 7200 {
		t.Fatalf("TotalCents = %d, want 7200", q.TotalCents)
	}
}

func TestCheckoutQuoteInternational(t *testing.T) {
	s := newCheckoutService(t)
	q, err := s.Quote(Request{Lines: []pricing.Line{{SKU: "keyboard", Quantity: 1}}, Zone: shipping.International})
	if err != nil {
		t.Fatal(err)
	}
	if q.TotalCents != 14500 {
		t.Fatalf("TotalCents = %d, want 14500", q.TotalCents)
	}
}

func TestCheckoutRejectsEmptyCart(t *testing.T) {
	s := newCheckoutService(t)
	_, err := s.Quote(Request{Zone: shipping.Domestic})
	if !errors.Is(err, ErrEmptyCart) {
		t.Fatalf("empty cart error = %v", err)
	}
}
