
package buggy

import "testing"

func TestTransfer(t *testing.T) {
	a := 1
	b := 2
	want := a + b
	got := transfer(a, b)
	if got != want {
		t.Errorf("got %d but want %d", got, want)
	}
}
