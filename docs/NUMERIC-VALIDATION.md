# Numeric validation matrix

## Scope and acceptance boundary

This is a bounded audit of the seven calculator operations under the stated
finite IEEE-754 binary64 contract. It provides deterministic fixtures and
properties for Go unit/API tests; it is not a proof over all real numbers,
every binary64 pair, or discrete/arbitrary-precision mathematics. The public
API accepts JSON numeric values only, while the pure calculator must reject
non-finite `float64` inputs in either argument as well so its direct callers
cannot bypass the same domain boundary. This includes `b` for `sqrt`, even
though that operation's public request shape rejects `b`.

The numerical authority is Go's standard library:
[`strconv.ParseFloat`](https://pkg.go.dev/strconv@go1.27.2#ParseFloat) returns the nearest binary64
value using IEEE-754 unbiased rounding and returns `±Inf` with `ErrRange` for
out-of-range values. [`math.Pow`](https://pkg.go.dev/math@go1.27.2#Pow) and
[`math.Sqrt`](https://pkg.go.dev/math@go1.27.2#Sqrt) document their special
cases. Go publishes `MaxFloat64` and `SmallestNonzeroFloat64` in
[`math`](https://pkg.go.dev/math@go1.27.2#pkg-constants). The math package does
not promise bit-identical results across architectures, so tests assert exact
values only for exact/reproducible fixtures and otherwise assert finiteness,
sign, or a documented tolerance.

## Required contract outcomes

| Area                      | Fixtures / property                                                                                                     | Expected public result                                                                                                                                                     |
| ------------------------- | ----------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Finite input boundary     | `1e400`, `-1e400`; direct `NaN`, `+Inf`, `-Inf` in each `Calculate` argument                                            | HTTP: `422 validation_error`; pure calculator: `ErrDomain`.                                                                                                                |
| Signed zero               | `add(-0,-0)`, `multiply(-0,2)`, `divide(-0,2)`, `sqrt(-0)`                                                              | Success; preserve the binary64 negative-zero sign where Go defines it. JSON may represent it as `-0`.                                                                      |
| Zero divisor              | `0/0`, `1/+0`, `1/-0`, `-1/+0`, `-1/-0`                                                                                 | `ErrDomain` / HTTP `422 domain_error`, regardless of the infinity sign Go would otherwise produce.                                                                         |
| Zero numerator            | `+0/2`, `-0/2`, `0/-2`                                                                                                  | Success with the binary64 zero sign.                                                                                                                                       |
| Addition/subtraction      | normal `2+3`, cancellation `1 + -1`, overflow `MaxFloat64 + MaxFloat64`, `MaxFloat64 - -MaxFloat64`                     | Normal/cancellation finite success; overflow is `ErrDomain`.                                                                                                               |
| Multiplication            | `-2*3`, `MaxFloat64*2`, `SmallestNonzeroFloat64*0.5`                                                                    | Normal success; overflow `ErrDomain`; finite underflow to signed zero succeeds.                                                                                            |
| Division                  | `7/2`, `0/2`, `MaxFloat64/0.5`, `SmallestNonzeroFloat64/2`                                                              | Normal/zero-numerator success; overflow `ErrDomain`; underflow-to-zero success.                                                                                            |
| Power conventions         | `0^0`, `2^-2`, `(-2)^3`, `0^-1`, `(-0)^-3`, `MaxFloat64^2`                                                              | `0^0=1`; finite results succeed; zero base with negative exponent and overflow are `ErrDomain`.                                                                            |
| Negative-base power       | `(-8)^3`, `(-8)^0.5`, `(-8)^(1.0/3.0)`                                                                                  | Integral exponent succeeds; any non-integral binary64 exponent is `ErrDomain`, even if a real-number simplification might exist.                                           |
| Square root               | `sqrt(9)`, `sqrt(-0)`, `sqrt(-1)`, `sqrt(SmallestNonzeroFloat64)`                                                       | Finite nonnegative values succeed; `sqrt(-0)` keeps `-0`; negative finite input is `ErrDomain`; no finite input can overflow sqrt.                                         |
| Percentage                | `25% of 80`, `100% of MaxFloat64`, `1e-308% of 1e308`, an actually overflowing finite pair such as `MaxFloat64% of 200` | Normal and avoidable-intermediate cases succeed; a non-finite final result is `ErrDomain`. The evaluation may scale one factor by `100` first.                             |
| Representational rounding | `9007199254740992 + 1`; JSON input `9007199254740993`                                                                   | The first sum remains `9007199254740992`; the decimal JSON token rounds to the nearest binary64. No exact-integer guarantee exists beyond the binary64 precision boundary. |
| Subnormals                | `SmallestNonzeroFloat64`, operations that keep it finite, and operations that round it to zero                          | Subnormal finite values are accepted. Rounding to finite zero is a success, not an error.                                                                                  |

Every successful test must also assert `!math.IsNaN(result)` and
`!math.IsInf(result, 0)`. For negative-zero fixtures, assert
`math.Signbit(result)` rather than comparing only with `== 0`.

## Domain classification

`domain_error` intentionally combines two categories under the existing stable
error envelope:

- **Undefined or non-real request:** division by either zero, zero to a
  negative power, negative square root, and a negative base with a
  non-integral binary64 exponent.
- **Unrepresentable finite-contract result:** a finite-real calculation that
  overflows to infinity, plus every other NaN/Infinity result.

Finite underflow is different: it produces a representable zero and therefore
is successful. This contract does not introduce decimal, rational, complex, or
arbitrary-precision arithmetic.

## Resolved audit findings

The following were findings when this matrix was written, not current TODOs.

1. **Pure calculator input invariant — resolved.** `calculator.Calculate` now
   guards `NaN` and infinity before evaluation, preventing special cases such
   as `math.Pow(NaN, 0)` from returning a finite success. The direct-input
   fixtures are in [numeric_test.go](../backend/internal/calculator/numeric_test.go).
2. **Boundary coverage — resolved as bounded evidence.** The matrix fixtures
   cover signed zero, zero-divisor signs, zero negative powers, approximate
   one-third exponents, overflow/underflow, subnormals, percentage final
   overflow, and the `2^53` boundary in
   [numeric_test.go](../backend/internal/calculator/numeric_test.go) and
   [numeric_test.go](../backend/internal/httpapi/numeric_test.go). The fuzz
   invariant is in [fuzz_test.go](../backend/internal/calculator/fuzz_test.go).
   These tests are intentionally bounded; they do not claim exhaustive coverage.

## Implemented validation commands

The targeted guard and tests are implemented. Re-run these commands after a
numeric change:

```bash
(
  cd backend
  go test ./internal/calculator ./internal/httpapi -count=1
  go test -race ./...
  go vet ./...
)
make test
```

The final review should additionally inspect that each HTTP boundary fixture
uses the documented status/error code and that successful JSON results remain
finite. The implementation links are in the resolved findings above; the
current backend coverage status is tracked in [COVERAGE.md](COVERAGE.md).
