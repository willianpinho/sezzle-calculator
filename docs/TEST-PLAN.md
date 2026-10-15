# Test plan

## Acceptance contract

`POST /api/calculate` accepts JSON `{ "operation": string, "a": number, "b"?: number }` and returns `{ "result": number }` with HTTP 200. The supported operations are `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, and `percentage`; `percentage` means `a%` of `b`.

Invalid JSON, a missing or unsupported operation, a missing required operand, non-numeric operands, division by zero, and square root of a negative number return a non-2xx JSON error without a result field. The final API wording and status codes must match `docs/SPEC.md`.

## Automated test matrix

| Layer                     | Cases                                                                                                    | Machine-checkable acceptance                                                                                                                                                                                         |
| ------------------------- | -------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Go domain unit tests      | add/subtract/multiply with positive, negative, decimal, and zero values; divide; power; sqrt; percentage | Every listed operation returns its exact expected result. Invalid domain inputs produce the specified typed error. `go test ./... -cover` exits 0.                                                                   |
| Go HTTP integration tests | Valid requests for all seven operations; malformed JSON; validation errors; response `Content-Type`      | Requests exercise the real handler and calculator service without a mocked API. Each valid request is HTTP 200 and response JSON is exactly `{result}`; invalid input has a non-2xx JSON error.                      |
| React unit tests          | Operand controls, operation selection, submit state, rendered result, API error, client/network error    | `pnpm test -- --run` exits 0. Assertions observe visible UI or request construction, never implementation-only state.                                                                                                |
| Browser E2E               | Desktop Chromium, mobile Chromium, Firefox where supported                                               | `pnpm exec playwright test` exits 0 against the Go server and production frontend build. Normal calculation traffic uses the real `/api/calculate` route; only the explicit network-failure case intercepts traffic. |
| Container smoke test      | Build image, start container, request API, load UI                                                       | `docker build` succeeds; a running container serves the UI and an in-container or host request to `/api/calculate` succeeds.                                                                                         |

## API data table

| Operation  | Request operands | Expected result |
| ---------- | ---------------- | --------------- |
| add        | `2`, `3`         | `5`             |
| subtract   | `2`, `3`         | `-1`            |
| multiply   | `-2`, `3`        | `-6`            |
| divide     | `7`, `2`         | `3.5`           |
| power      | `2`, `8`         | `256`           |
| sqrt       | `9`              | `3`             |
| percentage | `25`, `80`       | `20`            |

The Go unit suite also covers zero values, negative bases and exponents where valid, decimal operands, `0 / nonzero`, `0^0` (which returns `1`), and a sufficiently large finite result. It rejects `b = 0` for division, a negative radicand, missing `b` for binary operations, and supplied `b` for `sqrt`.

## Browser journeys

1. At desktop width, select each operation and submit the values in the API data table. Confirm the displayed answer and that it came from the real backend.
2. For `sqrt`, confirm the second operand is not required or displayed according to the final UI contract; submit `9` and see `3`.
3. Submit an invalid calculation such as division by zero. Confirm a readable error is shown, no stale successful result is presented as current, and the form remains usable.
4. Submit once with the API deliberately unavailable. Confirm a readable network-error state and successful recovery on the next real request.
5. Repeat one valid calculation at a narrow mobile viewport; confirm all controls and result remain visible and usable.

## Test implementation constraints

- Keep test data deterministic. Use numeric assertions with a documented tolerance only when floating-point representation requires it.
- Start the E2E server from the production frontend build served by Go, on the fixed local test port `18080`, which must be free. Playwright must wait for readiness and terminate spawned processes on success or failure.
- Capture Playwright traces, screenshots, and video only on retry/failure. Avoid fixed delays; wait for the result or error visible to the user.
- Test files must fail when an operation is removed, a request is sent to the wrong path, validation is bypassed, or a displayed result is stale.

## Commands and evidence

```bash
(cd backend && go test ./... -cover)
(cd frontend && pnpm test:coverage)
(cd frontend && pnpm test:e2e)
docker build -t sezzle-calculator:test .
node scripts/checks/run-all.mjs --diff
```

The final verification receipt records the commands, pass counts, coverage output for Go and React, the browser projects run, and container image tag. Any unavailable tool is reported as unavailable, never as a pass.
