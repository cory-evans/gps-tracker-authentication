package service_test

import "testing"

func TestFunc(t *testing.T) {
	a := 1 + 1

	if a != 2 {
		t.Errorf("Expected 2 got: %d", a)
	}
}
