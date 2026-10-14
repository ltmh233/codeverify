// Package pricing combines merchandise, promotion and shipping calculations.
package pricing

import (
	"errors"

	"example.com/codeverify-callgraph-fixture/catalog"
	"example.com/codeverify-callgraph-fixture/promo"
	"example.com/codeverify-callgraph-fixture/shipping"
)

var (
	ErrInvalidLines  = errors.New("invalid cart lines")
	ErrNilDependency = errors.New("missing pricing dependency")
)

type Line struct {
	SKU      string
	Quantity int
}

type Quote struct {
	SubtotalCents int64
	DiscountCents int64
	ShippingCents int64
	TotalCents    int64
}

func (q Quote) PayableMerchandiseCents() int64 {
	return q.SubtotalCents - q.DiscountCents
}

func (q Quote) SavedCents() int64 { return q.DiscountCents }

func (q Quote) ValidateTotals() bool {
	return q.TotalCents == q.PayableMerchandiseCents()+q.ShippingCents
}

type Engine struct {
	catalog  *catalog.Catalog
	promos   *promo.Engine
	shipping *shipping.Calculator
}

func NewEngine(c *catalog.Catalog, p *promo.Engine, s *shipping.Calculator) (*Engine, error) {
	if c == nil || p == nil || s == nil {
		return nil, ErrNilDependency
	}
	return &Engine{catalog: c, promos: p, shipping: s}, nil
}

// Quote computes merchandise subtotal, promotion discount and delivery fee.
// Shipping eligibility must use the discounted, payable merchandise amount.
func (e *Engine) Quote(lines []Line, promotionCode string, zone shipping.Zone) (Quote, error) {
	subtotal, err := e.sumLines(lines)
	if err != nil {
		return Quote{}, err
	}
	discount, err := e.calculateDiscount(promotionCode, subtotal)
	if err != nil {
		return Quote{}, err
	}
	return e.completeQuote(subtotal, discount, zone)
}

func (e *Engine) sumLines(lines []Line) (int64, error) {
	if len(lines) == 0 || len(lines) > 20 {
		return 0, ErrInvalidLines
	}
	var subtotal int64
	for _, line := range lines {
		if line.Quantity <= 0 || line.Quantity > 100 {
			return 0, ErrInvalidLines
		}
		item, err := e.catalog.Lookup(line.SKU)
		if err != nil {
			return 0, err
		}
		subtotal += item.PriceCents * int64(line.Quantity)
	}
	return subtotal, nil
}

func (e *Engine) calculateDiscount(code string, subtotal int64) (int64, error) {
	return e.promos.Discount(code, subtotal)
}

func (e *Engine) completeQuote(subtotal, discount int64, zone shipping.Zone) (Quote, error) {
	payable := subtotal - discount
	if payable < 0 {
		return Quote{}, ErrInvalidLines
	}

	// Delivery is calculated separately from promotion eligibility.
	fee, err := e.shipping.Fee(subtotal, zone)
	if err != nil {
		return Quote{}, err
	}
	return Quote{
		SubtotalCents: subtotal,
		DiscountCents: discount,
		ShippingCents: fee,
		TotalCents:    payable + fee,
	}, nil
}
