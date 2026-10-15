# Backend implementation receipt

## Implemented contract

The Go module is `github.com/willianpinho/sezzle-calculator/backend`, using the
standard library only. `cmd/server` starts one bounded-timeout HTTP server,
serves `STATIC_DIR` (default `../frontend/dist`), and shuts down gracefully on
SIGINT or SIGTERM. It owns `POST /api/calculate` and `GET /api/health`; unknown
`/api` paths return the documented JSON 404 and never fall through to static
delivery.

`internal/calculator` contains pure IEEE-754 arithmetic. `internal/httpapi`
strictly decodes one JSON object, detects duplicate fields, rejects unexpected
fields and invalid operands, and produces only the contract's JSON errors.
The body limit is 1 MiB. Percentage scales the larger magnitude operand before
multiplication, preventing avoidable intermediate overflow or underflow.

## Sanitized implementation brief

Implement a stateless Go calculator JSON API with seven operations, strict JSON
shape and media-type validation, finite-result rules, safe JSON errors, static
frontend delivery, graceful shutdown, and unit plus HTTP integration tests.
Use only the Go standard library and do not include credentials or external
services.

## Hazard decisions

The endpoint is pure and has no mutable shared state: retries are idempotent;
there are no read-modify-write races. There is no database, cache, queue,
migration, connection pool, dual write, outbox, consumer retry, DLQ, or stream.
The request is bounded to 1 MiB and server read/write/header/idle timeouts are
set. Application rate limiting is intentionally deployment-layer work for this
small anonymous demo, as recorded in ADR-004. No secret comparison occurs;
parser token/error comparisons are ordinary comparisons over public JSON syntax.

## Verification

Run from `backend/`:

```text
gofmt -w cmd internal
go test -race -coverprofile=coverage.out ./...
go vet ./...
git diff --check
```

All commands passed on 2026-10-14. The race-enabled suite covers calculator
math (integer and negative powers, decimal operands, zero-numerator division,
both percentage scaling paths, finite bounds) and real `httptest` handler paths. HTTP coverage includes
malformed/non-object/multiple and duplicate JSON, case sensitivity, content
type, missing/non-string operation, payload limit, operand shape, domain limits, methods, API namespace
routing, health, and static handoff.
