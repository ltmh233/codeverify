package catalog

import (
	"errors"
	"testing"
)

func TestCatalogLookup(t *testing.T) {
	c, err := NewCatalog([]Product{{SKU: "pen", PriceCents: 1500, Active: true}})
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Lookup("pen")
	if err != nil || got.PriceCents != 1500 {
		t.Fatalf("Lookup(pen) = %+v, %v", got, err)
	}
	if _, err := c.Lookup("missing"); !errors.Is(err, ErrUnknownSKU) {
		t.Fatalf("unknown SKU error = %v", err)
	}
}

func TestDuplicateSKU(t *testing.T) {
	_, err := NewCatalog([]Product{
		{SKU: "pen", PriceCents: 1500, Active: true},
		{SKU: "pen", PriceCents: 1500, Active: true},
	})
	if !errors.Is(err, ErrDuplicateSKU) {
		t.Fatalf("duplicate SKU error = %v", err)
	}
}
