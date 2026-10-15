# Historical acceptance verification

The integration coordinator independently repeated the specialist checks before
the initial publication. This receipt records observed results, not planned
tests. It predates the numeric-boundary, Swagger, and expanded responsive
follow-up; use [TASKS.md](TASKS.md) for the active phase and re-run the listed
commands before treating it as current evidence.

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
- The initial delivery's 20 direct frontend dependencies are accepted for this
  new repository: React/React DOM implement the required UI; the remaining packages are build,
  TypeScript, lint, unit/browser-test, DOM simulation or coverage tooling.
  Versions are pinned and the lockfile is frozen. The follow-up adds the 21st,
  `swagger-ui-dist`, solely to bundle local API-explorer assets; it is not
  calculator runtime logic. A full `pnpm audit` across all dependencies reported
  zero known vulnerabilities. The backend has no third-party dependency.
- Aggregate CRAP discovery cannot locate nested coverage automatically. The
  historical scoped run covered six frontend and fourteen Go production
  functions. The current scoped rerun covers six frontend files and twenty Go
  production functions; all pass the threshold of 30.
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

## Follow-up evidence in progress

The numerical, OpenAPI, Swagger, and responsive follow-up has not yet produced
a replacement consolidated receipt. The following fresh local checks are
recorded without treating them as a hosted-CI result:

OpenAPI schema validation passes. Redocly's remaining minimal-ruleset warning
is adjudicated: `/api/docs` intentionally returns only its documented `307`
redirect and must not advertise a fictitious `2xx` response.

- A fresh Docker image served the UI and `/api/health` on port 18081. Its served
  OpenAPI document parsed identically to the checked-in source; `/api/docs`
  returned `307` to `/docs/`; the local Swagger JavaScript, CSS, bundle-license
  reference, `LICENSE`, and `NOTICE` returned nonempty `200` responses.
- Container requests confirmed `0^-1` and `(-8)^(1/3)` return `422
domain_error`, while `0^0` returns `1` and `100% of 1e308` succeeds.
- A 15-second, two-worker calculator fuzz run completed 1,463,339 executions
  without a failure. It is bounded invariant evidence, not exhaustive proof.

The final browser run passed **54/54** after the local Swagger assets were
installed and bundled. It includes the 2 numeric boundary tests; 7 responsive
projects at 320, 375, 390, 768, 1024, and landscape Chromium viewports plus
WebKit iPhone emulation; and baseline calculator/API/Swagger journeys in
Chromium, Firefox, and WebKit. Swagger acceptance blocks external requests,
loads the real same-origin contract, follows `/api/docs`, and sends live add
and divide-by-zero requests through Try it out. Responsive screenshots are
retained as ignored local artifacts by [NUMERIC-UI-RECEIPT.md](NUMERIC-UI-RECEIPT.md).
Hosted CI remains pending D5 publication.

## Follow-up dependency and coverage checks

Fresh follow-up audits reported zero known vulnerabilities from `pnpm audit`
across all frontend dependencies and from
`go run golang.org/x/vuln/cmd/govulncheck@latest -json ./...` (zero finding
records). The repository-local mandated checker remains **UNAVAILABLE** because
its script is absent. Its installed counterpart still reports OSV metadata and
a JSON-delimiter comparison false positive; the added `swagger-ui-dist`
dependency is deliberate local API-documentation tooling. These are explicit
adjudications; C3's final scoped CRAP and regression run are recorded below.

Backend C1 now enforces 100.0% total Go statement coverage through the standard
race profile: `cmd/server`, `internal/calculator`, and `internal/httpapi` each
report 100.0%. The frontend coverage command enforces 100% global V8 thresholds
over its six handwritten runtime modules, including bootstrap and the Swagger
asset builder. Its current 22-test local run measures 99/99 statements, 72/72
branches, 20/20 functions, and 95/95 lines, as documented in
[FRONTEND-COVERAGE-RECEIPT.md](FRONTEND-COVERAGE-RECEIPT.md). The independent
C3 review repeated the unit, regression, browser, and Docker checks and
approved the final code/test/CI delta; see
[FOLLOWUP-CODE-REVIEW.md](FOLLOWUP-CODE-REVIEW.md). Hosted CI publication
remains separate and pending D5.

Follow-up API review approved the binary64/OpenAPI contract. The adversarial
review found no high- or medium-severity issue. Its one reproduced low local
race—lint encountering Playwright output during an intentionally parallel
run—is resolved: ESLint now ignores only generated output directories, and an
independent concurrent lint/E2E rerun passed. This does not change what lint or
coverage measures; see [FOLLOWUP-ADVERSARIAL-REVIEW.md](FOLLOWUP-ADVERSARIAL-REVIEW.md).
