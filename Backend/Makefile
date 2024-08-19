SHELL := /bin/bash

include .env

# --- HELPERS ---
## help: display this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## help/api: display API usage
.PHONY: help/api
help/api:
	@go run ./cmd/api --help

# confirm: display confirmation prompt
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# --- DEVELOPMENT ---
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo "building cmd/api..."
	# build for local machine
	go build ./cmd/api

.PHONY: audit
audit: vendor
	@echo "formatting codebase..."
	@go fmt ./...

	@echo "vetting code..."
	@go vet ./...
	@staticcheck ./...

	@echo "running tests..."
	@go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo "tidying and verifying module dependencies..."
	@go mod tidy
	@go mod verify

	@echo "vendoring dependencies..."
	@go mod vendor

# --- STAGING---
remote = ${STAGING_USER}@${STAGING_HOST}
remote_dir = ${remote}:~/project/

## staging/connect: connect to the production server
.PHONY: staging/connect
staging/connect:
	ssh -i ${STAGING_KEY_PATH} ${remote}

# --- DATABASE MIGRATION ---
## db/migrations/create name=$1: create a new database migration
.PHONY: db/migrations/create
db/migrations/create: confirm
	@echo "creating migration files for ${name}..."
	migrate create -seq -ext .sql -dir ./migrations ${name}

## db/migrations/version: display current database migration version
.PHONY: db/migrations/version
db/migrations/version:
	@echo -n "database migration version: "
	@migrate -path ./migrations -database ${DB_URI} version

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: 
	@echo "running up migrations..."
	@migrate -path ./migrations -database ${DB_URI} up

## db/migrations/reset: rollback all database migrations
.PHONY: db/migrations/reset
db/migrations/reset: confirm
	@echo "running down migrations..."
	@migrate -path ./migrations -database ${DB_URI} down

## db/migrations/goto version=$1: move to a specified database migration version
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	@echo "migrating database to version ${version}..."
	@migrate -path ./migrations -database ${DB_URI} goto ${version}

## db/migrations/force version=$1: force a migration to the given version
.PHONY: db/migrations/force
db/migrations/force: confirm
	@echo "migrating database to version ${version}..."
	@migrate -path ./migrations -database ${DB_URI} force ${version}
