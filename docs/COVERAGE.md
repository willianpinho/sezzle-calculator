# Historical coverage report

Measured on 2026-10-14 using the then-current test suites. This is retained as
baseline evidence, not a claim about the numeric, Swagger, or responsive
follow-up. Re-run the commands below after those changes; the final aggregate
will be recorded in [VERIFICATION.md](VERIFICATION.md). Coverage describes unit
and handler tests; browser and container execution are reported separately.

## Go

| Package               | Statement coverage |
| --------------------- | -----------------: |
| `internal/calculator` |             100.0% |
| `internal/httpapi`    |              92.1% |
| `cmd/server`          |               0.0% |
| All backend packages  |              83.0% |

`cmd/server` is exercised as a running process by the browser and Docker tests,
but those runs are not included in the Go unit coverage profile. The zero is
reported rather than excluded from the aggregate.

```bash
(
  cd backend
  go test -race -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  go tool cover -html=coverage.out
)
```

## React/TypeScript

Eight tests in two files passed. These metrics cover `src/App.tsx`,
`src/api.ts`, and `src/calculator.ts`, the production modules imported by the
unit suite. Bootstrap `src/main.tsx`, CSS, the Swagger HTML/build script, and
vendored Swagger assets are exercised by build, browser, and Docker checks and
are outside this Vitest coverage report. The configured thresholds are 85%
statements, lines and functions, and 80% branches.

| Metric     | Covered / total | Coverage |
| ---------- | --------------: | -------: |
| Statements |         76 / 76 |     100% |
| Lines      |         72 / 72 |     100% |
| Functions  |         18 / 18 |     100% |
| Branches   |         54 / 63 |   85.71% |

```bash
(
  cd frontend
  pnpm test:coverage
)
```

Reports: `frontend/coverage/index.html`, `frontend/coverage/lcov.info`, and
`frontend/coverage/coverage-summary.json`. Generated files are ignored by Git;
CI uploads backend and frontend coverage artifacts on each run.

## Current C1/C2 measurement

C1 passes with 100.0% Go statements for every owned production package:
`cmd/server`, `internal/calculator`, and `internal/httpapi`, as well as 100.0%
total (186/186 statements in the hosted artifact). It includes server lifecycle
and startup-failure behavior, not a package exclusion.

```bash
(
  cd backend
  go test -race -count=1 -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out | awk '/^total:/{seen=1; if ($3 != "100.0%") exit 1} END {if (!seen) exit 1}'
  go vet ./...
)
```

C2 passes 22 tests in 5 files and covers 6 handwritten runtime files at 100%:
99/99 statements, 72/72 branches, 20/20 functions, and 95/95 lines. Its global
V8 threshold applies to `src/` and `scripts/build-api-docs.mjs`; the details
are in [FRONTEND-COVERAGE-RECEIPT.md](FRONTEND-COVERAGE-RECEIPT.md). C3
independently repeated the unit, regression, browser, and Docker checks and
approved the final code/test/CI delta; its receipt is
[FOLLOWUP-CODE-REVIEW.md](FOLLOWUP-CODE-REVIEW.md).

## Coverage-extension status

The requested 100% production-code target is measured and enforced without
configuration exclusions or assertion-only tests. Historical measurements above
remain dated baseline evidence; rerun the guarded commands after a change.
