# Delivery notes

The application ships as one image: Vite produces `frontend/dist`, then the Go
server serves it from `STATIC_DIR` while also exposing the JSON API. The image
runs as an unprivileged user and exposes port 8080.

## Commands

Use `make install`, `make test`, `make lint`, `make typecheck`, `make build`,
and `make e2e` for local validation. Install Chromium, Firefox, and WebKit
before E2E with `(cd frontend && pnpm exec playwright install chromium firefox
webkit)`. `make docker-build` creates the production image; `make docker-up`
starts it with a health check at `/api/health`.

CI repeats frontend lint, type checking, coverage tests, and build; Go vet and
race-enabled coverage tests; Chromium, Firefox, and WebKit browser tests; and
an actual Docker build plus API health/calculation and rendered-UI smoke checks.
It uploads coverage and Playwright artifacts. The workflow passes actionlint.

## Delivery prompts

Implementation followed the assignment constraints: React with TypeScript and
Vite, a Go standard-library API, all seven arithmetic operations, validation,
unit/integration/browser coverage, and an optional Docker delivery. The API is
stateless and Compose binds it to loopback for local use; add proxy-level rate
limiting before exposing it to untrusted traffic.

## Historical local verification

`make lint`, `make typecheck`, `make test`, `make build`, and `make e2e`
passed for the initial delivery. The image passed eight calculation cases and six error regressions,
served HTML and its JavaScript asset, and ran as UID 100.
`docker compose up --build --wait` produced a healthy service and a real
percentage calculation returned 20 for 25% of 80; the service was then stopped.
See [VERIFICATION.md](VERIFICATION.md) and [COVERAGE.md](COVERAGE.md) for that
historical evidence. The numeric, Swagger, and responsive follow-up must record
a fresh complete run before it is described as final.

The follow-up Docker smoke has verified the same-origin OpenAPI JSON,
`/api/docs` redirect, local Swagger assets and notices, and nine representative
numeric boundary outcomes. The follow-up browser run passed 54/54. These local
receipts are recorded in [VERIFICATION.md](VERIFICATION.md); hosted CI and the
final delta review remain pending.
