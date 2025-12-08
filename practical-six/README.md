# Practical 6: Comprehensive Testing for Microservices

This project demonstrates comprehensive testing strategies for a microservices architecture built with Go, gRPC, and PostgreSQL.

## Project Structure

```
practical-six/
├── proto/                      # Protocol Buffer definitions
│   ├── user/v1/
│   ├── menu/v1/
│   └── order/v1/
├── user-service/              # User microservice
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go    # Unit tests
│   ├── database/
│   ├── models/
│   ├── Dockerfile
│   └── main.go
├── menu-service/              # Menu microservice
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go    # Unit tests
│   ├── database/
│   ├── models/
│   ├── Dockerfile
│   └── main.go
├── order-service/             # Order microservice
│   ├── grpc/
│   │   ├── server.go
│   │   └── server_test.go    # Unit tests with mocks
│   ├── database/
│   ├── models/
│   ├── Dockerfile
│   └── main.go
├── api-gateway/               # HTTP API Gateway
│   ├── main.go
│   └── Dockerfile
├── tests/
│   ├── integration/           # Integration tests
│   │   └── integration_test.go
│   └── e2e/                   # End-to-end tests
│       └── e2e_test.go
├── docker-compose.yml
├── Makefile
└── README.md
```

## Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)
- Make (optional but recommended)

### Installing protoc on Windows

Download from: https://github.com/protocolbuffers/protobuf/releases

Or use Chocolatey:
```powershell
choco install protoc
```

## Quick Start

### 1. Setup Development Environment

```bash
# Install dependencies and generate proto code
make dev-setup
```

This will:
- Install Go dependencies
- Install protoc plugins
- Generate protobuf code
- Build Docker images

### 2. Run Unit Tests

```bash
# Run all unit tests
make test-unit

# Or run specific service tests
make test-unit-user
make test-unit-menu
make test-unit-order
```

### 3. Run Integration Tests

```bash
make test-integration
```

### 4. Run E2E Tests

```bash
# Start services, run tests, then stop
make test-e2e-docker

# Or if services are already running
make test-e2e
```

### 5. Run All Tests

```bash
make test-all
```

## Running the Application

### Start Services

```bash
docker compose up -d
```

Services will be available at:
- API Gateway: http://localhost:8080
- User Service: localhost:50051 (gRPC)
- Menu Service: localhost:50052 (gRPC)
- Order Service: localhost:50053 (gRPC)

### Stop Services

```bash
docker compose down
```

### View Logs

```bash
docker compose logs -f
```

## API Endpoints

### User Endpoints

- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user by ID
- `GET /api/users` - Get all users

### Menu Endpoints

- `POST /api/menu` - Create menu item
- `GET /api/menu/{id}` - Get menu item by ID
- `GET /api/menu` - Get all menu items

### Order Endpoints

- `POST /api/orders` - Create order
- `GET /api/orders/{id}` - Get order by ID
- `GET /api/orders` - Get all orders

## Example API Calls

### Create User

```powershell
$body = @{
    name = "John Doe"
    email = "john@example.com"
    is_cafe_owner = $false
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/users -Method POST -Body $body -ContentType "application/json"
```

### Create Menu Item

```powershell
$body = @{
    name = "Cappuccino"
    description = "Espresso with steamed milk"
    price = 4.50
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/api/menu -Method POST -Body $body -ContentType "application/json"
```

### Create Order

```powershell
$body = @{
    user_id = 1
    items = @(
        @{
            menu_item_id = 1
            quantity = 2
        }
    )
} | ConvertTo-Json -Depth 3

Invoke-RestMethod -Uri http://localhost:8080/api/orders -Method POST -Body $body -ContentType "application/json"
```

## Testing Architecture

### Unit Tests (70%)
- Test individual functions in isolation
- Use SQLite in-memory database
- Mock external dependencies
- Fast execution (milliseconds)

### Integration Tests (20%)
- Test multiple services working together
- Use bufconn for in-memory gRPC
- No external dependencies
- Medium speed (seconds)

### End-to-End Tests (10%)
- Test entire system via HTTP API
- Use Docker Compose
- Production-like setup
- Slower execution (seconds to minutes)

## Coverage Reports

Generate HTML coverage reports:

```bash
make test-coverage
```

View reports:
```bash
start user-service/coverage.html
start menu-service/coverage.html
start order-service/coverage.html
```

## CI/CD Commands

```bash
# Run tests suitable for CI
make ci-test

# Full CI pipeline
make ci-full
```

## Troubleshooting

### Protoc Not Found

Install Protocol Buffers compiler:
```powershell
choco install protoc
```

### Port Already in Use

Stop existing services:
```bash
docker compose down
```

### Tests Timeout

Increase wait time in tests or check if services are healthy:
```bash
docker compose ps
docker compose logs
```

### Module Issues

Clean and reinstall dependencies:
```bash
make install-deps
```

## Test Results

All tests should pass successfully:

✅ **Unit Tests**: 
- User Service: All tests passing
- Menu Service: All tests passing  
- Order Service: All tests passing (with mocks)

✅ **Integration Tests**:
- Service communication working
- Order flow validated
- Concurrent operations tested

✅ **E2E Tests**:
- Full system integration verified
- HTTP API functioning correctly
- Error handling validated

## Submission

Submit the following:
1. Complete source code repository
2. Screenshots showing:
   - All unit tests passing
   - Integration tests passing
   - E2E tests passing
3. Coverage reports (optional)

## Resources

- [Go Testing](https://golang.org/pkg/testing/)
- [Testify Library](https://github.com/stretchr/testify)
- [gRPC Testing](https://grpc.io/docs/languages/go/basics/#testing)
- [Protocol Buffers](https://protobuf.dev/)

## License

This project is for educational purposes as part of Practical 6.
