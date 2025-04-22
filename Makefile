# Paths
APP_NAME := social-app
MAIN_PKG := github.com/iammrsea/social-app/cmd
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

# Default environment file
ENV_FILE := .env

# The Go binary to use
GO := go

# Delve config
DLV := dlv
DLV_PORT := 2345

# The main Go test and build targets
TEST := $(GO) test -v
BUILD := $(GO) build

# Test tags
UNIT_TAG := unit
INTEGRATION_TAG := integration

SRC_DIR := ./internal

.PHONY: run build debug clean fmt tidy dev generate test-unit test-integration test

## === Commands ===

run:
	$(GO) run $(MAIN_PKG)

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(MAIN_PKG)

debug:
	$(DLV) debug $(MAIN_PKG) \
		--headless \
		--listen=:$(DLV_PORT) \
		--api-version=2 \
		--log

debug-tags:
	$(DLV) debug -tags=debug $(MAIN_PKG) \
		--headless \
		--listen=:$(DLV_PORT) \
		--api-version=2 \
		--log

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR)

dev:
	air

generate:
	$(GO) generate github.com/iammrsea/social-app/cmd/server/graphql

test-unit:
	$(GO) test $(SRC_DIR)/... -tags $(UNIT_TAG)

test-integration:
	$(GO) test $(SRC_DIR)/... -tags $(INTEGRATION_TAG)

test:
	$(GO) test $(SRC_DIR)/...

coverage:
	$(GO) test $(SRC_DIR)/... -coverprofile=coverage.out && $(GO) tool cover -html=coverage.out
