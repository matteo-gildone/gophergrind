MODULE  := $(shell go list -m)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test test-coverage clean clean-all install lint check pkgdev help

.DEFAULT_GOAL := help

help: ## Show available targets
	@echo "Gopher Grinder Makefile"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run tests with race detector
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "→ Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

clean: ## Clean build artifacts
	@rm -rf bin/ dist/ coverage.out coverage.html
	@go clean -cache -testcache

clean-all: clean ## Also clean module cache (WARNING: affects all Go projects)
	@echo "Cleaning module cache..."
	@go clean -modcache

lint: ## Check formatting and vet code
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	@go vet ./...

format: ## Format code
	@gofmt -w .

check: lint test ## Run all checks

release: ## Check tests pass and working tree is clean before releasing
	@go test -race ./...
	@git diff --exit-code || (echo "Working tree is dirty — commit or stash changes first"; exit 1)
	@echo "Ready to release. Run:"
	@echo "  git tag <version>"
	@echo "  git push origin <version>"

pkgdev: ## Trigger pkg.go.dev indexing
	@test -n "$(VERSION)" || (echo "No git tag found"; exit 1)
	@echo "https://proxy.golang.org/$(MODULE)/@v/$(VERSION).info"
	@echo "→ https://pkg.go.dev/$(MODULE)@$(VERSION)"