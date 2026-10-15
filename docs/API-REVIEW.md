# API and domain review

**Scope:** independent review of the calculator contract, Go HTTP/domain
implementation, client API use, tests, README, and boundary behavior on
2026-10-14. This review made no implementation changes.

## Verdict

**APPROVE WITH COMMENT.** The seven operations and their edge policies agree
with the contract. The handler consistently returns the documented JSON error
envelope for calculation and `/api`-namespace routes, and the client correctly
uses that contract. No implementation change is required before the next review
gate.

## Finding

### Low — clarify the scope of the documented 404 rule

`docs/SPEC.md` says that `not_found` applies to “Any unknown route.” The
production server intentionally delegates non-`/api` paths to Go's static file
server. As a result, an unknown static path such as `GET /not-a-route` returns
the file server's normal `404 text/plain` body, while `GET /api/nope` returns
the documented JSON `{ "error": { "code": "not_found", ... } }` response.

This is appropriate for a single-page/static delivery boundary and does not
affect calculation consumers. Update the spec wording to say “any unknown API
route” (or explicitly define a JSON response for unknown static paths) so the
contract matches the deployed behavior exactly.

## Verified behavior

- All seven operations return the expected finite binary64 result. Percentage
  is implemented as `a% of b`, with reordered scaling to avoid avoidable
  intermediate overflow; `0^0` returns `1`.
- Division by zero, negative square root, non-real negative-base fractional
  power, and non-finite outcomes return `422 domain_error`. Finite underflow
  to zero remains valid.
- The parser accepts exactly one object, catches duplicate keys and multiple
  values as `400 invalid_json`, rejects unknown fields and invalid operand
  shapes as `422 validation_error`, enforces a 1 MiB body cap, and rejects
  non-finite numeric inputs.
- `POST /api/calculate` is the only calculation method; its other methods
  return JSON `405` with `Allow: POST`. Unknown `/api` routes return JSON 404.
  `GET /api/health` is correctly documented as outside the calculation
  contract.
- The React client sends the specified request shape, omits `b` for square
  root, validates finite browser input, displays server errors, and accepts
  only finite numeric success results.
- Direct checks passed: `go test -race ./...`, `go vet ./...`, frontend unit
  tests (8), lint, and type-check. Focused live requests confirmed 0^0,
  percentage at `MaxFloat64`, finite underflow, duplicate keys, multiple JSON
  values, numeric overflow input, and the routing distinction above.

## Accepted assumptions

- The error-message text is explanatory rather than a compatibility guarantee;
  the status and `error.code` are the stable contract.
- The 404/error envelope applies to API endpoints. Static assets follow normal
  static-server semantics once the specification is clarified.
- The percentage implementation may choose a numerically safer equivalent
  binary64 evaluation order; its required semantic result is `a * b / 100`
  when finite.

## Sanitized review prompt

> Review a small stateless Go JSON calculator API and React client. Verify the
> seven operations, strict one-object JSON parsing (including duplicate and
> unknown fields), finite binary64/domain-error policy, method and API-route
> errors, and client request/response handling. Compare specification, tests,
> README, and deployed static routing. Report only concrete contract or domain
> discrepancies; keep the review bounded to the stated calculator scope.

## Finding resolution

The specification now explicitly scopes JSON `not_found` errors to the `/api`
namespace and documents ordinary static-file 404 responses. The documentation
comment is resolved; no runtime behavior changed.
