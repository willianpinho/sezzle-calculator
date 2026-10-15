# QA receipt

Date: 2026-10-14

## Browser acceptance

`cd frontend && pnpm exec playwright test` passed: **36 tests** in **12.5 seconds**.

| Project           | Engine and viewport   | Result    |
| ----------------- | --------------------- | --------- |
| `chromium`        | Desktop Chromium      | 12 passed |
| `mobile-chromium` | Chromium at 320 x 700 | 12 passed |
| `firefox`         | Desktop Firefox       | 12 passed |

The Playwright web server first ran `pnpm build`, then started the Go server
with `PORT=18080` and `STATIC_DIR=../frontend/dist`. Browser calculation tests
waited for actual `POST /api/calculate` responses. The only route mock aborts
that endpoint in the explicit network-failure-and-recovery test.

The suite covers all seven operations, blank client validation, division by
zero, negative square root, overflow, keyboard submission, stale-result
clearing after an edit, network recovery, the health endpoint, malformed JSON,
unknown operations and fields, and an invalid `b` sent with `sqrt`.

Visual smoke screenshots from the production stack were captured and manually
inspected at desktop and 320px mobile widths. They are ignored test artifacts:
`frontend/test-results/desktop-calculator.png` and
`frontend/test-results/mobile-calculator.png`. Both show every control and the
result panel without overlap or clipping.

## Supporting checks

`cd frontend && pnpm typecheck` passed. Playwright Chromium and Firefox
browsers were installed locally with `pnpm exec playwright install chromium
firefox` before the run. The coordinator subsequently repeated lint, types,
unit coverage, Go race tests/vet and all 36 browser tests, then rebuilt and
validated the final Docker image and healthy Compose service. See
[VERIFICATION.md](VERIFICATION.md) for command outcomes and explicit optional
checker limitations.
