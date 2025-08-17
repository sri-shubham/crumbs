.PHONY: test benchmark coverage lint clean

# Default target
all: test lint benchmark

# Run all tests
test:
	go test -v ./...

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./...

# Generate test coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# Clean up generated files
clean:
	rm -f coverage.out coverage.html
	
# Show help
help:
	@echo "Available targets:"
	@echo "  all        : Run tests, linter and benchmarks"
	@echo "  test       : Run all tests"
	@echo "  benchmark  : Run all benchmarks"
	@echo "  coverage   : Generate test coverage report"
	@echo "  lint       : Run golangci-lint"
	@echo "  clean      : Remove generated files"
