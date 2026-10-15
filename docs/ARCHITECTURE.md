# Architecture

The production application is one stateless Go process. It serves compiled
frontend files from `STATIC_DIR`, the calculator API, OpenAPI JSON, and health
checks on one same-origin port. Vite's development proxy is only for local UI
development.

```mermaid
flowchart LR
  subgraph Browser
    CalculatorUI[Calculator UI]
    Swagger[Bundled Swagger UI]
  end
  Browser -->|GET / and /docs/ assets| Go[Go HTTP server]
  CalculatorUI -->|POST /api/calculate| Go
  Swagger -->|GET /api/openapi.json| Go
  Swagger -->|Try it out: POST /api/calculate| Go
  Go --> Static[STATIC_DIR: Vite dist]
  Go --> Calculator[Finite binary64 calculator]
```

```mermaid
sequenceDiagram
  participant U as Browser or Swagger UI
  participant S as Go server
  participant C as Calculator
  U->>S: POST /api/calculate JSON
  S->>S: Parse, validate shape and finite operands
  alt validation failure
    S-->>U: 4xx JSON error
  else valid request
    S->>C: Calculate(operation, a, b)
    C-->>S: finite result or domain error
    S-->>U: 200 result or 422 JSON error
  end
```

There is no database, authentication, queue, cache, or background worker:
calculation is synchronous, deterministic for the documented binary64 inputs,
and has no side effects. [DECISIONS.md](DECISIONS.md) explains the measured
scope decision; [SPEC.md](SPEC.md) and [NUMERIC-VALIDATION.md](NUMERIC-VALIDATION.md)
define the behavior.
