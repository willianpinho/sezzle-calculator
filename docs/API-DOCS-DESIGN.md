# API documentation design

## Decision

`backend/internal/httpapi/openapi.json` is the authoritative, hand-reviewed
OpenAPI 3.1 contract for the calculator API. The Go API serves that exact
embedded file from `GET /api/openapi.json`. `GET /api/docs` redirects to
`/docs/`, where the built frontend serves a bundled Swagger UI using local
`swagger-ui-dist` files. The browser UI loads the same-origin specification;
no CDN, external runtime network dependency, or generated server contract is used.

The contract declares the same-origin server as `/` and root-level
`security: []`: this anonymous calculator intentionally has no authentication
scheme.

This remains a compact Level-2 REST surface: one calculation resource action,
one health read, explicit HTTP status codes, and a stable JSON error envelope.
The static-file boundary stays outside the JSON API contract.

## Contract representation

The document models calculation requests as a `oneOf`:

- binary operations (`add`, `subtract`, `multiply`, `divide`, `power`, and
  `percentage`) require `operation`, `a`, and `b`;
- `sqrt` requires `operation: "sqrt"` and `a`, and each branch has
  `additionalProperties: false`.

OpenAPI/JSON Schema cannot reliably require a server to reject duplicate JSON
object keys, so duplicate-key rejection remains stated in the response
description and enforced by the Go parser. `FiniteNumber` is descriptive
because JSON Schema has no portable finite-number keyword: the handler remains
the authority for ParseFloat range checking, signed zero, and calculation
result validity.

## Error and version policy

Every documented API failure returns `{ "error": { "code", "message" } }`.
Consumers should branch on HTTP status and `error.code`; `message` is safe for
display but is not a stable programmatic value. There is no URL version segment
for this initial internal/take-home surface. A breaking request or response
change requires `/api/v2`, a migration note, and a published deprecation window
before retiring `/api`.

## Implemented delivery checks

The backend now embeds and serves the document and the frontend bundles the
viewer. Verify the current implementation with:

```bash
node -e 'JSON.parse(require("node:fs").readFileSync("backend/internal/httpapi/openapi.json", "utf8"))'
(cd backend && go test ./...)
(cd frontend && pnpm lint && pnpm typecheck && pnpm build)
make e2e
make docker-build
```

The implemented API tests compare the served OpenAPI document and verify the
docs redirect in [docs_test.go](../backend/internal/httpapi/docs_test.go).
Browser acceptance loads `/docs/` from the production-shaped stack and confirms
the calculator paths remain available in
[swagger-ui.spec.ts](../frontend/e2e/swagger-ui.spec.ts).

Redocly's minimal ruleset reports one intentional warning: `/api/docs` has a
`307` redirect instead of a success response. The redirect is the endpoint's
documented behavior, so the warning is adjudicated rather than hidden or
changed into a fictitious `2xx` response.
