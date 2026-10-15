# Acceptance verification

The integration coordinator independently repeated the specialist checks before
publication. This receipt records observed results, not planned tests.

| Requirement                                | Evidence                                                                                                                                                                   | Result |
| ------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| Go backend and seven arithmetic operations | Go domain tables, HTTP handler tests, real-browser requests, container requests                                                                                            | Pass   |
| JSON validation and edge cases             | malformed/non-object/multiple/duplicate JSON; case sensitivity; numeric strings/null; missing operands; overflow; div0; negative sqrt; body limit; method and route errors | Pass   |
| React/TypeScript UI                        | ESLint with zero warnings, TypeScript check, production Vite build                                                                                                         | Pass   |
| Frontend unit tests                        | 8 tests; response rendering, request shape, validation, stale-result clearing, malformed payloads and full-response timeout                                                | Pass   |
| Go unit and integration tests              | `go test -race -coverprofile=coverage.out ./...`; `go vet ./...`                                                                                                           | Pass   |
| Browser integration                        | 36 Playwright tests: 12 each on Chromium desktop, Chromium at 320×700, and Firefox                                                                                         | Pass   |
| Responsive UI and keyboard use             | Mobile screenshot inspection; same browser journeys at 320px; Enter submits a calculation                                                                                  | Pass   |
| Network error handling                     | Explicitly aborted API request shows an error; next real request succeeds                                                                                                  | Pass   |
| Docker full stack                          | Multi-stage build; container UID 100; 8 successful calculations, 6 error regressions, HTML, JS asset and health requests                                                   | Pass   |
| Coverage                                   | Measured frontend and backend summaries in COVERAGE.md, generated reports and CI artifacts                                                                                 | Pass   |
| Setup/API/rationale/prompts                | README, SPEC.md, DECISIONS.md and PROMPTS.md, plus specialist receipts                                                                                                     | Pass   |

Normal E2E calculations use the real Go server and production frontend build.
Only the explicit network-failure test mocks transport. See TEST-PLAN.md for
individual scenarios and QA-RECEIPT.md for the browser configuration.

## Dependency and static analysis

`pnpm audit` reports zero vulnerabilities. Direct `govulncheck ./...` with
Go 1.27.2 exits successfully with “No vulnerabilities found”; its JSON stream
contains zero finding messages. Dependency-structure and source-size checks
pass. Production source files remain below 300 lines.

The optional, machine-installed engineering checker is not bundled with this
standalone take-home: `node scripts/checks/run-all.mjs --diff` is therefore
**UNAVAILABLE** in the repository. Its installed counterpart was also run.
The following outcomes are retained rather than described as a blanket pass:

- The hazard heuristic flags a comparison of a JSON delimiter as a secret
  comparison. This is a false positive: the operand is public request syntax,
  not a credential; the parser documents this at the call site.
- Its Go vulnerability adapter counts OSV advisory metadata as high findings,
  even though the authoritative scanner emits no actual finding messages.
  Direct scanner output adjudicates this as an adapter false positive.
- All 20 new direct frontend dependencies are accepted for this new repository:
  React/React DOM implement the required UI; the remaining packages are build,
  TypeScript, lint, unit/browser-test, DOM simulation or coverage tooling.
  Versions are pinned and the lockfile is frozen. The backend has no third-party
  dependency. No dependency is introduced for calculator arithmetic or styling.
- Aggregate CRAP discovery cannot locate nested coverage automatically. Explicit
  scoped runs with frontend LCOV and a line mapping of Go coverprofile blocks
  pass: six frontend and fourteen Go production functions, no score above 30.
  The aggregate discovery result itself remains **UNAVAILABLE**.
- Brand-token drift and harness fix-recurrence checks are **UNAVAILABLE** because
  their product/harness configuration is outside this application's scope.

These are explicit adjudications, not suppressed findings. Required application
checks remain reproducible through the Makefile and GitHub Actions workflow.

## Final review corrections

The final Go test revision additionally proves negative-base integer power,
negative exponents, decimal arithmetic, zero-numerator division, both percentage
scaling branches and missing/nonstring operations. Success assertions reject
NaN and extra/trailing JSON instead of accepting weak comparisons.
`actionlint .github/workflows/ci.yml` passes after fixing its unused loop variable.
The documented E2E port is fixed at 18080 and must be available.

## Hosted execution

[GitHub Actions run 38672809082](https://github.com/willianpinho/sezzle-calculator/actions/runs/38672809082)
passed the frontend, backend, and browser/container jobs for the published
implementation. The run includes all 36 browser cases, race-enabled Go tests,
frontend coverage, production builds, and a real container API/UI smoke test.
Its artifacts contain backend coverage, frontend coverage, and the browser
report. [PR #1](https://github.com/willianpinho/sezzle-calculator/pull/1) records
checks for subsequent documentation/workflow refinements and the final merge.

The generic local orchestration verifier does not discover this repository's
nested Go and frontend modules: its acceptance probe skips and its full probe
executes zero checks. Those empty results are not acceptance evidence. The
explicit module commands above and the hosted workflow provide the evidence.
