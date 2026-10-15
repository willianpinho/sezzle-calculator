// Package calculator contains the pure arithmetic contract for the API.
package calculator

import (
	"errors"
	"math"
)

var ErrDomain = errors.New("calculation has no finite real result")

// Calculate evaluates an operation using IEEE-754 binary64 arithmetic.
func Calculate(operation string, a, b float64) (float64, error) {
	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return 0, ErrDomain
		}
		result = a / b
	case "power":
		if a < 0 && math.Trunc(b) != b {
			return 0, ErrDomain
		}
		result = math.Pow(a, b)
	case "sqrt":
		if a < 0 {
			return 0, ErrDomain
		}
		result = math.Sqrt(a)
	case "percentage":
		// Scale the larger factor first to avoid an avoidable intermediate overflow.
		if math.Abs(a) >= math.Abs(b) {
			result = (a / 100) * b
		} else {
			result = a * (b / 100)
		}
	default:
		return 0, errors.New("unknown operation")
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrDomain
	}
	return result, nil
}
