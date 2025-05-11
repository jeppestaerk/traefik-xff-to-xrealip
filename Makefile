.PHONY: all lint test yaegi_test vendor clean

# Default target: runs lint and standard tests
default: lint test

# Lint the code
lint:
	@echo "==> Linting code..."
	golangci-lint run

# Run standard Go tests
test:
	@echo "==> Running standard Go tests..."
	go test -v -cover ./...

# Run tests with Yaegi
# This assumes your plugin's .traefik.yml is in the root and tests are in *_test.go files.
# Yaegi needs to interpret the plugin configuration and the test files.
yaegi_test:
	@echo "==> Running tests with Yaegi..."
	yaegi test -v .

# The 'vendor' directory contains project dependencies managed by Go modules.
vendor:
	@echo "==> Creating vendor directory..."
	go mod vendor

# Clean up (optional)
clean:
	@echo "==> Cleaning up..."
	go clean
	rm -rf ./vendor

# Help target (optional)
help:
	@echo "Available targets:"
	@echo "  all          : Run lint and standard tests (default)"
	@echo "  lint         : Lint the Go code"
	@echo "  test         : Run standard Go tests"
	@echo "  yaegi_test   : Run tests with Yaegi"
	@echo "  vendor       : Create vendor directory with dependencies"
	@echo "  clean        : Clean up build artifacts"
