.PHONY: lint format test tools

GOLANGCI_LINT ?= golangci-lint
# Match golangci/golangci-lint-action@v7 (see .github/workflows/lint.yml).
GOLANGCI_LINT_VERSION ?= v2.1.6

.DEFAULT_GOAL := lint

# The GitHub Actions "format" job runs golangci-lint at the repository root.
lint format:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { \
		echo "$(GOLANGCI_LINT) not found; install with: make tools"; \
		exit 1; \
	}
	$(GOLANGCI_LINT) run ./...

test:
	go test ./...

tools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
