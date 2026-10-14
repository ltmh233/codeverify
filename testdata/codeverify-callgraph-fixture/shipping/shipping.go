// Package shipping calculates shipping fees for different regions.
package shipping

import "errors"

var (
	ErrInvalidRates  = errors.New("invalid shipping rates")
	ErrUnknownZone   = errors.New("unknown shipping zone")
	ErrInvalidAmount = errors.New("negative merchandise amount")
)

type Zone string

const (
	Domestic      Zone = "domestic"
	International Zone = "international"
)

type Calculator struct {
	domesticFlatCents      int64
	internationalFlatCents int64
	freeThresholdCents     int64
}

func NewCalculator(domesticFlat, internationalFlat, freeThreshold int64) (*Calculator, error) {
	if domesticFlat < 0 || internationalFlat < 0 || freeThreshold <= 0 {
		return nil, ErrInvalidRates
	}
	return &Calculator{
		domesticFlatCents:      domesticFlat,
		internationalFlatCents: internationalFlat,
		freeThresholdCents:     freeThreshold,
	}, nil
}

// Fee computes delivery cost from the PAYABLE merchandise amount,
// AFTER applying all discounts. Domestic shipping is free when that amount
// reaches the configured threshold; international delivery has a flat fee.
func (c *Calculator) Fee(payableMerchandiseCents int64, zone Zone) (int64, error) {
	if payableMerchandiseCents < 0 {
		return 0, ErrInvalidAmount
	}
	switch zone {
	case Domestic:
		return c.domesticFee(payableMerchandiseCents), nil
	case International:
		return c.internationalFee(), nil
	default:
		return 0, ErrUnknownZone
	}
}

func (c *Calculator) domesticFee(payableMerchandiseCents int64) int64 {
	if c.EligibleForFreeDomestic(payableMerchandiseCents) {
		return 0
	}
	return c.domesticFlatCents
}

func (c *Calculator) internationalFee() int64 {
	return c.internationalFlatCents
}

func (c *Calculator) EligibleForFreeDomestic(payableMerchandiseCents int64) bool {
	return payableMerchandiseCents >= c.freeThresholdCents
}

func (c *Calculator) Threshold() int64 { return c.freeThresholdCents }

func SupportedZones() []Zone { return []Zone{Domestic, International} }
