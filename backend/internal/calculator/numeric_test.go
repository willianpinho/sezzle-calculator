package calculator

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

func TestNumericBoundaries(t *testing.T) {
	for _, test := range []struct {
		name, op          string
		a, b, want, delta float64
	}{
		{"positive fractional power", "power", 9, .5, 3, 0},
		{"negative base even exponent", "power", -2, 4, 16, 0},
		{"negative base odd negative exponent", "power", -2, -3, -.125, 0},
		{"sqrt smallest subnormal", "sqrt", math.SmallestNonzeroFloat64, 0, math.Ldexp(1, -537), 0},
		{"sqrt largest finite", "sqrt", math.MaxFloat64, 0, bigSqrt(math.MaxFloat64), 0},
		{"multiply underflow", "multiply", math.SmallestNonzeroFloat64, .5, 0, 0},
		{"divide underflow", "divide", math.SmallestNonzeroFloat64, 2, 0, 0},
		{"power underflow", "power", math.SmallestNonzeroFloat64, 2, 0, 0},
		{"percentage underflow", "percentage", math.SmallestNonzeroFloat64, 50, 0, 0},
		{"two to fifty-three rounding", "add", 1 << 53, 1, 1 << 53, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := Calculate(test.op, test.a, test.b)
			if err != nil || !finite(got) || math.Abs(got-test.want) > test.delta {
				t.Fatalf("Calculate() = %g, %v; want %g", got, err, test.want)
			}
		})
	}
}

func TestDomainBoundaries(t *testing.T) {
	for _, test := range []struct {
		name, op string
		a, b     float64
	}{
		{"positive zero denominator", "divide", 1, 0}, {"negative zero denominator", "divide", 1, math.Copysign(0, -1)},
		{"zero over zero", "divide", 0, 0}, {"zero negative power", "power", 0, -2},
		{"negative zero negative power", "power", math.Copysign(0, -1), -3},
		{"negative fractional third", "power", -8, 1.0 / 3}, {"add overflow", "add", math.MaxFloat64, math.MaxFloat64},
		{"subtract overflow", "subtract", math.MaxFloat64, -math.MaxFloat64}, {"multiply overflow", "multiply", math.MaxFloat64, 2},
		{"divide overflow", "divide", math.MaxFloat64, math.SmallestNonzeroFloat64}, {"power overflow", "power", math.MaxFloat64, 2},
		{"percentage overflow", "percentage", math.MaxFloat64, 200},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Calculate(test.op, test.a, test.b); !errors.Is(err, ErrDomain) {
				t.Fatalf("error = %v; want ErrDomain", err)
			}
		})
	}
}

func TestSignedZeroAndNonFiniteInputs(t *testing.T) {
	got, err := Calculate("add", math.Copysign(0, -1), math.Copysign(0, -1))
	if err != nil || got != 0 || !math.Signbit(got) {
		t.Fatalf("-0 + -0 = %g, %v; want -0", got, err)
	}
	got, err = Calculate("multiply", math.Copysign(0, -1), 2)
	if err != nil || got != 0 || !math.Signbit(got) {
		t.Fatalf("-0 * 2 = %g, %v; want -0", got, err)
	}
	got, err = Calculate("divide", math.Copysign(0, -1), 2)
	if err != nil || got != 0 || !math.Signbit(got) {
		t.Fatalf("-0 / 2 = %g, %v; want -0", got, err)
	}
	got, err = Calculate("sqrt", math.Copysign(0, -1), 0)
	if err != nil || got != 0 || !math.Signbit(got) {
		t.Fatalf("sqrt(-0) = %g, %v; want -0", got, err)
	}
	for _, operation := range []string{"add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"} {
		for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := Calculate(operation, value, 1); !errors.Is(err, ErrDomain) {
				t.Errorf("%s(%g, 1) error = %v", operation, value, err)
			}
			if _, err := Calculate(operation, 1, value); !errors.Is(err, ErrDomain) {
				t.Errorf("%s(1, %g) error = %v", operation, value, err)
			}
		}
	}
}

func TestBasicArithmeticMatchesBigFloat(t *testing.T) {
	values := []float64{math.SmallestNonzeroFloat64, -3.5, .1, 1 << 53, math.MaxFloat64 / 4}
	for _, op := range []string{"add", "subtract", "multiply", "divide"} {
		for _, a := range values {
			for _, b := range values {
				want := bigResult(op, a, b)
				got, err := Calculate(op, a, b)
				if !finite(want) {
					if !errors.Is(err, ErrDomain) {
						t.Fatalf("%s(%g, %g) = %g, %v; want domain error", op, a, b, got, err)
					}
					continue
				}
				if err != nil || got != want {
					t.Fatalf("%s(%g, %g) = %g, %v; big oracle %g", op, a, b, got, err, want)
				}
			}
		}
	}
}

func bigResult(op string, a, b float64) float64 {
	x := new(big.Float).SetPrec(1024).SetFloat64(a)
	y := new(big.Float).SetPrec(1024).SetFloat64(b)
	switch op {
	case "add":
		x.Add(x, y)
	case "subtract":
		x.Sub(x, y)
	case "multiply":
		x.Mul(x, y)
	case "divide":
		x.Quo(x, y)
	}
	value, _ := x.Float64()
	return value
}
func bigSqrt(value float64) float64 {
	root := new(big.Float).SetPrec(1024).Sqrt(new(big.Float).SetPrec(1024).SetFloat64(value))
	result, _ := root.Float64()
	return result
}
func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
