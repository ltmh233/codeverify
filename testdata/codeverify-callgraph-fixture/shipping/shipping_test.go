package shipping

import "testing"

func TestDomesticShippingThreshold(t *testing.T) {
	c, err := NewCalculator(1200, 2500, 10000)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ payable, want int64 }{{9000, 1200}, {9999, 1200}, {10000, 0}, {11000, 0}} {
		got, err := c.Fee(tc.payable, Domestic)
		if err != nil || got != tc.want {
			t.Fatalf("Fee(%d,domestic) = %d, %v, want %d", tc.payable, got, err, tc.want)
		}
	}
}

func TestInternationalFlatFee(t *testing.T) {
	c, err := NewCalculator(1200, 2500, 10000)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.Fee(50000, International)
	if err != nil || got != 2500 {
		t.Fatalf("Fee(international) = %d, %v", got, err)
	}
}
