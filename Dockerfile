# syntax=docker/dockerfile:1
FROM node:22.23.2-bookworm-slim AS frontend-build
WORKDIR /src/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.27.2-bookworm AS backend-build
WORKDIR /src/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=backend-build --chown=app:app /out/server /app/server
COPY --from=frontend-build --chown=app:app /src/frontend/dist /app/static
USER app
ENV PORT=8080 STATIC_DIR=/app/static
EXPOSE 8080
ENTRYPOINT ["/app/server"]
