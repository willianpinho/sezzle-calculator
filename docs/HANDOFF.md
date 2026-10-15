# Maintenance handoff

Start with [README.md](README.md) for the current document map. Read
`docs/SPEC.md` for calculator semantics and HTTP behavior,
`docs/NUMERIC-VALIDATION.md` for binary64 boundaries, and
`docs/ARCHITECTURE.md` for the topology and request flow. Reconstruct the
implementation from these focused entry points:

- `backend/internal/calculator/calculator.go`: pure arithmetic and domain rules.
- `backend/internal/httpapi/api.go`: strict JSON parsing, routing and errors.
- `backend/cmd/server/main.go`: startup, static files, timeouts and shutdown.
- `frontend/src/App.tsx`: form, request state and accessible result/error UI.
- `frontend/src/api.ts`: API payload validation, cancellation and timeout.
- `frontend/e2e/`: real full-stack browser and API acceptance tests.
- `frontend/documentation/` and `frontend/scripts/build-api-docs.mjs`: bundled Swagger UI.
- `Dockerfile`, `compose.yaml`, `.github/workflows/ci.yml`: delivery and CI.

Use the README setup instructions and `make lint`, `make typecheck`, `make test`,
`make build`, `make e2e` to verify a change. Install Chromium, Firefox, and
WebKit first. `/api/openapi.json` is served by Go and `/docs/` is the bundled
browser viewer; do not let generated Swagger assets become source of truth.
Keep the Go API authoritative, preserve the single-process production topology,
and use pnpm's frozen lockfile. Do not add persistent state for arithmetic.

`docs/WORKLOG.md` records decisions and resolved regressions. `docs/TASKS.md`
tracks the delivery phase. `docs/VERIFICATION.md` and `docs/COVERAGE.md`
contain current measured evidence with clearly marked historical sections;
re-run changed behavior rather than treating an older receipt as proof.
