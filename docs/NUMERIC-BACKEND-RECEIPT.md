# Numeric backend receipt

## Scope completed

`calculator.Calculate` now rejects NaN and either infinity in either direct
argument before evaluating an operation. This closes the finite-result bypass
where `math.Pow(NaN, 0)` or `math.Pow(+Inf, -1)` could otherwise return a
finite value. The HTTP API was already protected because JSON parsing rejects
out-of-range numeric tokens; supplied `b` remains invalid for public `sqrt`
requests.

The boundary suite covers all seven operations across signed zero, zero
denominators, zero negative powers, integer and fractional powers, subnormal
and maximum square roots, finite underflow, representational rounding at 2^53,
and finite-input overflow. It verifies all seven operations reject direct NaN,
`+Inf`, and `-Inf` in both arguments. HTTP integration covers the numerical
error envelope for each applicable overflow/domain case, JSON `1e400` and
`-1e400`, plus finite power, subnormal square-root, rounding, and percentage
underflow responses.

`math/big` with 1024-bit precision independently checks a bounded grid of
finite add, subtract, multiply, and divide results. Square-root expectations
use an exact power-of-two identity for the smallest subnormal and a separate
high-precision `big.Float` square root for `MaxFloat64`.

## Fuzz evidence

`FuzzCalculateFiniteInvariant` maps fuzzed indices across every valid
operation and fuzzes raw binary64 operand bits, including NaN and infinities.
It requires non-finite direct input to produce `ErrDomain`, and every finite
successful result to be finite. The 15-second, two-worker run completed 21
baseline inputs and 1,749,433 executions with no failures.

This is bounded evidence, not an exhaustive proof of every binary64 pair or a
promise of cross-architecture bit-identical `math.Pow`/`math.Sqrt` output.
The contract remains finite IEEE-754 binary64 arithmetic; it does not add
decimal, rational, complex, or arbitrary-precision calculator behavior.

## Verification

Executed from `backend/` on 2026-10-15:

```text
go test ./internal/calculator ./internal/httpapi -count=1
go test -race -coverprofile=coverage.out ./...
go vet ./...
go test ./internal/calculator -run=^$ -fuzz=FuzzCalculateFiniteInvariant -fuzztime=15s -parallel=2
git diff --check
```

All commands passed. The deterministic suite has seven calculator tests and
five HTTP API tests, plus one native fuzz target. Race coverage is 100.0% for
`internal/calculator`, 92.1% for `internal/httpapi`, and 83.2% overall because
the process-lifecycle `cmd/server` entrypoint has no direct test.

## Historical C1 coverage plan (2026-10-15)

The remaining uncovered production paths are the real server lifecycle in
`cmd/server` and HTTP parser/router branches not represented by the original
request table. C1 will preserve the process contract (environment defaults,
timeouts, same static delivery, and SIGINT/SIGTERM shutdown) while extracting a
small listener-based serve function. Tests will exercise it with an ephemeral
listener and real HTTP requests, including graceful cancellation and bind
failure. The main entrypoint will be covered by starting it in-process on a
reserved loopback port and delivering SIGTERM, rather than by a direct fake
call. HTTP additions will use malformed decoder boundaries, health methods, and
static fallback/404 behavior. No valid error path or production assertion will
be removed merely to change the coverage number.

## C1 completion (2026-10-15)

The production entrypoint keeps `log.Fatal` for failed startup while delegating
listener creation and serving to `run` and `serve`; the production shutdown
timeout remains ten seconds. Real lifecycle tests use ephemeral listeners and
HTTP readiness checks for graceful shutdown, server close, listen failure, and
a blocked request whose deliberately short test deadline returns
`context.DeadlineExceeded`. Normal SIGTERM and fatal-startup paths run in
bounded child test processes; the startup child must exit 1 and emit the
startup failure message.

The HTTP suite adds health method, nil-static fallback, malformed object-close,
and malformed trailing-value coverage. Two defensive branches were removed as
unreachable under documented `encoding/json` invariants: object member tokens
are strings, and a `RawMessage` emitted after successful strict decoding cannot
fail JSON syntax parsing again. Valid error handling remains intact.

The reproducible Makefile/CI command is:

```text
cd backend && go test -race -count=1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

It passed with 100.0% statement coverage in `cmd/server`,
`internal/calculator`, and `internal/httpapi`, and 100.0% overall. `go vet
./...` and `git diff --check` also passed. No extra coverage tool or Go
dependency is required.
