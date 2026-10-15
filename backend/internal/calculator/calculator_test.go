package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name              string
		op                string
		a, b, want, delta float64
	}{
		{"add", "add", 2, 3, 5, 0}, {"subtract", "subtract", 2, 3, -1, 0},
		{"multiply", "multiply", -2, 3, -6, 0}, {"divide", "divide", 7, 2, 3.5, 0},
		{"divide zero numerator", "divide", 0, 2, 0, 0},
		{"power", "power", 2, 8, 256, 0}, {"zero power zero", "power", 0, 0, 1, 0},
		{"negative base integer exponent", "power", -2, 3, -8, 0},
		{"negative exponent", "power", 2, -2, .25, 0},
		{"decimal addition", "add", .1, .2, .3, 1e-15},
		{"sqrt", "sqrt", 9, 0, 3, 0}, {"percentage", "percentage", 25, 80, 20, 0},
		{"percentage scales a", "percentage", 200, 20, 40, 0},
		{"percentage avoids intermediate overflow", "percentage", 100, math.MaxFloat64, math.MaxFloat64, 0},
		{"percentage avoids intermediate underflow", "percentage", 1e-308, 1e308, .01, 1e-15},
		{"underflow is finite", "multiply", math.SmallestNonzeroFloat64, 0.5, 0, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Calculate(test.op, test.a, test.b)
			if err != nil || math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-test.want) > test.delta {
				t.Fatalf("Calculate() = %v, %v; want %v, nil", got, err, test.want)
			}
		})
	}
}

func TestCalculateUnknownOperation(t *testing.T) {
	if _, err := Calculate("unknown", 1, 2); err == nil {
		t.Fatal("Calculate() error = nil; want an error")
	}
}

func TestCalculateDomainErrors(t *testing.T) {
	for _, test := range []struct {
		op   string
		a, b float64
	}{
		{"divide", 1, 0}, {"sqrt", -1, 0}, {"power", -1, .5}, {"power", math.MaxFloat64, 2},
	} {
		if _, err := Calculate(test.op, test.a, test.b); !errors.Is(err, ErrDomain) {
			t.Errorf("Calculate(%q) error = %v; want domain error", test.op, err)
		}
	}
}
