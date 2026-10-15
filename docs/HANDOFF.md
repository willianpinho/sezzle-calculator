# Maintenance handoff

Read `docs/SPEC.md` for calculator semantics and HTTP behavior, then
`docs/DECISIONS.md` for architecture and numeric assumptions. Reconstruct the
implementation from these focused entry points:

- `backend/internal/calculator/calculator.go`: pure arithmetic and domain rules.
- `backend/internal/httpapi/api.go`: strict JSON parsing, routing and errors.
- `backend/cmd/server/main.go`: startup, static files, timeouts and shutdown.
- `frontend/src/App.tsx`: form, request state and accessible result/error UI.
- `frontend/src/api.ts`: API payload validation, cancellation and timeout.
- `frontend/e2e/`: real full-stack browser and API acceptance tests.
- `Dockerfile`, `compose.yaml`, `.github/workflows/ci.yml`: delivery and CI.

Use the README setup instructions and `make lint`, `make typecheck`, `make test`,
`make build`, `make e2e` to verify a change. Install Playwright browsers first.
Keep the Go API authoritative, preserve the single-process production topology,
and use pnpm's frozen lockfile. Do not add persistent state for arithmetic.

`docs/WORKLOG.md` records decisions and resolved regressions. `docs/TASKS.md`
tracks the delivery phase. `docs/VERIFICATION.md` and `docs/COVERAGE.md` contain
reproducible evidence and limitations; review receipts document independent
feedback. Re-run changed behavior rather than treating an old receipt as proof.
