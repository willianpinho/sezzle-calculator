# Independent code review

**Scope:** `feat/calculator` working tree against `a2ace4b`, reviewed on
2026-10-14. The stated intent is a compact Sezzle calculator with a React +
TypeScript UI, a Go standard-library JSON API, all seven required operations,
strict validation, real integration/E2E coverage, Docker delivery, and clear
documentation.

## Findings

### Resolved concern — verification receipt was stale

The prior text at [`docs/QA-RECEIPT.md:35`](QA-RECEIPT.md#L35) said the final
delivery gate still needed to run despite completed evidence. It now records the
completed coordinator checks and links to the public consolidated receipt.
**No further action is required.**

### Praise — contract boundaries are enforced rather than only documented

[`backend/internal/httpapi/api.go:94`](../backend/internal/httpapi/api.go#L94)
parses exactly one object, detects duplicate keys, rejects extra fields, bounds
the body, and produces the documented JSON error envelope. The accompanying
handler tests cover malformed JSON, shape errors, domain errors, routing, and
the byte limit. The browser tests exercise the production-built frontend against
the live Go endpoint for every required operation.

## Correctness and security assessment

The calculator implementation covers add, subtract, multiply, divide, power,
square root, and percentage. The API checks media type, request shape,
finite-number input/result rules, domain errors, method handling, API routing,
and a 1 MiB body cap. The UI has labelled native controls, keyboard form
submission, loading/error/result states, stale-response cancellation, and
mobile browser journeys at 320px with manual screenshot inspection.

No credentials, authentication state, persistence, migrations, queues, or
cross-service writes exist in scope. Retrying the pure calculate request is
idempotent; race, outbox, N+1, locking-migration, pool, cache-stampede, DLQ,
and streaming hazards are not applicable. The documented proxy-level rate-limit
decision is proportionate for this anonymous, bounded take-home service.

The production-hazard scan's high finding at
[`backend/internal/httpapi/api.go:98`](../backend/internal/httpapi/api.go#L98)
is a false positive: `token != json.Delim('{')` compares a parsed public JSON
grammar token with a delimiter. It does not compare secret material, so
constant-time comparison is neither relevant nor useful.

## Verification evidence

- Directly run: `go test -race -coverprofile=coverage.out ./...`, `go vet ./...`,
  focused HTTP contract tests, frontend lint, type-check, coverage tests, and
  `make build`; all passed. Go coverage: calculator 90.5%, HTTP API 88.9%.
  Frontend coverage: 100% statements/functions/lines and 85.71% branches.
- Review evidence records 36 passing Playwright tests across desktop Chromium,
  320px mobile Chromium, and Firefox, plus a real Docker/UI/API smoke test.
- `git diff --check HEAD` passed. `pnpm audit --prod --audit-level=high` found
  no known vulnerabilities.
- The targeted CRAP checks passed: frontend 6 functions and Go 14 functions,
  with zero breaches. The aggregate auto-discovery CRAP gate remains
  **UNAVAILABLE** because it cannot locate the two coverage formats together;
  that is a harness limitation, not a coverage pass.
- Dependency structure and source-size gates passed. The 20 added npm packages
  are the declared React/Vite/test/lint toolchain required by this new project.
- The dependency gate's reported 180 Go highs are not vulnerabilities. Its
  parser counted `govulncheck`'s OSV metadata records; the tool emitted zero
  `finding` records and exited successfully against Go 1.27.2.
- Design-token and recurrence checks are **UNAVAILABLE**, respectively because
  this take-home has no configured brand-token source and no harness component
  in scope. They are not represented as passing checks.

## Review prompt for the next pass

Confirm that a future change preserves the listed validation outcomes and does
not weaken the JSON contract or live E2E path.

**Verdict: approve.** The documentation concern was resolved; no implementation,
security, or API-contract blocker was found.

## Final narrow delta review

Reviewed after the initial approval: only backend test coverage, workflow lint
cleanup, and coverage documentation changed; production source is unchanged.

- [`backend/internal/calculator/calculator_test.go:17`](../backend/internal/calculator/calculator_test.go#L17)
  now covers zero numerator division, signed integer power, negative exponent,
  decimal arithmetic, both percentage scaling branches, and explicitly rejects
  non-finite success values.
- [`backend/internal/httpapi/api_test.go:31`](../backend/internal/httpapi/api_test.go#L31)
  now verifies that each successful HTTP body is exactly the single `result`
  object; error coverage adds missing and non-string operations.
- [`.github/workflows/ci.yml:85`](../.github/workflows/ci.yml#L85) replaces the
  unused shell loop variable with `_`; the smoke assertion remains exact.
- [`docs/COVERAGE.md:11`](COVERAGE.md#L11) correctly reflects the rerun: 100.0%
  calculator, 92.1% HTTP API, and 83.0% aggregate backend statement coverage.

Direct rerun of `go test -race -coverprofile=coverage.out ./...` and `go vet
./...` passed; `git diff --check HEAD` passed. **Final delta verdict: approve.**

## Publication delta review

Reviewed the final publication-only delta: README shell-directory corrections,
hosted-CI evidence, and the workflow runtime update. No production code changed.

- [`README.md:42`](../README.md#L42) keeps browser installation in a subshell,
  so the following Make targets still execute from the repository root.
- [`README.md:48`](../README.md#L48) runs the Go HTML coverage viewer from its
  module directory, where the generated coverage profile exists.
- [`.github/workflows/ci.yml:18`](../.github/workflows/ci.yml#L18) pins actions
  to immutable release commits and Ubuntu 24.04. Direct `git ls-remote`
  verification matches the stated v7/v6 release refs, including the dereferenced
  annotated `pnpm/action-setup` v6.1.0 tag. `actionlint` and `git diff --check`
  pass.
- [`docs/VERIFICATION.md:68`](VERIFICATION.md#L68) accurately distinguishes the
  successful initial hosted run from the final workflow/doc revision, which must
  run in CI before merge. GitHub confirms run 38672809082 succeeded for all
  three jobs and published three artifacts.

**Publication delta verdict: approve.** This approval covers the uncommitted
documentation/workflow delta; the hosted CI run on its final head remains the
required merge gate.
