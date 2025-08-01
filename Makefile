include .env
LOCAL_BIN:=$(CURDIR)/bin

ENV_FILE=.env

MIGRATION_DIR := $(shell grep MIGRATION_DIR $(ENV_FILE) | cut -d '=' -f2)

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0

migration-add:
	goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v