.DEFAULT_GOAL := run

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# This here is done to allow goose to run migrations, since it won't work using db as the value in a development environment
ifeq ($(ENV),production)
    PGHOST_OVERRIDE=$(PGHOST)
else
    PGHOST_OVERRIDE=localhost
endif

PG_CONN_STRING=postgresql://$(PGUSER):$(PGPASSWORD)@$(PGHOST_OVERRIDE):$(PGPORT)/$(PGDATABASE)?sslmode=disable

# API
vet: 
	go vet ./... 
	go fmt ./...
	staticcheck ./...
	gosec ./...
.PHONY:fmt

vendor:
	go mod vendor
.PHONY: vendor

build: vendor
	go build -mod=vendor -o luso-wiki
.PHONY: build

run:
	air
.PHONY:run

test:
	go test ./... -count=1
.PHONY: test

m-create:
ifndef name
	$(error name is required, e.g., `make m-create name=article_alter_table_add_column_content`)
endif
	goose create -s -dir internal/database/sql/postgres/migrations $(name) sql
.PHONY: m-create

m-up:
	goose -dir internal/database/sql/postgres/migrations postgres "$(PG_CONN_STRING)" up
.PHONY: m-up

m-down:
	goose -dir internal/database/sql/postgres/migrations postgres "$(PG_CONN_STRING)" down
.PHONY: m-down

m-status:
	goose -dir internal/database/sql/postgres/migrations postgres "$(PG_CONN_STRING)" status
.PHONY: m-status
