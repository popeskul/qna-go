.PHONY: postgres adminer migrate-install migrate-up migrate-down

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

migrate-install:
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash

migrate-up:
	migrate -path ./schema -database 'postgres://$(db_user):$(db_password)@$(db_host):$(db_port)/$(db_name)?sslmode=disable' -verbose up

migrate-down:
	migrate -path ./schema -database 'postgres://$(db_user):$(db_password)@$(db_host):$(db_port)/$(db_name)?sslmode=disable' down -all

postgres:
	docker run --rm -ti --network host -e POSTGRES_PASSWORD=12345 -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres postgres

api:
	docker run --rm -ti --network host qna-go

# rules for compiling the golangci-lint
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.46.2

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml
