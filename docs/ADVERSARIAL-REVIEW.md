# Adversarial review

**Scope:** bounded pre-publication review of the calculator specification,
implementation, tests, delivery tooling, and evidence. This is structured
dissent, not an approval, recommendation, or vote.

## Steelman

The strongest version of the submission is a deliberately small stateless
monolith. Go owns strict request validation and arithmetic, while React is an
accessible same-origin client. The contract explicitly resolves ambiguous
calculator behavior: finite IEEE-754 binary64 values, `0^0 = 1`, percentage as
`a% of b`, and domain errors for non-finite or non-real results. One production
process serves the API and compiled frontend, so Playwright exercises the same
topology shipped in Docker.

The evidence is broader than a typical take-home: domain and HTTP tests, React
tests, three Playwright projects using the live Go API, production builds, a
non-root container smoke test, measured coverage with the uncovered server
entrypoint disclosed, and direct dependency scans. The design avoids databases,
queues, authentication, and service boundaries that the assignment does not
need.

## Assumption audit and attack

Likelihood is qualitative because the repository has no measured
failure-frequency dataset. `Reproduced` means the discrepancy was directly
observed; `low` names a plausible but unobserved residual risk.

| Assumption                                             | Disconfirming case and evidence                                                                                                                                                                                                                                                 | Likelihood                        | Status                                                     |
| ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- | ---------------------------------------------------------- |
| Passing tests prove finite arithmetic results          | The success-table assertion compared `math.Abs(got-want) > delta`. A `NaN` result would make that comparison false and could pass. The implementation rejected `NaN`, so this was a false-green test path rather than an observed product defect.                               | Reproduced test weakness          | Resolved: success assertions now reject `NaN` and infinity |
| The test plan describes checked-in coverage exactly    | The plan claimed valid negative-base power, negative exponents, decimal operands, and `0 / nonzero`, but those cases were absent. It also claimed every HTTP success response was verified as exactly `{result}`, while the test decoded into a struct that ignored extra keys. | Reproduced evidence mismatch      | Resolved: cases and whole-object assertion added           |
| The GitHub Actions workflow is lint-clean              | `actionlint` initially failed with ShellCheck `SC2034` because the retry-loop variable was unused. The workflow behavior was valid, but the warning contradicted the zero-warning quality gate.                                                                                 | Reproduced lint failure           | Resolved: fresh `actionlint` exits 0                       |
| Browser tests use an ephemeral port                    | Playwright intentionally hard-codes port `18080` and refuses to reuse a server. This prevents stale-server false greens at the cost of a possible local collision.                                                                                                              | Low portability risk              | Accepted: the test plan now requires that port to be free  |
| The frontend could silently diverge from the Go API    | Normal browser journeys wait for real `POST /api/calculate` responses from the production frontend served by Go. A fresh run passed all 36 cases across desktop Chromium, 320px Chromium, and Firefox.                                                                          | Low within tested browsers        | Disconfirming case not found                               |
| Coverage headlines imply complete product verification | The report separates unit coverage from browser/container execution, states 0% for `cmd/server`, and reports frontend branch coverage as well as line coverage. The overstatement was in the scenario list, not the coverage table.                                             | Low after reading the full report | Transparently bounded                                      |
| Public delivery is already complete                    | The GitHub repository is public and uses `main`, but the implementation remained on the feature branch during this review. Publication, hosted CI, merge, and post-merge verification belong to the next stage.                                                                 | Reproduced phase gap              | Explicit handoff item                                      |

## Findings that required closure

1. Make arithmetic success assertions fail on `NaN` or infinity and add the
   edge cases already claimed by the test plan.
2. Verify the complete HTTP success object when claiming it is exactly
   `{result}`.
3. Clear the reproducible `actionlint` / `SC2034` workflow warning.

All three were changed during the review. Fresh `actionlint`, Go race tests, Go
vet, and the strengthened tests passed. These findings concerned test and
delivery evidence. The review did not reproduce an incorrect result, validation
bypass, frontend/backend disconnect, failed production build, or container
runtime defect.

## Accepted tradeoffs and residual cracks

- Binary64 displays ordinary floating-point artifacts such as the result of
  `0.1 + 0.2`. That behavior is explicit and consistent across Go and
  JavaScript. An evaluator could still expect decimal semantics.
- `0^0 = 1` and `a% of b` are defensible but not universal calculator choices.
  Labels, examples, and the contract make the choices visible.
- The fixed E2E port can collide. This affects test startup rather than
  application correctness, and the failure is immediate rather than silent.
- Unit coverage does not instrument `cmd/server`, CSS, or the React bootstrap.
  Live browser and container runs exercise those paths, and the coverage report
  discloses the boundary.
- Rate limiting is absent. The service is an anonymous, constant-time,
  body-bounded take-home endpoint bound to loopback by Compose.

## Independent evidence from this pass

- Go race tests, vet, and coverage passed. After the adversarial test additions,
  the fresh aggregate was 83.0%, with 100% calculator and 92.1% HTTP API
  statement coverage.
- Frontend lint, type-check, eight unit tests, coverage thresholds, and the
  production build passed. Coverage was 100% statements/functions/lines and
  85.71% branches.
- All 36 live Playwright cases passed on three projects in 10.6 seconds.
- A fresh Docker build completed. A container returned healthy JSON, calculated
  25% of 80 as 20, and ran as the unprivileged `app` user.
- `git diff --check` passed. `actionlint` failed on the initial review state and
  passed after the targeted correction.

## Sanitized review prompt

> Steelman and then attack a completed Go and React TypeScript calculator
> take-home. Check seven-operation arithmetic and numerical edge policies,
> strict JSON parsing, client/server integration, coverage claims, test
> sensitivity, install and CI portability, Docker runtime behavior, responsive
> browser evidence, documentation prompts, and public-repository delivery.
> Reproduce concrete discrepancies. Keep the review within a 2-4 hour exercise,
> separate product bugs from evidence defects, avoid invented precision, and
> return the decision to the proposer without a vote.

**Structural verdict:** steelmanned before attack: yes; reproduced failure modes
were stale or weak test evidence and workflow-lint failure, both resolved; the
remaining fixed-port collision risk is low and explicit; no incorrect
calculator, API, browser, build, or container behavior was reproduced.

## Decision handoff

The three evidence defects were closed locally and the fixed-port tradeoff was
made explicit. The final delivery owner must verify the resolved diff, hosted
CI, merge state, and public default branch; this review does not make that
decision.
