package promo

import "testing"

func TestDiscountAndCap(t *testing.T) {
	e, err := NewEngine([]Rule{{Code: "SAVE25", Percent: 25, MinSubtotalCents: 5000, MaxDiscountCents: 4000}})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ subtotal, want int64 }{{4000, 0}, {12000, 3000}, {20000, 4000}} {
		got, err := e.Discount("SAVE25", tc.subtotal)
		if err != nil || got != tc.want {
			t.Fatalf("Discount(%d) = %d, %v, want %d", tc.subtotal, got, err, tc.want)
		}
	}
}

func TestEmptyCode(t *testing.T) {
	e, err := NewEngine(nil)
	if err != nil {
		t.Fatal(err)
	}
	got, err := e.Discount("", 1234)
	if err != nil || got != 0 {
		t.Fatalf("Discount(empty) = %d, %v", got, err)
	}
}
