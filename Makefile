# Paths
APP_NAME := social-app
MAIN_PKG := github.com/iammrsea/social-app/cmd
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

# Default environment file
ENV_FILE := .env

# Delve config
DLV := dlv
DLV_PORT := 2345

.PHONY: run build debug clean fmt tidy

## === Commands ===

run:
	go run $(MAIN_PKG)

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(MAIN_PKG)

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
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)

dev:
	air

generate:
	go generate github.com/iammrsea/social-app/cmd/server/graphql