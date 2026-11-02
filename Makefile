MIGRATIONS_DIR=./migrations
APP_MAIN=./cmd/api/main.go
GO_VERSION=$(shell go list -m -f '{{.GoVersion}}')
LINT=GOTOOLCHAIN=go$(GO_VERSION) golangci-lint

LOCAL_COMPOSE_FILE=docker-compose.local.yaml
LOCAL_PROJECT=shop_api_local
LOCAL_COMPOSE=docker compose -f $(LOCAL_COMPOSE_FILE)
LOCAL_NETWORK=$(LOCAL_PROJECT)_default
MIGRATE_NETWORK=$(notdir $(CURDIR))_default

ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: help up down build restart logs ps migrate-up migrate-down api-restart \
	swag clean nginx-restart db-create-readonly-user db-build nginx-balancing-test \
	nginx-cache-test nginx-compression-test generate-local-cert test lint lint-fix fmt \
	local-up local-down local-build local-restart local-logs local-ps local-clean \
	local-migrate-up local-migrate-down local-config

help:
	@echo "How to build:"
	@echo "  1. make swag                       - generate swagger docs"
	@echo "  2. make db-build                   - build and start database"
	@echo "  3. make migrate-up                 - apply migrations"
	@echo "  4. make db-create-readonly-user    - create readonly user"
	@echo "  5. make generate-local-cert        - generate local https certificate"
	@echo "  6. make build                      - rebuild and start all services"
	@echo "  ----------------------------"
	@echo "Other available commands:"
	@echo "  make up                         - start all services"
	@echo "  make down                       - stop all services"
	@echo "  make restart                    - restart all services"
	@echo "  make api-restart                - restart api services"
	@echo "  make nginx-restart              - restart nginx service"
	@echo "  make nginx-balancing-test       - nginx balancing test"
	@echo "  make nginx-cache-test           - nginx cache test"
	@echo "  make nginx-compression-test     - nginx compression test"
	@echo "  make migrate-down               - rollback last migration"
	@echo "  make test                       - run integration tests (requires docker)"
	@echo "  make lint                       - run golangci-lint"
	@echo "  make lint-fix                   - apply automatic lint fixes"
	@echo "  make fmt                        - format code (gofmt + goimports)"
	@echo "  make logs                       - show logs"
	@echo "  make ps                         - show containers"
	@echo "  make clean                      - remove containers and volumes"
	@echo "  ----------------------------"
	@echo "Local HTTP stack (http://localhost:8080):"
	@echo "  make local-build                - build and start local stack"
	@echo "  make local-up                   - start local stack"
	@echo "  make local-down                 - stop local stack"
	@echo "  make local-restart              - rebuild and restart local stack"
	@echo "  make local-migrate-up           - apply migrations to local database"
	@echo "  make local-migrate-down         - rollback last local migration"
	@echo "  make local-logs                 - show local logs"
	@echo "  make local-ps                   - show local containers"
	@echo "  make local-clean                - remove local containers and volumes"
	@echo "  ----------------------------"

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose up -d --build

restart:
	docker compose down
	docker compose up -d --build

nginx-restart:
	docker compose down nginx
	docker compose up -d nginx

nginx-balancing-test:
	sh scripts/nginx_balancing_test.sh

nginx-cache-test:
	sh scripts/nginx_cache_test.sh

nginx-compression-test:
	sh scripts/nginx_compression_test.sh

api-restart:
	docker compose down api api-read-1 api-read-2
	docker compose up -d api api-read-1 api-read-2

generate-local-cert:
	sh scripts/generate_local_cert.sh

logs:
	docker compose logs -f

ps:
	docker compose ps

db-build:
	docker compose up -d --build postgres

db-create-readonly-user:
	docker compose exec postgres psql \
		-U $(POSTGRES_USER) \
		-d $(POSTGRES_DB) \
		-v ro_user="$(POSTGRES_READONLY_USER)" \
		-v ro_password="$(POSTGRES_READONLY_PASSWORD)" \
		-f /scripts/create_readonly_user.sql

migrate-up:
	docker run --rm \
		--network $(MIGRATE_NETWORK) \
		-v "$$(pwd)/$(MIGRATIONS_DIR):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_DSN)" \
		up

migrate-down:
	docker run --rm \
		--network $(MIGRATE_NETWORK) \
		-v "$$(pwd)/$(MIGRATIONS_DIR):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_DSN)" \
		down 1

swag:
	swag init -g $(APP_MAIN)

test:
	go test -count=1 -timeout 15m -tags integration ./internal/integration/...

lint:
	$(LINT) run ./...

lint-fix:
	$(LINT) run --fix ./...

fmt:
	$(LINT) fmt

local-config:
	$(LOCAL_COMPOSE) config --quiet

local-up:
	$(LOCAL_COMPOSE) up -d

local-build:
	$(LOCAL_COMPOSE) up -d --build

local-down:
	$(LOCAL_COMPOSE) down

local-restart:
	$(LOCAL_COMPOSE) down
	$(LOCAL_COMPOSE) up -d --build

local-logs:
	$(LOCAL_COMPOSE) logs -f

local-ps:
	$(LOCAL_COMPOSE) ps

local-clean:
	$(LOCAL_COMPOSE) down -v --remove-orphans

local-migrate-up:
	$(MAKE) migrate-up MIGRATE_NETWORK=$(LOCAL_NETWORK)

local-migrate-down:
	$(MAKE) migrate-down MIGRATE_NETWORK=$(LOCAL_NETWORK)

clean:
	docker compose down -v --remove-orphans
