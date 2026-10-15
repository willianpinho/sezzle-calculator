# Sezzle Calculator Specification

## Purpose and scope

Deliver a small, self-contained calculator take-home: a React + TypeScript user
interface backed by a Go JSON API. It supports binary addition, subtraction,
multiplication, division, power, percentage, and unary square root. The goal is
correctness, clarity, accessibility, responsive use, and an easy local/Docker
run path; persistence, accounts, history, and external services are out of
scope.

## Acceptance matrix

| Area       | Acceptance criterion                                                                                                                                                                           |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| UI         | Keyboard-operable, labelled controls; clear operation/result/error state; works at narrow mobile and desktop widths.                                                                           |
| Operations | `+`, `-`, `*`, `/`, power, square root, and percentage calculate according to the contract below.                                                                                              |
| API        | Go serves `POST /api/calculate`, JSON only, with documented success and error responses.                                                                                                       |
| Validation | Invalid JSON, wrong request shapes, unknown operations, missing operands, divide-by-zero, negative square root, and non-finite results produce safe JSON errors.                               |
| Tests      | Go unit and HTTP integration tests; React component/unit tests; browser E2E covers a successful calculation and a visible validation error. Coverage commands produce reports for both layers. |
| Delivery   | `pnpm` is used for Node work. Docker Compose starts the complete stack. README explains setup, test/coverage commands, architecture, API, and tradeoffs.                                       |

## HTTP contract

`POST /api/calculate` accepts `Content-Type: application/json` and exactly one
JSON object:

```json
{ "operation": "add", "a": 2, "b": 3 }
```

`operation` is one of `add`, `subtract`, `multiply`, `divide`, `power`,
`sqrt`, or `percentage`. Field names and operation values are case-sensitive.
`a` and `b` are JSON numbers decoded as IEEE-754
binary64 values. They must be present for every binary operation. `sqrt` needs
only `a` and rejects a supplied `b`. `percentage` means **a percent of b**:
`a * b / 100`.

Success is `200 OK`:

```json
{ "result": 5 }
```

Calculation and API-route failures always use:

```json
{
  "error": { "code": "validation_error", "message": "b is required for divide" }
}
```

| Status | `error.code`             | Conditions                                                                                                  |
| ------ | ------------------------ | ----------------------------------------------------------------------------------------------------------- |
| 400    | `invalid_json`           | Malformed JSON, duplicate keys, multiple JSON values, or a non-object body.                                 |
| 413    | `payload_too_large`      | Body exceeds the server's documented request-body limit.                                                    |
| 415    | `unsupported_media_type` | Request content type is not JSON.                                                                           |
| 422    | `validation_error`       | Unknown/extra fields; unknown operation; missing, null, or non-number operands; operand supplied to `sqrt`. |
| 422    | `domain_error`           | Division by zero; negative square root; a calculation whose result is NaN or infinity.                      |
| 405    | `method_not_allowed`     | Any method other than POST on this route.                                                                   |
| 404    | `not_found`              | Unknown API route (`/api` or `/api/*`).                                                                                          |

The API rejects unknown request fields to catch client mistakes. It does not
accept `NaN` or infinity as inputs; JSON itself has no valid representation for
them. Results follow Go/IEEE-754 binary64 arithmetic, then must be finite.
Overflow and indeterminate operations therefore return `422 domain_error`.
Underflow to a finite zero is valid. `0^0` explicitly returns `1`, matching
Go's `math.Pow(0, 0)`. A negative base with a non-integral exponent has no real
result and returns `422 domain_error`.

## Design and delivery boundaries

One stateless Go process exposes the API and serves the built frontend from
`STATIC_DIR` in production, giving browser requests a same-origin `/api` path.
Vite proxies `/api` in development, so no browser CORS policy is needed. No database, cache, auth,
background job, queue, migration, or separate service is justified by the
exercise's measured need. The Go server owns validation and math; the client
may provide immediate form guidance but treats the API as authoritative.

The service should set a bounded request-body size, send `application/json` on
all API responses, and never expose Go/internal errors. `GET /api/health`
supports Docker health checks and returns a JSON status object; it is outside
the calculation contract.

Static files outside the `/api` namespace use standard HTTP file-serving
semantics, including plain-text 404 responses for missing assets.
