.PHONY: help proto-generate install-deps test-unit test-unit-user test-unit-menu test-unit-order test-integration test-e2e test test-all test-coverage test-e2e-docker docker-build docker-up docker-down docker-logs ci-test ci-full dev-setup

help:
	@echo "Available commands:"
	@echo "  make proto-generate      - Generate protobuf code"
	@echo "  make install-deps        - Install testing dependencies"
	@echo "  make test-unit          - Run all unit tests"
	@echo "  make test-unit-user     - Run user service unit tests"
	@echo "  make test-unit-menu     - Run menu service unit tests"
	@echo "  make test-unit-order    - Run order service unit tests"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-e2e           - Run E2E tests (services must be running)"
	@echo "  make test-e2e-docker    - Start services, run E2E tests, stop services"
	@echo "  make test               - Run unit and integration tests"
	@echo "  make test-all           - Run all tests including E2E"
	@echo "  make test-coverage      - Generate coverage reports"
	@echo "  make docker-build       - Build Docker images"
	@echo "  make docker-up          - Start all services"
	@echo "  make docker-down        - Stop all services"
	@echo "  make docker-logs        - Show logs from all services"
	@echo "  make ci-test            - Run tests suitable for CI"
	@echo "  make ci-full            - Full CI pipeline"
	@echo "  make dev-setup          - Complete dev setup (deps + proto + docker)"

proto-generate:
	@echo "=== Generating protobuf code ==="
	@powershell -Command "if (-not (Get-Command protoc -ErrorAction SilentlyContinue)) { Write-Host 'Error: protoc not found. Please install Protocol Buffers compiler.' -ForegroundColor Red; exit 1 }"
	cd proto/user/v1 && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
	cd proto/menu/v1 && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative menu.proto
	cd proto/order/v1 && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative order.proto
	@echo "Protobuf code generated successfully"

install-deps:
	@echo "=== Installing dependencies ==="
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	cd user-service && go mod tidy
	cd menu-service && go mod tidy
	cd order-service && go mod tidy
	cd api-gateway && go mod tidy
	cd tests/integration && go mod tidy
	cd tests/e2e && go mod tidy
	@echo "Dependencies installed successfully"

test-unit:
	@echo "=== Running all unit tests ==="
	@$(MAKE) test-unit-user
	@$(MAKE) test-unit-menu
	@$(MAKE) test-unit-order

test-unit-user:
	@echo "=== User Service Unit Tests ==="
	cd user-service && go test -v ./grpc/...

test-unit-menu:
	@echo "=== Menu Service Unit Tests ==="
	cd menu-service && go test -v ./grpc/...

test-unit-order:
	@echo "=== Order Service Unit Tests ==="
	cd order-service && go test -v ./grpc/...

test-integration:
	@echo "=== Running integration tests ==="
	cd tests/integration && go test -v ./...

test-e2e:
	@echo "=== Running E2E tests ==="
	cd tests/e2e && go test -v ./...

test-e2e-docker:
	@echo "=== Running E2E tests with Docker ==="
	@$(MAKE) docker-up
	@powershell -Command "Start-Sleep -Seconds 15"
	@$(MAKE) test-e2e
	@$(MAKE) docker-down

test:
	@echo "=== Running unit and integration tests ==="
	@$(MAKE) test-unit
	@$(MAKE) test-integration

test-all:
	@echo "=== Running all tests ==="
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@$(MAKE) test-e2e

test-coverage:
	@echo "=== Generating coverage reports ==="
	cd user-service && go test -coverprofile=coverage.out ./grpc/... && go tool cover -html=coverage.out -o coverage.html
	cd menu-service && go test -coverprofile=coverage.out ./grpc/... && go tool cover -html=coverage.out -o coverage.html
	cd order-service && go test -coverprofile=coverage.out ./grpc/... && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage reports generated: */coverage.html"

docker-build:
	@echo "=== Building Docker images ==="
	docker compose build

docker-up:
	@echo "=== Starting services ==="
	docker compose up -d
	@echo "Services started. Waiting for initialization..."

docker-down:
	@echo "=== Stopping services ==="
	docker compose down

docker-logs:
	@echo "=== Service logs ==="
	docker compose logs -f

ci-test:
	@echo "=== Running CI tests ==="
	@$(MAKE) test-unit
	@$(MAKE) test-integration

ci-full:
	@echo "=== Running full CI pipeline ==="
	@$(MAKE) install-deps
	@$(MAKE) proto-generate
	@$(MAKE) test-unit
	@$(MAKE) test-integration
	@$(MAKE) docker-build
	@$(MAKE) test-e2e-docker

dev-setup:
	@echo "=== Setting up development environment ==="
	@$(MAKE) install-deps
	@$(MAKE) proto-generate
	@$(MAKE) docker-build
	@echo "Development environment ready!"
