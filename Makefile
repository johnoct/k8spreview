.PHONY: build test clean release run dev

# Build variables
BINARY_NAME=k8spreview
VERSION=$(shell grep -E "Version = \".*\"" pkg/version/version.go | cut -d'"' -f2)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X k8spreview/pkg/version.Version=${VERSION} -X k8spreview/pkg/version.Commit=${COMMIT} -X k8spreview/pkg/version.Date=${BUILD_DATE}"

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BINARY_NAME} cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -f cassette.tape
	@rm -f coverage.out
	@rm -rf dist/
	@rm -rf tmp/

# Run the application
run: build
	@./${BINARY_NAME} examples/multi-resource.yaml

# Run with hot reload using Air
dev:
	@air -c .air.toml

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Verify before release
verify: fmt lint test

# Create a new release
release:
	@if [ "${v}" = "" ]; then \
		echo "Please specify version: make release v=x.x.x"; \
		exit 1; \
	fi
	@echo "Creating release v${v}..."
	@echo "Updating version in version.go..."
	@sed -i '' 's/Version = ".*"/Version = "v${v}"/' pkg/version/version.go
	@git add pkg/version/version.go
	@git commit -m "chore: bump version to v${v}"
	@git tag -a "v${v}" -m "Release v${v}"
	@git push origin main "v${v}"
	@echo "Release v${v} created and pushed!"

# Show help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make run           - Build and run the application"
	@echo "  make dev           - Run with hot reload"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make verify        - Run format, lint and tests"
	@echo "  make release v=x.x.x - Create and push a new release"
	@echo "  make help          - Show this help message"
