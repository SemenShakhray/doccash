YAML_FILE=./config.yaml	

MIGRATION_DIR:=$(shell yq eval '.migration_dir' $(YAML_FILE))

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	
migration-add:
	goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} sqlite3 ${STORAGE_PATH} down -v