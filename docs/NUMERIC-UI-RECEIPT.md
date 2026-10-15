# Numeric and responsive browser receipt

Date: 2026-10-15

## Numerical acceptance

`pnpm exec playwright test --project=numeric-chromium` passed **2/2** against
the production frontend served by the Go binary. The requests use the live API;
they do not mock calculation responses.

The fixtures verify documented binary64 behavior: `0^0 = 1`, finite zero and
subnormal outcomes, signed zero preserved by `sqrt(-0)`, finite negative-base
integer power, square-root power, binary64 rounding at `2^53 + 1`, percentage
scaling, and the largest finite percentage result. They verify `422
domain_error` for `0/0`, division by `-0`, `0` to a negative power, negative
base with a fractional exponent, negative square root, and representative
overflow in add, subtract, multiply, divide, power, and percentage. JSON
lexical overflow (`1e400` and `-1e400`) returns `422 validation_error`.

This is a set of documented boundary fixtures, not an exhaustive proof over
real or discrete mathematics. JSON cannot represent `NaN` or infinity, so
browser requests cannot exercise those direct-calculator guards.

## Responsive acceptance

The responsive production-stack test passed **7/7**:

| Project                 | Engine                     | CSS viewport |
| ----------------------- | -------------------------- | ------------ |
| `responsive-320`        | Chromium                   | 320 x 700    |
| `responsive-375`        | Chromium                   | 375 x 760    |
| `responsive-390`        | Chromium                   | 390 x 844    |
| `responsive-768`        | Chromium                   | 768 x 1024   |
| `responsive-1024`       | Chromium                   | 1024 x 768   |
| `responsive-landscape`  | Chromium                   | 844 x 390    |
| `responsive-webkit-390` | WebKit iPhone 13 emulation | 390 x 844    |

At every viewport the test proves no document horizontal overflow, visible and
in-viewport labels, controls, button, result/error region, and result text. It
changes operations, renders the 22-character `1.2345678901234567e+308` result,
checks a small scientific result, and displays a server-side square-root error.

Manual inspection covered the 320px long result and error states, Chromium
landscape error state, and WebKit 390px error state. They showed wrapping
inside the result panel and no clipping or overlap. The recaptured ignored
artifacts are retained for coordinator inspection at:

```text
.claude/nexus/n3-screenshots/320-long-result.png
.claude/nexus/n3-screenshots/320-error-state.png
.claude/nexus/n3-screenshots/landscape-error-state.png
.claude/nexus/n3-screenshots/webkit-390-error-state.png
```

WebKit coverage is browser emulation, not a claim of physical-device testing.
The browser was installed locally with `pnpm exec playwright install webkit`.

## Verification status

The final command `cd frontend && pnpm test:e2e` passed **54/54** after local
Swagger assets were installed and bundled. That run includes the numeric and
responsive cases above, three baseline browser projects (Chromium, Firefox,
and WebKit iPhone emulation), and Swagger UI acceptance in each baseline
project. The Swagger checks load the live OpenAPI document without external
requests, follow the `/api/docs` redirect, and use the UI to send a real
`add(2, 3)` request returning `{ "result": 5 }` and a division-by-zero request
returning HTTP 422 with `domain_error`.
