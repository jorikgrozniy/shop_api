MIGRATIONS_DIR=./migrations
APP_MAIN=./cmd/api/main.go

ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: help up down build restart logs ps migrate-up migrate-down api-restart \
	swag clean nginx-restart db-create-readonly-user db-build nginx-balancing-test

help:
	@echo "Available commands:"
	@echo "  ----------------------------"
	@echo "  1. make swag                       - generate swagger docs"
	@echo "  2. make db-build                   - build and start database"
	@echo "  3. make migrate-up                 - apply migrations"
	@echo "  4. make db-create-readonly-user    - create readonly user"
	@echo "  5. make build                      - rebuild and start all services"
	@echo "  ----------------------------"
	@echo "  make up                         - start all services"
	@echo "  make down                       - stop all services"
	@echo "  make restart                    - restart all services"
	@echo "  make api-restart                - restart api services"
	@echo "  make nginx-restart              - restart nginx service"
	@echo "  make nginx-balancing-test       - nginx balancing test"
	@echo "  make logs                       - show logs"
	@echo "  make ps                         - show containers"
	@echo "  make migrate-down               - rollback last migration"
	@echo "  make clean                      - remove containers and volumes"

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

api-restart:
	docker compose down api
	docker compose down api-read-1
	docker compose down api-read-2
	docker compose up -d api
	docker compose up -d api-read-1
	docker compose up -d api-read-2

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
		--network $$(basename "$$(pwd)")_default \
		-v "$$(pwd)/$(MIGRATIONS_DIR):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_DSN)" \
		up

migrate-down:
	docker run --rm \
		--network $$(basename "$$(pwd)")_default \
		-v "$$(pwd)/$(MIGRATIONS_DIR):/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database "$(POSTGRES_DSN)" \
		down 1

swag:
	swag init -g $(APP_MAIN)

clean:
	docker compose down -v --remove-orphans