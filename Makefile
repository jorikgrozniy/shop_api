MIGRATIONS_DIR=./migrations
APP_MAIN=./cmd/api/main.go

ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: help up down build restart logs ps \
	migrate-up migrate-down swag clean

help:
	@echo "Available commands:"
	@echo "  make up                         - start all services"
	@echo "  make down                       - stop all services"
	@echo "  make build                      - rebuild and start all services"
	@echo "  make restart                    - restart all services"
	@echo "  make logs                       - show logs"
	@echo "  make ps                         - show containers"
	@echo "  make migrate-up                 - apply migrations"
	@echo "  make migrate-down               - rollback last migration"
	@echo "  make swag                       - generate swagger docs"
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

logs:
	docker compose logs -f

ps:
	docker compose ps

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