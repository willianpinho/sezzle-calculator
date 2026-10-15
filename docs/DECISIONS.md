# Decisions

## ADR-001: Stateless modular monolith

**Decision:** Ship one Go API and one React/Vite client in a single repository,
with Docker Compose for local full-stack execution.

**Trigger:** This is a 2–4 hour calculator exercise with one synchronous
calculation endpoint and no durable state. No latency, load, ownership, or
deployment evidence warrants a service split.

**Rationale:** A standard-library Go HTTP server keeps the required backend
small and observable. React/TypeScript provides the requested browser UI.
Composition keeps each layer independently testable without operational cost.

## ADR-002: Contract-first synchronous JSON API

**Decision:** Use `POST /api/calculate` with operation enum and numeric
operands, returning `{ "result": number }` or `{ "error": { "code", "message" } }`.

**Trigger:** The UI needs one request/response calculation action. There is no
measured benefit for GraphQL, RPC tooling, persistence, async work, or command
and query separation.

**Rationale:** This shape makes validation, API tests, and error presentation
straightforward. The Go server is the single arithmetic authority.

## ADR-003: Binary64, finite real-result policy

**Decision:** Decode JSON numbers to `float64`; reject non-finite outcomes;
return `1` for `0^0`; define percentage as `a% of b`.

**Trigger:** The assignment asks for common calculator operations but gives no
precision model or edge-case definitions.

**Rationale:** Binary64 is native to Go and JavaScript and needs no dependency.
Finite-result validation gives callers deterministic errors for overflow and
undefined real arithmetic. Decimal financial arithmetic would add complexity
without a stated currency or precision requirement.

## ADR-004: Single-process static delivery and no authentication

**Decision:** The Go process serves the built frontend from `STATIC_DIR` and
its own `/api` routes in production; Vite proxies `/api` in development. The
endpoint has no auth or session.

**Trigger:** The anonymous calculator has no user data or protected action.

**Rationale:** Same origin avoids CORS configuration and removes a production
proxy process. Authentication and
sessions would have no security subject to protect in this scope.

## Production-hazard receipt

| Hazard                                        | Applicability and decision                                                                                                                                                                       |
| --------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Idempotency / races                           | Applicable to HTTP retries, but calculation is pure and has no side effect or shared mutable state; repeat requests return the same result.                                                      |
| Dual write / outbox                           | Not applicable: no database, queue, cache, or third-party write.                                                                                                                                 |
| Poison messages / DLQ                         | Not applicable: no asynchronous consumers or jobs.                                                                                                                                               |
| N+1 / indexes / connection pools / migrations | Not applicable: no database.                                                                                                                                                                     |
| Cache stampede                                | Not applicable: no cache; calculations are cheap and computed inline.                                                                                                                            |
| Rate limit                                    | Public endpoint is bounded by body size and constant-time local arithmetic. Application rate limiting is out of scope for this demo; deployer/gateway can add it when exposure/load requires it. |
| Load testing                                  | No measured traffic target exists for the exercise. Unit, integration, and E2E tests are the proportional verification.                                                                          |
| Middleware order                              | If middleware is used: request ID/logging first, request-body limit/parser before handler, routing, then centralized JSON error handling. No auth exists.                                        |
| Streaming cleanup                             | Not applicable: no streaming endpoint.                                                                                                                                                           |
| Timing attacks                                | Not applicable: no secret comparison, credentials, or signatures.                                                                                                                                |
| CQRS                                          | Rejected: one synchronous calculation command and response need no divergent models.                                                                                                             |
