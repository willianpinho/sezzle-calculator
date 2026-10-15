# Follow-up adversarial review

**Scope:** final bounded review of `test/numeric-boundaries` against
`origin/main` on 2026-10-15. This receipt covers the seven-operation finite
binary64 contract, meaningful coverage claims, embedded OpenAPI and offline
Swagger delivery, responsive browser evidence, server lifecycle behavior, and
CI portability. It is structured dissent, not an approval, recommendation, or
vote.

## Steelman

The strongest version of this change treats the calculator as a small,
explicitly bounded system. Go remains the arithmetic authority; React only
parses finite browser inputs and renders API results. The specification defines
signed zero, subnormals, overflow, underflow, `0^0`, negative-base powers, and
percentage semantics without claiming decimal, financial, or arbitrary-
precision behavior.

The added evidence is behavior-oriented rather than a raw coverage exercise.
Go tests drive numerical boundaries, strict HTTP parsing, documentation routes,
listener failures, real shutdown, timeout, and fatal startup behavior. Frontend
tests cover all six handwritten runtime/tooling modules, while 54 browser tests
exercise the production build and live Go API across Chromium, Firefox, WebKit
emulation, numerical boundaries, Swagger Try it out, and responsive widths.
The container path bundles Swagger locally and serves the same embedded OpenAPI
document used by the explorer.

## Adversarial findings

### Resolved low: parallel local lint and E2E raced over generated test output

A deliberately parallel verification run started `pnpm lint` while Playwright
was deleting and recreating `frontend/test-results`. ESLint traversed that
directory between those operations and exited with `ENOENT`. The supported
local sequence and GitHub Actions do not overlap these commands: the browser
job depends on the completed frontend job. A sequential rerun of lint,
type-check, coverage, and build passed. This is a low-likelihood developer-
tooling concurrency edge, not a product or hosted-CI failure.

The frontend ESLint configuration now ignores only generated directories:
`dist`, `coverage`, `test-results`, `playwright-report`, and `public/docs`.
It does not change the Vitest coverage denominator. A direct repeat launched
lint and E2E concurrently; lint passed with zero warnings while all 54 browser
tests passed (16.9 seconds in the implementation receipt and 19.2 seconds in
this independent repeat). The reproduced race is resolved.

No unresolved finding remains from this review.

## Strongest disconfirming cases checked

| Assumption under attack                                     | Disconfirming case                                                                                               | Evidence and disposition                                                                                                                                                                                                                                                  |
| ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Finite input policy is enforced only at HTTP                | A direct caller supplies `NaN` or infinity, including the unused `b` argument for square root                    | `Calculate` rejects non-finite values in either argument; table and fuzz seeds cover the boundary. Resolved risk.                                                                                                                                                         |
| The boundary matrix overclaims mathematical proof           | Power and transcendental behavior are judged by a circular oracle, or a few fixtures are presented as exhaustive | The independent `big.Float` oracle is limited to basic arithmetic. Power follows the documented Go binary64 contract, and the matrix and receipts explicitly disclaim exhaustive proof. Accepted bounded evidence.                                                        |
| 100% coverage is assertion theater                          | Tests execute lines without checking externally visible behavior                                                 | Server tests assert process exits and HTTP lifecycle outcomes; API tests assert status, envelope, and routes; frontend tests assert DOM and requests; browser tests use the live server. No tautological blocker was found.                                               |
| OpenAPI can diverge from runtime unnoticed                  | Swagger displays a schema that the handler does not implement                                                    | The embedded document is served verbatim, references resolve, required paths are checked, unary and binary shapes match the parser, and browser Try it out verifies both success and domain-error requests. Future drift remains possible but was not found in this diff. |
| “Offline Swagger” still contacts a CDN or validator         | The explorer renders only when external traffic succeeds                                                         | Browser tests abort every non-application-origin request and still load the UI, contract, and Try it out flow. Docker checks the bundled assets. Resolved risk.                                                                                                           |
| Responsive evidence is stronger than the environment tested | Emulation is described as proof on physical phones                                                               | The docs call WebKit coverage iPhone emulation and do not claim physical-device proof. Chromium viewport tests also check horizontal overflow and long results. Accepted limitation.                                                                                      |
| Hosted CI is already proven                                 | Local success is presented as a completed GitHub run                                                             | Verification and task records keep hosted CI and publication pending as the next gate. Explicit handoff item.                                                                                                                                                             |

## Independent evidence from this pass

- `go test -race -count=1 -coverprofile=coverage.out ./...` and `go vet ./...`
  passed; every Go package and the aggregate reported 100.0% statements.
- Sequential frontend lint, type-check, 22 tests, coverage, and production build
  passed; all six included files reported 100% statements, branches, functions,
  and lines.
- All 54 Playwright tests passed in 18.2 seconds against the production build
  and live Go server.
- After the targeted ESLint ignore change, a second deliberately concurrent
  lint and E2E run passed with zero lint warnings and all 54 browser tests.
- `actionlint`, JSON parsing of the OpenAPI document, and
  `git diff --check origin/main` passed.
- The consolidated verification receipt records a passing non-root Docker
  image smoke test and zero actual findings from the frontend and Go dependency
  audits.

## Residual accepted risks

- Binary64 rounding and the chosen `0^0` and percentage conventions may differ
  from another calculator's conventions; the UI and contract disclose them.
- `math.Pow` and `math.Sqrt` rely on Go's runtime behavior. Fixtures are bounded
  to documented, reproducible outcomes rather than cross-architecture claims
  about every input.
- The E2E server reserves port `18080`; a collision fails fast, and the README
  states the requirement.
- The OpenAPI file remains a manually maintained contract. Current parity has
  multiple checks, but future changes still require contract review.

**Structural verdict:** steelmanned before attack: yes; the only reproduced
failure mode was a low-likelihood parallel local-tooling race, now resolved;
numeric correctness, coverage sensitivity, API-document parity, offline assets,
responsive behavior, and CI sequencing produced no unresolved disconfirming
evidence.

## Decision handoff

The adversarial finding is resolved locally. Publication still requires the
planned hosted CI, merge-state, and post-merge verification gates; this review
does not make that decision.
