.PHONY: build test test-unit test-vm lint clean fmt check all

# Build the project
build:
	go build -v ./...

# Run unit tests only  
test-unit:
	go test -v -race -short ./truenas/

# Run VM integration tests
test-vm:
	go test -v -timeout=30m -run TestClientWithVM ./truenas/

# Run all tests
test: test-unit

# Run all tests including VM tests
test-all: test-unit test-vm

# Run tests with coverage
test-coverage:
	go test -v -short -coverprofile=coverage.out ./truenas/
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint: fmt
	golangci-lint run ./...

# Format code
fmt:
	goimports -local github.com/715d/go-truenas -w .

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html coverage-vm.out coverage-vm.html

# Run all checks (format, lint, test)
check: fmt lint test

# Run full checks including VM tests
check-full: fmt lint test-all

# Validate VM setup
validate-vm-setup:
	@echo "Validating VM test setup..."
	@which qemu-system-x86_64 > /dev/null || (echo "ERROR: QEMU not found. Install with: brew install qemu" && exit 1)
	@which qemu-img > /dev/null || (echo "ERROR: qemu-img not found. Install with: brew install qemu" && exit 1)
	@qemu-system-x86_64 --version | head -1 | grep -q "version 1[0-9]\." || (echo "ERROR: QEMU version 10.x or higher required" && exit 1)
	@echo "VM setup is valid"

# Setup development environment
setup-dev:
	@echo "Setting up development environment..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest)
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@which qemu-system-x86_64 > /dev/null || echo "Warning: QEMU not found. VM tests will not work without QEMU installed."
	@echo "Development environment ready"

# Default target
all: check build
