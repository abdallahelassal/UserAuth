include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATE=migrate -path ./migrations -database "$(DB_URL)"

.PHONY: migrate-up migrate-down migrate-create migrate-version migrate-force db-start

## Apply all pending migrations
migrate-up:
	$(MIGRATE) up

## Roll back last migration
migrate-down:
	$(MIGRATE) down 1

## Roll back ALL migrations (destructive!)
migrate-down-all:
	@echo "WARNING: This will drop all tables!"
	$(MIGRATE) down -all

## Create a new migration pair (usage: make migrate-create name=add_posts)
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

## Show current migration version
migrate-version:
	$(MIGRATE) version

## Force set version (escape dirty state)
migrate-force:
	$(MIGRATE) force $(version)

## Start postgres via Docker
db-start:
	docker compose up postgres -d

## Run the full stack
up:
	docker compose up --build