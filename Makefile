.PHONY: install dev test lint typecheck build api e2e docker-build docker-up docker-down

install:
	cd frontend && pnpm install --frozen-lockfile

dev:
	@echo "Run 'make api' and 'cd frontend && pnpm dev' in separate terminals."

test:
	cd backend && go test -race -count=1 -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | awk '/^total:/{seen=1; if ($$3 != "100.0%") exit 1} END {if (!seen) exit 1}' && go vet ./...
	cd frontend && pnpm test:coverage

lint:
	cd frontend && pnpm lint

typecheck:
	cd frontend && pnpm typecheck

build:
	cd backend && go build ./cmd/server
	cd frontend && pnpm build

api:
	cd backend && go run ./cmd/server

e2e:
	cd frontend && pnpm test:e2e

docker-build:
	docker build --tag calculator:local .

docker-up:
	docker compose up --build

docker-down:
	docker compose down --remove-orphans
