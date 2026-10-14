// Package promo applies percentage-based promotional rules.
package promo

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrInvalidRule = errors.New("invalid promotion rule")
	ErrUnknownCode = errors.New("unknown promotion code")
	ErrBadSubtotal = errors.New("subtotal must be non-negative")
)

type Rule struct {
	Code             string
	Percent          int
	MinSubtotalCents int64
	MaxDiscountCents int64 // Zero means no cap.
}

func (r Rule) Validate() error {
	if strings.TrimSpace(r.Code) == "" || r.Percent <= 0 || r.Percent > 80 || r.MinSubtotalCents < 0 || r.MaxDiscountCents < 0 {
		return ErrInvalidRule
	}
	return nil
}

func (r Rule) Eligible(subtotalCents int64) bool {
	return subtotalCents >= r.MinSubtotalCents
}

type Engine struct {
	rules map[string]Rule
}

func NewEngine(rules []Rule) (*Engine, error) {
	e := &Engine{rules: make(map[string]Rule, len(rules))}
	for _, r := range rules {
		if err := r.Validate(); err != nil {
			return nil, err
		}
		if _, duplicate := e.rules[r.Code]; duplicate {
			return nil, ErrInvalidRule
		}
		e.rules[r.Code] = r
	}
	return e, nil
}

func (e *Engine) Lookup(code string) (Rule, error) {
	r, ok := e.rules[code]
	if !ok {
		return Rule{}, ErrUnknownCode
	}
	return r, nil
}

// Discount returns the discount on the merchandise subtotal, in cents.
// An empty code means no promotion.
func (e *Engine) Discount(code string, subtotalCents int64) (int64, error) {
	if subtotalCents < 0 {
		return 0, ErrBadSubtotal
	}
	if code == "" {
		return 0, nil
	}
	r, err := e.Lookup(code)
	if err != nil {
		return 0, err
	}
	if !r.Eligible(subtotalCents) {
		return 0, nil
	}
	return limitDiscount(percentOff(subtotalCents, r.Percent), r.MaxDiscountCents), nil
}

func percentOff(amount int64, percent int) int64 {
	return amount * int64(percent) / 100
}

func limitDiscount(discount, cap int64) int64 {
	if cap > 0 && discount > cap {
		return cap
	}
	return discount
}

// Codes is a utility for displaying configured promotional codes.
func (e *Engine) Codes() []string {
	codes := make([]string, 0, len(e.rules))
	for code := range e.rules {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
