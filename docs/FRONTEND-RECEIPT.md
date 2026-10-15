# Frontend receipt

## Scope

Implemented the React/Vite calculator in `frontend/`. It uses the agreed
`POST /api/calculate` contract, sends no `b` value for square root, and relies
on the server for arithmetic validation. Inputs are labelled native controls;
the result is a polite live region. Editing an input cancels an active request
and clears its prior result. Requests have a 10-second timeout, including while
waiting for a JSON response body.

## Verification

Run from `frontend/` on 2026-10-14:

| Command                                  | Result                                                               |
| ---------------------------------------- | -------------------------------------------------------------------- |
| `pnpm install --frozen-lockfile`         | passed; lockfile current                                             |
| `pnpm lint`                              | passed with zero warnings                                            |
| `pnpm typecheck`                         | passed                                                               |
| `pnpm test -- --run`                     | passed: 2 files, 8 tests                                             |
| `pnpm build`                             | passed: Vite 8.3.4 production build                                  |
| `pnpm test:coverage`                     | passed: statements 100%, branches 85.71%, functions 100%, lines 100% |
| `node scripts/checks/run-all.mjs --diff` | unavailable: repository script did not exist when checked            |

Coverage artifacts are emitted as `frontend/coverage/lcov.info` and
`frontend/coverage/coverage-summary.json`. Unit tests include the success path,
square-root payload shape, server and local errors, stale-result clearing,
malformed API bodies, and timeouts before headers or during JSON parsing.

## Sanitized implementation brief

> Build the isolated React TypeScript calculator under `frontend/` with a
> polished, responsive, accessible English interface. Support add, subtract,
> multiply, divide, power, square root, and percentage. Call
> `POST /api/calculate` using `{ operation, a, b? }`, omit `b` for square root,
> and display `{ result }` or a safe API/network error. Use React 19.3.1,
> Vite 8.3.4, strict TypeScript, pnpm 10.33.0, Vitest and Testing Library;
> include lint, typecheck, build, unit-test, coverage, and E2E package scripts.
