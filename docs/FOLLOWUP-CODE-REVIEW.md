# Follow-up code review

**Scope:** `test/numeric-boundaries` against `origin/main` on 2026-10-15.
This receipt records the completed frontend, numerical-contract, OpenAPI,
Swagger, server-lifecycle, and final coverage review. It retains the preliminary
finding and its resolution for the review trail.

## Preliminary findings

### Resolved concern — importing the Swagger helper rebuilt real generated assets

The earlier module-evaluation call deleted and recreated the repository's
ignored `public/docs` directory when the test imported its helper. The current
[`frontend/scripts/build-api-docs.mjs:45`](../frontend/scripts/build-api-docs.mjs#L45)
uses an ESM entry-point guard, so imports are side-effect-free, normal CLI
execution still builds the assets, and the no-script-argument case is covered.

### Praise — numeric verification uses independent and bounded evidence

[`backend/internal/calculator/numeric_test.go:85`](../backend/internal/calculator/numeric_test.go#L85)
compares add/subtract/multiply/divide fixtures with high-precision `big.Float`,
and keeps power behavior tied to the documented Go binary64 contract rather than
pretending an oracle proves transcendental behavior. The boundary matrix defines
signed zero, underflow, overflow, non-real powers, and input finiteness without
claiming exhaustive proof.

## Review evidence

- Focused Go numerical, HTTP-boundary, OpenAPI, and redirect tests passed.
- At the preliminary review, frontend coverage, lint, and type-check passed:
  19 tests and 100% of the then-six included production modules on statements,
  branches, functions, and lines. The final count and metrics appear below.
- The OpenAPI document is a valid JSON, single canonical contract. It models the
  unary/binary request distinction with `oneOf`, matches the served routes, and
  the UI uses same-origin, locally bundled Swagger assets with validation
  disabled.
- Manual HTTP verification against a separate local server accepted valid JSON
  whitespace around numeric fields and returned `{"result":5}`.
- Source code files are below 300 lines; the 370-line OpenAPI JSON is the
  canonical machine-readable contract and is not source code subject to that
  limit.

## Final C3 review

The lifecycle suite exercises successful shutdown, listener and startup failure,
closed-server handling, forced shutdown timeout, and the real `log.Fatal` exit
path. The JSON-object parser's string assertion is safe: after `Decoder.More`
within an object, `Decoder.Token` returns an object-member string or an error.

The finite-input guard rejects direct `NaN` and infinity before arithmetic; the
numeric tests use an independent `big.Float` oracle only where it is suitable
and document the remaining binary64 contract as bounded evidence. The embedded
OpenAPI contract, redirect, and same-origin offline Swagger assets match their
routes and browser coverage. No secret, migration, or architecture issue was
introduced.

I independently ran the Go race/coverage/vet command, frontend coverage/lint/
type-check, the production frontend build, and `git diff --check`; they passed.
The resulting coverage was 100.0% for every owned Go package and total, and
100% across the frontend V8 metrics (99 statements, 72 branches, 20 functions,
95 lines). The final regression receipt records 54/54 browser tests, Docker
smoke, and a bounded fuzz run. The workflow enforces the coverage checks and
uses immutable action pins.

**Final verdict: approve.** The final diff is surgical for the stated
numeric-validation, offline-API-explorer, responsive-verification, and
meaningful-coverage scope. It may proceed to commit; hosted CI remains the
separate publication gate.

## Post-review configuration delta

[`frontend/eslint.config.js:8`](../frontend/eslint.config.js#L8) now ignores
only generated E2E reports (`test-results`, `playwright-report`) and generated
Swagger output (`public/docs`), alongside the existing build and coverage
directories. This prevents ESLint from traversing files that Playwright or the
asset-builder can replace concurrently; application, test-source, and authored
Swagger-template files remain linted. `pnpm lint` passed after the change.

**Delta verdict: approve.** This is a targeted generated-artifact exclusion and
does not weaken production-source linting or coverage.

## Publication-documentation delta

The publication documentation correctly records PR #2 at
`3cc9051b0c0ccbb6cb5da7df7bf0b889d274fcc6` and its successful hosted Actions
run 38681961976. The four affected records distinguish historical checkpoints
from current evidence, retain the hosted artifact totals, and direct merge and
post-merge status to the live PR. They do not claim that merge or a local
refresh has occurred. The PR is open and the hosted run reports success for the
frontend, backend, and browser/container jobs.

**Publication-documentation verdict: approve.** The delta is accurate,
internally consistent, and may proceed to its focused documentation commit.
