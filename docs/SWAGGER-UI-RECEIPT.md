# Swagger UI receipt

## Delivery

`frontend/documentation/swagger.html` is the tracked local Swagger UI template.
`frontend/scripts/build-api-docs.mjs` recreates ignored `frontend/public/docs/`
from that template and the pinned `swagger-ui-dist@5.33.2` files, including
the bundle license reference and package `LICENSE` and `NOTICE`. Vite copies
the resulting local assets into `frontend/dist/docs/` for production serving.

The viewer requests the same-origin `/api/openapi.json`, sets
`validatorUrl: null`, enables Try it out for GET and POST, and exposes server
response status and bodies in the standard Swagger UI response panel. It has a
link back to `/`.

## Verification

Run from `frontend/` on 2026-10-15:

| Command                          | Result                                                                                                  |
| -------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `pnpm install --frozen-lockfile` | passed; lockfile is current                                                                             |
| `pnpm audit --prod`              | passed; no known vulnerabilities                                                                        |
| `pnpm build:api-docs`            | passed; recreated local Swagger HTML, CSS, JavaScript, and favicons                                     |
| `pnpm lint`                      | passed with zero warnings                                                                               |
| `pnpm typecheck`                 | passed                                                                                                  |
| `pnpm test -- --run`             | passed: 2 files, 8 tests                                                                                |
| `pnpm build`                     | passed; `dist/docs/` contains `index.html`, `swagger-ui.css`, `swagger-ui-bundle.js`, and both favicons |
| local-reference scan             | passed; template, generated public HTML, and built HTML contain no `http://` or `https://` reference    |

The browser acceptance handoff is `/docs/`. Useful selectors after Swagger UI
loads are `.swagger-ui`, `.opblock-summary-description`, `.try-out__btn`,
`.execute-wrapper .btn.execute`, and `.responses-wrapper`.

## Later dependency check

The follow-up coordinator ran `pnpm audit` for all dependencies on 2026-10-15
and reported zero known vulnerabilities. This supplements the production-only
check above; it is not a replacement for the final integrated verification.

## Sanitized implementation brief

> Bundle a same-origin Swagger UI at `/docs/` using a pinned local
> `swagger-ui-dist` package. Do not use a runtime CDN or validator. Generate
> ignored public assets from a tracked template before development and builds;
> load `/api/openapi.json`, provide a calculator link, and support Try it out
> and visible API errors.
