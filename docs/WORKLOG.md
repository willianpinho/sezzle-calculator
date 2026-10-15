# Work log

## 2026-10-14 — specification baseline

The assignment requires a Go backend and React/TypeScript frontend for a small
calculator. The chosen scope is a stateless standard-library Go JSON API and a
Vite React client. The API contract defines every requested operation, response
shape, error category, finite-number policy, `0^0`, and the percentage meaning.

The repository will use `pnpm` for all Node commands. Use Go 1.26 or later with
the current patch toolchain available to CI/container; keep the exact image
version in delivery configuration. No database means no migrations, cache,
queues, sessions, or job topology. This is deliberate scope control, not a
deferred implementation.

Production delivery has the Go process serve the built frontend via `STATIC_DIR`
and its own same-origin API routes; Vite uses a development proxy. Implementation
must preserve server authority for validation, add tests at both
layers plus browser E2E, provide coverage reports, and use Docker Compose for
the complete stack. API contract source: `docs/SPEC.md`.

## Integration checkpoint

The production topology was clarified before implementation: Go serves the
compiled frontend directly from `STATIC_DIR`; `/api/health` is the health
route. The development proxy is Vite-only. Dependency versions were verified
against package registries, and the Go 1.27.2 toolchain and Docker image tags
were independently checked.

Early validation found a quoted-number decoding regression in Go: decoding a
JSON string into `json.Number` can accept a numeric string. A failing integration
test confirmed it; the backend specialist owns the correction. Review also
identified non-object JSON double-response handling, oversized leading
whitespace, and percentage intermediate overflow as boundary cases to cover.
Frontend review added malformed-response validation and a request timeout that
must cover both response headers and body. These are regression fixes within
the original validation/error-handling criteria.

Initial frontend build, lint, type-check, and coverage execution passed. Initial
pnpm audit and Go vulnerability checks reported no vulnerabilities. Final
numbers and final-state verification belong in the acceptance receipt; these
intermediate checks are not a completion claim.

## Acceptance and review checkpoint

The coordinator independently repeated eight frontend unit tests and all 36
Playwright cases across desktop Chromium, 320px Chromium and Firefox. Go race
tests and vet passed. Coverage is published in COVERAGE.md, including the
server entrypoint's zero unit coverage rather than hiding it. Visual desktop
and mobile inspection found no clipping or overlapping controls.

The final Docker build ran as UID 100 and passed eight calculation cases, six
error regressions, HTML/JavaScript delivery and health requests. A separate
Compose startup reached healthy state and served a real calculation before
cleanup. No deployed service or external infrastructure is needed to evaluate
the assignment.

Independent code review approved with comments. Its stale QA-receipt concern
was corrected and the public command/result summary consolidated in
VERIFICATION.md and DELIVERY.md. The local engineering checker's JSON-token
and Go-advisory-metadata findings were investigated and adjudicated as false
positives; optional unavailable checks remain explicitly marked unavailable.
The remaining delivery stages are API-domain review, adversarial review,
publication and CI verification.

API-domain review approved the implementation. Its only comment clarified
that JSON 404 errors belong to the API namespace; static-file errors retain
ordinary HTTP file-server semantics. SPEC.md now makes that boundary explicit.

Adversarial review exposed an unused workflow loop variable and missing cases
claimed by the test plan. The workflow now passes actionlint. Added explicit
negative-base power, negative exponent, decimal, zero-numerator, missing/invalid
operation and percentage-scaling tests. Math assertions now reject nonfinite
results; HTTP success assertions require exactly one complete JSON result
object. Independent Go race tests/vet pass, with calculator coverage 100%, HTTP
coverage 92.1% and total backend coverage 83.0%. The E2E port is deliberately
fixed at 18080 and must be free; the README and test plan now state that
constraint. Production source behavior did not change in these final fixes.

All three review stages are complete. The code reviewer also approved the final
test/workflow delta after the adversarial corrections. The final local gate run
reproduced only the already-adjudicated checker false positives and explicitly
unavailable optional checks. Publication and hosted CI are the remaining steps.

## GitHub publication checkpoint

The implementation is published in the public repository through
[PR #1](https://github.com/willianpinho/sezzle-calculator/pull/1).
[The initial hosted CI run](https://github.com/willianpinho/sezzle-calculator/actions/runs/38672809082)
passed all three jobs: frontend, backend, and browser/container validation.
Coverage reports and the Playwright report were uploaded as run artifacts.

The final documentation read found two shell-directory issues: installing
Playwright browsers must not leave subsequent Makefile commands in the frontend
directory, and Go's HTML coverage viewer must run inside the backend module.
Both commands were corrected and the coverage HTML was generated successfully.
The runner's action-runtime deprecation annotations prompted an update to
official supported action releases and an explicit Ubuntu 24.04 runner. The
final documentation/workflow revision is reviewed separately and must pass the
same hosted jobs before merge; the PR records the final head and merge status.
