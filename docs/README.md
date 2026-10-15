# Documentation index

## Current sources of truth

| Need                                            | Document                                                                        |
| ----------------------------------------------- | ------------------------------------------------------------------------------- |
| HTTP request, response, and error contract      | [SPEC.md](SPEC.md) and [openapi.json](../backend/internal/httpapi/openapi.json) |
| Interactive API explorer                        | [API-DOCS-DESIGN.md](API-DOCS-DESIGN.md); live at `http://localhost:8080/docs/` |
| Binary64 rules and mathematical boundary matrix | [NUMERIC-VALIDATION.md](NUMERIC-VALIDATION.md)                                  |
| System topology and request flow                | [ARCHITECTURE.md](ARCHITECTURE.md)                                              |
| Design decisions and hazard receipt             | [DECISIONS.md](DECISIONS.md)                                                    |
| Test scope                                      | [TEST-PLAN.md](TEST-PLAN.md)                                                    |
| Coverage commands and limits                    | [COVERAGE.md](COVERAGE.md)                                                      |
| Current verification status                     | [VERIFICATION.md](VERIFICATION.md)                                              |
| Delivery setup                                  | [DELIVERY.md](DELIVERY.md)                                                      |
| Work phases and open items                      | [TASKS.md](TASKS.md), [WORKLOG.md](WORKLOG.md), [HANDOFF.md](HANDOFF.md)        |

When the application is running, OpenAPI is available at
`http://localhost:8080/api/openapi.json` and `/api/docs` redirects to the
explorer at `http://localhost:8080/docs/`. The bundled UI uses local assets and
fetches that same-origin document; it has no external runtime CDN or validator
dependency.

## Evidence and history

The receipt files preserve evidence from the phase named in each file. They are
historical records, not substitutes for a later full-suite run:

- [BACKEND-RECEIPT.md](BACKEND-RECEIPT.md) and [FRONTEND-RECEIPT.md](FRONTEND-RECEIPT.md): initial delivery.
- [QA-RECEIPT.md](QA-RECEIPT.md): initial browser acceptance.
- [NUMERIC-BACKEND-RECEIPT.md](NUMERIC-BACKEND-RECEIPT.md) and [NUMERIC-UI-RECEIPT.md](NUMERIC-UI-RECEIPT.md): numerical and responsive follow-up evidence.
- [SWAGGER-UI-RECEIPT.md](SWAGGER-UI-RECEIPT.md): local Swagger asset packaging evidence.
- [FRONTEND-COVERAGE-RECEIPT.md](FRONTEND-COVERAGE-RECEIPT.md): current C2
  frontend coverage measurement.
- [CODE-REVIEW.md](CODE-REVIEW.md), [API-REVIEW.md](API-REVIEW.md), and [ADVERSARIAL-REVIEW.md](ADVERSARIAL-REVIEW.md): independent reviews.
- [FOLLOWUP-CODE-REVIEW.md](FOLLOWUP-CODE-REVIEW.md): follow-up review and final
  delta approval; earlier preliminary findings remain preserved in the receipt.
- [FOLLOWUP-API-REVIEW.md](FOLLOWUP-API-REVIEW.md) and
  [FOLLOWUP-ADVERSARIAL-REVIEW.md](FOLLOWUP-ADVERSARIAL-REVIEW.md): follow-up
  domain and structured-dissent reviews. The latter records the resolved
  generated-output tooling advisory, with no unresolved finding.

[PROMPTS.md](PROMPTS.md) records sanitized task prompts. Do not treat a count
or coverage figure in an older receipt as current after source changes; use the
commands in [VERIFICATION.md](VERIFICATION.md) instead.

The follow-up publication record is [PR #2](https://github.com/willianpinho/sezzle-calculator/pull/2).
Its [hosted CI run](https://github.com/willianpinho/sezzle-calculator/actions/runs/38681961976)
is green; [VERIFICATION.md](VERIFICATION.md) records that result, while PR #2
remains authoritative for merge and post-merge state.
