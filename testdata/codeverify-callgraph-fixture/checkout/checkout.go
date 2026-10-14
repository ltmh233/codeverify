// Package checkout exposes customer-facing quote operations.
package checkout

import (
	"errors"
	"fmt"

	"example.com/codeverify-callgraph-fixture/pricing"
	"example.com/codeverify-callgraph-fixture/shipping"
)

var (
	ErrNoPricingEngine = errors.New("checkout pricing engine is required")
	ErrEmptyCart       = errors.New("cart must contain at least one item")
)

type Request struct {
	Lines     []pricing.Line
	PromoCode string
	Zone      shipping.Zone
}

type Service struct {
	pricing *pricing.Engine
}

func NewService(engine *pricing.Engine) (*Service, error) {
	if engine == nil {
		return nil, ErrNoPricingEngine
	}
	return &Service{pricing: engine}, nil
}

// Quote returns the amount a customer should pay for an order.
// Domestic free-shipping eligibility is based on merchandise AFTER discounts,
// not the original price. Quote must satisfy this rule for all valid requests.
func (s *Service) Quote(req Request) (pricing.Quote, error) {
	if err := validateRequest(req); err != nil {
		return pricing.Quote{}, err
	}
	return s.pricing.Quote(req.Lines, req.PromoCode, req.Zone)
}

func validateRequest(req Request) error {
	if len(req.Lines) == 0 {
		return ErrEmptyCart
	}
	return nil
}

// Preview is a presentation helper, not a payment operation.
func (s *Service) Preview(req Request) (string, error) {
	q, err := s.Quote(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("total=%d discount=%d shipping=%d", q.TotalCents, q.DiscountCents, q.ShippingCents), nil
}

func (s *Service) IsFreeDelivery(req Request) (bool, error) {
	q, err := s.Quote(req)
	if err != nil {
		return false, err
	}
	return q.ShippingCents == 0, nil
}

func LineCount(req Request) int { return len(req.Lines) }
