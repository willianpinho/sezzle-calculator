package calculator

import (
	"errors"
	"math"
	"testing"
)

func FuzzCalculateFiniteInvariant(f *testing.F) {
	for index, seed := range []struct{ a, b uint64 }{
		{0, 0}, {0, 1}, {math.Float64bits(-2), math.Float64bits(3)},
		{math.Float64bits(math.SmallestNonzeroFloat64), 0}, {math.Float64bits(math.MaxFloat64), math.Float64bits(100)},
		{math.Float64bits(math.NaN()), 0}, {0, math.Float64bits(math.Inf(1))},
	} {
		f.Add(uint8(index), seed.a, seed.b)
	}
	f.Fuzz(func(t *testing.T, index uint8, aBits, bBits uint64) {
		operation := [...]string{"add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"}[index%7]
		a, b := math.Float64frombits(aBits), math.Float64frombits(bBits)
		result, err := Calculate(operation, a, b)
		if !finite(a) || !finite(b) {
			if !errors.Is(err, ErrDomain) {
				t.Fatalf("%s accepted non-finite input: %g, %g, err=%v", operation, a, b, err)
			}
			return
		}
		if err != nil && !errors.Is(err, ErrDomain) {
			t.Fatalf("%s error = %v", operation, err)
		}
		if err == nil && !finite(result) {
			t.Fatalf("%s returned non-finite %g", operation, result)
		}
	})
}
