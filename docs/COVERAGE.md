# Coverage report

Measured on 2026-10-14 using the checked-in test suites. Coverage describes unit
and handler tests; browser and container execution are reported separately in
[VERIFICATION.md](VERIFICATION.md).

## Go

| Package | Statement coverage |
| --- | ---: |
| `internal/calculator` | 100.0% |
| `internal/httpapi` | 92.1% |
| `cmd/server` | 0.0% |
| All backend packages | 83.0% |

`cmd/server` is exercised as a running process by the browser and Docker tests,
but those runs are not included in the Go unit coverage profile. The zero is
reported rather than excluded from the aggregate.

```bash
cd backend
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

## React/TypeScript

Eight tests in two files passed. These metrics cover `src/App.tsx`,
`src/api.ts`, and `src/calculator.ts`, the production modules imported by the
unit suite. Bootstrap `src/main.tsx` and CSS are exercised through browser
tests and are outside this unit coverage report. The configured thresholds
are 85% statements, lines and functions, and 80% branches.

| Metric | Covered / total | Coverage |
| --- | ---: | ---: |
| Statements | 76 / 76 | 100% |
| Lines | 72 / 72 | 100% |
| Functions | 18 / 18 | 100% |
| Branches | 54 / 63 | 85.71% |

```bash
cd frontend
pnpm test:coverage
```

Reports: `frontend/coverage/index.html`, `frontend/coverage/lcov.info`, and
`frontend/coverage/coverage-summary.json`. Generated files are ignored by Git;
CI uploads backend and frontend coverage artifacts on each run.
