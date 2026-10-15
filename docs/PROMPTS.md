# Prompt record

Tooling: OpenAI Codex, with delegated implementation and independent review
agents. AI-assisted changes were checked with unit, integration, browser and
container tests. This record preserves the assignment intent and sanitized
task prompts; it contains no credentials or private configuration.

## Sanitized assignment

Build a complete Sezzle calculator take-home in roughly 2–4 hours. A Go backend
is required, with a React TypeScript frontend. Support add, subtract, multiply,
divide, power, square root, and percentage. The UI must be intuitive,
responsive, and accessible. Provide a JSON API with validation/errors; unit
tests and coverage reports for both layers; integration and E2E tests; Docker
full-stack setup; README instructions and design decisions; and a GitHub-ready
repository. Use English. Prefer correctness and simplicity.

## Actual design instruction

Create a concise specification before implementation. Choose the simplest
stateless standard-library Go API and React Vite UI with no database. Define
`POST /api/calculate` as `{operation,a,b?}` to `{result:number}`; enumerate
operations as add, subtract, multiply, divide, power, sqrt, percentage; define
percentage as a percent of b; reject `b` for square root; use strict finite
binary64 real-number semantics; and make policies for `0^0` and overflow
explicit. Document API errors, architecture rationale, hazard applicability,
handoff context, and acceptance checks. Use `pnpm` for Node work.

The scoped implementation prompts and their verification receipts are recorded
in [BACKEND-RECEIPT.md](BACKEND-RECEIPT.md),
[FRONTEND-RECEIPT.md](FRONTEND-RECEIPT.md), and
[QA-RECEIPT.md](QA-RECEIPT.md).

## Implementation and verification prompts

The following are sanitized task prompts: machine-specific paths were replaced
with repository-relative ownership, and repeated execution-policy text was
omitted. They describe the AI work performed; they are not a claim that the
project was written without AI assistance.

### Frontend

> Implement the React/TypeScript frontend in frontend/. Use seven operations,
> accessible labeled controls, a clear loading/result/error state, and a
> responsive layout. Call POST /api/calculate; omit b for square root. Keep
> arithmetic authoritative in Go. Use pnpm, strict TypeScript, Vitest and
> Testing Library. Add tests for request construction, validation, API errors,
> malformed responses, timeouts, and clearing stale results. Keep source files
> below 300 lines and provide build, lint, typecheck and coverage commands.

### Backend

> Implement the Go standard-library API specified in docs/SPEC.md. Separate
> pure calculator logic from HTTP parsing. Reject invalid/duplicate JSON,
> nonnumeric or missing operands, unsupported operations and nonfinite results.
> Serve the compiled frontend from STATIC_DIR and provide /api/health. Add
> table-driven math tests and httptest integration tests, then run gofmt,
> race-enabled coverage tests, and go vet.

### Delivery

> Provide one multi-stage Docker image with Node for the frontend build and
> Go for the backend build, running as a nonroot user. Add Docker Compose,
> Makefile commands, a CI pipeline, and README setup/API/design instructions.
> CI must execute unit coverage, lint/types/build, browser tests, and a real
> container smoke test. Avoid extra services or deployment infrastructure.

### Independent acceptance verification

> Use Playwright against the actual Go server serving the production frontend
> build. Test all seven operations, form/server errors, overflow, keyboard
> submission, stale results, and network failure followed by recovery. Run
> desktop Chromium, mobile Chromium, and Firefox. Mock traffic only for the
> explicit network-error scenario. Save evidence and report failures honestly.

### Regression feedback

> Reject quoted numeric strings: decoding into json.Number alone accepts them.
> For a non-object JSON body, return a nonnil parsing error after writing 400
> so the handler cannot write another JSON response. Decode the entire response
> in tests to detect that regression. Test oversized leading whitespace and
> ensure percentage calculations avoid unnecessary intermediate overflow.

> Validate successful and error payloads as unknown objects before reading
> properties. Keep the request timeout active through response-body parsing,
> not just until response headers arrive.

Review prompts and final receipts are linked from the delivery ledger.
