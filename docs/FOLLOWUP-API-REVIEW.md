# Follow-up API and numeric-domain review

**Scope:** final bounded review of `test/numeric-boundaries` against
`origin/main` on 2026-10-15. Reviewed the seven-operation finite-binary64
policy, strict JSON parser, embedded OpenAPI contract, Swagger routing, and
their specification/README documentation. No code was changed by this review.

## Verdict

**APPROVE.** No API or numeric-domain blocker was found.

## Contract findings

- `calculator.Calculate` now rejects `NaN` and both infinities in either
  argument before evaluation. This closes the direct-call gap where documented
  `math.Pow` special cases could otherwise turn a non-finite input into a
  finite result.
- The seven operations preserve the agreed finite-binary64 policy: `0^0 = 1`;
  zero divisors of either sign, zero to a negative power, negative square root,
  non-integral negative-base power, and non-finite results produce the stable
  domain error. Finite underflow, subnormals, and signed zero remain valid.
- The strict decoder consumes exactly one object, detects duplicate keys,
  rejects unknown fields, validates JSON-number syntax/range before arithmetic,
  and returns one JSON error envelope. The refactored object-member assertion
  is safe because `Decoder.Token` yields a string for an object member name.
- `sqrt` requires `a` and rejects supplied `b`; binary operations require both
  operands. These rules agree across the handler, tests, OpenAPI `oneOf`,
  SPEC, and README.
- The OpenAPI 3.1 document matches the implemented calculation, health,
  OpenAPI, and documentation routes, including method errors and the
  same-origin public-server policy. It explicitly models strict request object
  properties; duplicate-key rejection remains a parser rule because JSON
  Schema cannot portably express it.

## Verification performed

- `go test -race ./...` and `go vet ./...` passed in `backend/`.
- OpenAPI request/routing invariants were parsed directly from
  `backend/internal/httpapi/openapi.json`.
- Redocly minimal lint reported no errors. Its sole warning is intentionally
  retained: `GET /api/docs` correctly has only a `307` redirect, so a made-up
  `2xx` response would make the machine contract less accurate.
- `git diff --check origin/main` passed.

The boundary matrix is intentionally finite, fixture-based evidence rather
than a claim to exhaust all real numbers or float64 values.
