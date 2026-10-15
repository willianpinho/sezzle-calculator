# Frontend coverage receipt

## Scope and configuration

The Vite/Vitest coverage configuration measures every handwritten runtime
TypeScript and TSX module in `frontend/src/`. It explicitly excludes only test
files and shared test setup. Each measured source module is imported through
behavioral unit or component tests, including `main.tsx`, whose test verifies
both mounting into `#root` and the missing-root failure.

Global V8 thresholds are now 100% for statements, branches, functions, and
lines. No coverage-ignore directive, deletion of error handling, or generated
Swagger asset instrumentation was used.

`frontend/scripts/build-api-docs.mjs` is included in the same V8 coverage
metric through a Node-environment Vitest suite. A standard ESM direct-execution
guard keeps helper imports side-effect-free. Tests cover both guard outcomes:
the helper import emits no CLI output, while a controlled `process.argv` test
runs the CLI entry and verifies its repository output. A missing script argument
is also harmless. Isolated temporary filesystem fixtures verify normal local
copying, stale-output replacement, and a missing vendored asset error. No
subprocess or third-party test dependency is used.

## Coverage behavior added

- Empty, finite, and non-finite number parsing.
- Every operation helper path, including the defensive unknown-operation path.
- Valid binary and square-root requests, local validation, percentage guidance,
  and API-visible errors.
- Non-`Error` network failures, silent aborts, and both late success and late
  error completions after an operand edit.
- API malformed responses and errors, request/header timeout, and response-body
  timeout coverage retained from the existing API tests.
- Swagger bundle copying, including direct-execution guard outcomes, normal
  local asset copy, stale output cleanup, and missing asset failure.

## Verification

Run from `frontend/` on 2026-10-15:

| Command                                             | Result                                                                                                 |
| --------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| `pnpm lint`                                         | passed with zero warnings                                                                              |
| `pnpm typecheck`                                    | passed                                                                                                 |
| `pnpm test -- --run`                                | passed: 5 files, 22 tests                                                                              |
| `pnpm test:coverage`                                | passed: 100% statements (99/99), branches (72/72), functions (20/20), and lines (95/95) across 6 files |
| `pnpm build`                                        | passed; includes local Swagger copy step and Vite production build                                     |
| `node scripts/build-api-docs.mjs` with asset checks | passed                                                                                                 |

ESLint ignores generated `test-results`, `playwright-report`, and `public/docs`
alongside `dist` and `coverage`, so concurrent Playwright cleanup cannot race
the lint file walker. These ignore paths do not affect V8 coverage collection.
