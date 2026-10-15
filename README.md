# Sezzle Calculator

A responsive calculator built with React, TypeScript, and Vite, backed by a Go
standard-library JSON API. It supports addition, subtraction, multiplication,
division, power, square root, and percentage (`a% of b`).

## Run locally

Prerequisites: Node 22.23.2, pnpm 10.33.0, and Go 1.27.2.

```bash
make install
make api                 # terminal 1: API on http://localhost:8080
cd frontend && pnpm dev  # terminal 2: Vite development UI
```

For the production-shaped application, which has Go serve the compiled React
files and API on one port:

```bash
make docker-up
# open http://localhost:8080
```

Stop it with `make docker-down`.

To run the production build without Docker, build the frontend first and then
start Go; it serves `frontend/dist` by default:

```bash
make build
make api
```

## Validate

```bash
make lint
make typecheck
make test
make build
(cd frontend && pnpm exec playwright install chromium firefox)
make e2e
make docker-build
```

`make test` runs race-enabled Go tests with coverage and frontend coverage
tests. Open Go's report with `(cd backend && go tool cover -html=coverage.out)`;
open the frontend report at `frontend/coverage/index.html`. `make e2e` exercises
the production build through Playwright after installing its browsers. CI also
builds and runs the container, requests the health endpoint, posts a real
calculation, loads the UI, and publishes coverage and Playwright artifacts.

## API

`POST /api/calculate` accepts JSON. Binary operations need both operands;
`sqrt` only accepts `a`.

```bash
curl --request POST http://localhost:8080/api/calculate \
  --header 'Content-Type: application/json' \
  --data '{"operation":"add","a":2,"b":3}'
# {"result":5}

curl --request POST http://localhost:8080/api/calculate \
  --header 'Content-Type: application/json' \
  --data '{"operation":"sqrt","a":9}'
# {"result":3}

curl --request POST http://localhost:8080/api/calculate \
  --header 'Content-Type: application/json' \
  --data '{"operation":"divide","a":1,"b":0}'
# {"error":{"code":"domain_error","message":"calculation has no finite real result"}}
```

| Operation   | Request example      | Result |
| ----------- | -------------------- | ------ |
| Add         | `add`, 2, 3          | 5      |
| Subtract    | `subtract`, 2, 3     | -1     |
| Multiply    | `multiply`, -2, 3    | -6     |
| Divide      | `divide`, 7, 2       | 3.5    |
| Power       | `power`, 2, 8        | 256    |
| Square root | `sqrt`, 9            | 3      |
| Percentage  | `percentage`, 25, 80 | 20     |

Successful requests return `200` with `{"result": number}`. Failures return
an error object with a stable code and a safe message. `GET /api/health`
supports the container health check. The API uses IEEE-754 binary64 numbers,
so `0.1 + 0.2` can produce `0.30000000000000004`; this is expected. `0^0`
returns `1`, and percentage means `a%` of `b`.

The complete request rules, numeric edge cases, and error statuses are in
[docs/SPEC.md](docs/SPEC.md).

## Architecture and tradeoffs

The browser calls the Go API on the same origin. In production, Go serves the
compiled Vite files from `STATIC_DIR`; the Docker image therefore has one
network endpoint and no CORS configuration. Compose binds it to loopback for
local use. The service is stateless: there is
no database, cache, login, queue, background job, or service split because the
assignment has one synchronous, side-effect-free calculation request.

Validation and arithmetic remain authoritative on the server. Inputs and
results use finite IEEE-754 binary64 numbers so that Go and JavaScript share a
simple numeric model. Details, including production-hazard decisions, are in
[docs/DECISIONS.md](docs/DECISIONS.md). The testing scope is documented in
[docs/TEST-PLAN.md](docs/TEST-PLAN.md), and the sanitized assignment/design
prompts and specialist receipts are recorded in [docs/PROMPTS.md](docs/PROMPTS.md),
[docs/BACKEND-RECEIPT.md](docs/BACKEND-RECEIPT.md), and
[docs/FRONTEND-RECEIPT.md](docs/FRONTEND-RECEIPT.md), and
[docs/QA-RECEIPT.md](docs/QA-RECEIPT.md).

## Verification evidence

See the measured [coverage report](docs/COVERAGE.md),
[acceptance verification](docs/VERIFICATION.md), and
[delivery ledger](docs/TASKS.md) for the tests, review status, and limitations.

Browser tests reserve local port `18080` and fail if it is occupied; they never
reuse an existing server, so a stale process cannot silently satisfy the tests.
