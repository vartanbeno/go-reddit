ROOT_PKG ?= "github.com/vartanbeno/geddit"
LIST_PKG := $(shell go list $(ROOT_PKG)/...)

# Tests
TEST_TIMEOUT ?= 20

.DEFAULT_GOAL := usage

# Print colorized log
define log
	echo "\n\033[1;32m--- [$(@)] $(1) ---\033[0m\n"
endef

all: lint fmt vet test test-coverage build

usage:
	@echo "make [all|fmt|vet|lint|test|test-coverage]"

fmt:
	@$(call log,"Running formatter")
	@go fmt $(LIST_PKG)

vet:
	@$(call log,"Running vet")
	@go vet -all $(LIST_PKG)

lint:
	@$(call log,"Running linter")
	@golint -set_exit_status $(LIST_PKG)

test: fmt vet lint
	@$(call log,"Running tests")
	@go test -v -race -short -timeout $(TEST_TIMEOUT)s $(ARGS) $(LIST_PKG)

test-coverage: fmt vet lint
	@$(call log,"Running tests with coverage")
	@go test -v -race -short -timeout $(TEST_TIMEOUT)s $(ARGS) -coverprofile=coverage.out $(LIST_PKG)
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
