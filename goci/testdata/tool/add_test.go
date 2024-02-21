package add

import "testing"

func TestAdd(t *testing.T) {
	a := 2
	b := 3

	expected := 5
	result := add(a, b)

	if expected != result {
		t.Errorf("Expected %d, got %d", expected, result)
	}
}
