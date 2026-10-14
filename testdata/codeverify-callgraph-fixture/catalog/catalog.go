// Package catalog stores items available for purchase.
package catalog

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrInvalidProduct = errors.New("invalid catalog product")
	ErrDuplicateSKU   = errors.New("duplicate SKU")
	ErrUnknownSKU     = errors.New("unknown SKU")
	ErrUnavailable    = errors.New("product unavailable")
)

// Product uses integer cents to avoid floating-point pricing errors.
type Product struct {
	SKU        string
	Title      string
	PriceCents int64
	Active     bool
}

// Validate checks the basic product definition.
func (p Product) Validate() error {
	if strings.TrimSpace(p.SKU) == "" || p.PriceCents <= 0 || p.PriceCents > 1_000_000 {
		return ErrInvalidProduct
	}
	return nil
}

func (p Product) Label() string {
	if p.Title != "" {
		return p.Title
	}
	return p.SKU
}

type Catalog struct {
	products map[string]Product
}

func NewCatalog(products []Product) (*Catalog, error) {
	c := &Catalog{products: make(map[string]Product, len(products))}
	for _, p := range products {
		if err := p.Validate(); err != nil {
			return nil, err
		}
		if _, exists := c.products[p.SKU]; exists {
			return nil, ErrDuplicateSKU
		}
		c.products[p.SKU] = p
	}
	return c, nil
}

// Lookup returns only active products.
func (c *Catalog) Lookup(sku string) (Product, error) {
	p, ok := c.products[sku]
	if !ok {
		return Product{}, ErrUnknownSKU
	}
	if !p.Active {
		return Product{}, ErrUnavailable
	}
	return p, nil
}

func (c *Catalog) Count() int { return len(c.products) }

func (c *Catalog) Has(sku string) bool {
	_, ok := c.products[sku]
	return ok
}

// SKUs returns catalog keys in stable order; this is an administrative helper.
func (c *Catalog) SKUs() []string {
	ids := make([]string, 0, len(c.products))
	for id := range c.products {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
