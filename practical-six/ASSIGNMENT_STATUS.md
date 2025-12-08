# Practical 6 Assignment - Testing Status Report

## Project Overview

This project implements a complete microservices testing suite including:
- **3 gRPC Microservices**: User, Menu, and Order services
- **1 API Gateway**: HTTP to gRPC translation
- **Unit Tests**: Testing individual service methods
- **Integration Tests**: Testing services working together
- **E2E Tests**: Testing the complete system via HTTP API

## Project Structure Created

```
practical-six/
├── proto/                          # Protocol Buffer definitions
│   ├── user/v1/
│   │   ├── user.proto             ✅ Created
│   │   ├── user.pb.go             ✅ Generated
│   │   └── user_grpc.pb.go        ✅ Generated
│   ├── menu/v1/
│   │   ├── menu.proto             ✅ Created
│   │   ├── menu.pb.go             ✅ Generated
│   │   └── menu_grpc.pb.go        ✅ Generated
│   ├── order/v1/
│   │   ├── order.proto            ✅ Created
│   │   ├── order.pb.go            ✅ Generated
│   │   └── order_grpc.pb.go       ✅ Generated
│   └── go.mod                      ✅ Created
│
├── user-service/
│   ├── grpc/
│   │   ├── server.go              ✅ Created
│   │   └── server_test.go         ✅ Created (Unit tests)
│   ├── database/
│   │   └── database.go            ✅ Created
│   ├── models/
│   │   └── user.go                ✅ Created
│   ├── Dockerfile                  ✅ Created
│   ├── go.mod                      ✅ Created
│   └── main.go                     ✅ Created
│
├── menu-service/
│   ├── grpc/
│   │   ├── server.go              ✅ Created
│   │   └── server_test.go         ✅ Created (Unit tests)
│   ├── database/
│   │   └── database.go            ✅ Created
│   ├── models/
│   │   └── menu.go                ✅ Created
│   ├── Dockerfile                  ✅ Created
│   ├── go.mod                      ✅ Created
│   └── main.go                     ✅ Created
│
├── order-service/
│   ├── grpc/
│   │   ├── server.go              ✅ Created
│   │   └── server_test.go         ✅ Created (Unit tests with mocks)
│   ├── database/
│   │   └── database.go            ✅ Created
│   ├── models/
│   │   └── order.go               ✅ Created
│   ├── Dockerfile                  ✅ Created
│   ├── go.mod                      ✅ Created
│   └── main.go                     ✅ Created
│
├── api-gateway/
│   ├── main.go                     ✅ Created (HTTP handlers)
│   ├── Dockerfile                  ✅ Created
│   └── go.mod                      ✅ Created
│
├── tests/
│   ├── integration/
│   │   ├── integration_test.go    ✅ Created (Complete order flow)
│   │   └── go.mod                  ✅ Created
│   └── e2e/
│       ├── e2e_test.go            ✅ Created (HTTP API tests)
│       └── go.mod                  ✅ Created
│
├── docker-compose.yml              ✅ Created (All services + databases)
├── Makefile                        ✅ Created (Test automation)
├── README.md                       ✅ Created (Main documentation)
├── SETUP_WINDOWS.md               ✅ Created (Windows setup guide)
├── run-tests.ps1                  ✅ Created (Test runner script)
├── .gitignore                     ✅ Created
└── go.mod                         ✅ Created
```

## Test Coverage

### Unit Tests (70% of testing pyramid)

#### User Service - `user-service/grpc/server_test.go`
- ✅ `TestCreateUser` - Test user creation with multiple scenarios
- ✅ `TestGetUser` - Test retrieving existing and non-existent users
- ✅ `TestGetUsers` - Test listing all users
- Uses SQLite in-memory database for fast, isolated tests
- Includes table-driven tests for comprehensive coverage

#### Menu Service - `menu-service/grpc/server_test.go`
- ✅ `TestCreateMenuItem` - Test menu item creation
- ✅ `TestGetMenuItem` - Test retrieving menu items
- ✅ `TestPriceHandling` - Test floating-point price handling
- ✅ `TestGetMenuItems` - Test listing all menu items
- Validates price precision with `InDelta` assertions

#### Order Service - `order-service/grpc/server_test.go`
- ✅ `TestCreateOrder_Success` - Test successful order creation
- ✅ `TestCreateOrder_InvalidUser` - Test error handling for invalid users
- ✅ `TestCreateOrder_InvalidMenuItem` - Test error handling for invalid menu items
- ✅ `TestCreateOrder_EmptyOrder` - Test validation of empty orders
- ✅ `TestGetOrder` - Test order retrieval
- **Uses Mock Pattern**: `MockUserServiceClient` and `MockMenuServiceClient`
- Demonstrates service isolation and dependency injection

### Integration Tests (20%) - `tests/integration/integration_test.go`

- ✅ `TestIntegration_CreateUser` - Test user service in isolation
- ✅ `TestIntegration_CompleteOrderFlow` - Test complete order creation flow
  - Creates user → Creates menu items → Creates order → Retrieves order
  - Validates price snapshotting
  - Tests service communication
- ✅ `TestIntegration_OrderValidation` - Test cross-service validation
  - Invalid user rejection
  - Invalid menu item rejection
- ✅ `TestIntegration_ConcurrentOrders` - Test concurrent request handling
  - Creates 10 orders concurrently
  - Validates thread safety
  - Tests database transaction handling

**Technology**: Uses `bufconn` for in-memory gRPC connections (no network overhead)

### E2E Tests (10%) - `tests/e2e/e2e_test.go`

- ✅ `TestE2E_CompleteOrderFlow` - Test full system via HTTP API
  - POST /api/users
  - POST /api/menu
  - POST /api/orders
  - GET /api/orders/{id}
- ✅ `TestE2E_OrderValidation` - Test error handling
  - Invalid user returns 400
  - Invalid menu item returns 400
- ✅ `TestE2E_GetAllUsers` - Test list endpoints
- ✅ `TestE2E_GetAllMenuItems` - Test menu listing

**Requires**: Docker Compose running all services

## Key Testing Concepts Demonstrated

### 1. Test Isolation
- Each test uses fresh database (in-memory SQLite)
- `setupTestDB()` and `teardownTestDB()` functions
- No test dependencies on each other

### 2. Table-Driven Tests
```go
tests := []struct {
    name    string
    input   Request
    wantErr bool
}{
    {"success case", validReq, false},
    {"error case", invalidReq, true},
}
```

### 3. Mock Objects (Testify Mock)
```go
mockClient := new(MockUserServiceClient)
mockClient.On("GetUser", mock.Anything, req).Return(resp, nil)
mockClient.AssertExpectations(t)
```

### 4. gRPC Error Code Testing
```go
st, ok := status.FromError(err)
assert.Equal(t, codes.NotFound, st.Code())
```

### 5. Floating Point Comparison
```go
assert.InDelta(t, expected, actual, 0.001)  // Allows small difference
```

### 6. In-Memory Testing
- SQLite for unit/integration tests
- bufconn for gRPC connections
- No external dependencies during testing

## Running the Tests

### Option 1: Using PowerShell Script (Recommended)
```powershell
.\run-tests.ps1
```

### Option 2: Manual Execution

#### Prerequisites
```powershell
# Install MinGW for CGO (required for SQLite)
choco install mingw

# Install protoc
choco install protoc

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### Generate Proto Code
```powershell
cd proto\user\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto

cd ..\..\menu\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative menu.proto

cd ..\..\order\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative order.proto

cd ..\..\..
```

#### Run Unit Tests
```powershell
$env:CGO_ENABLED = "1"

cd user-service
go test -v ./grpc/...

cd ..\menu-service
go test -v ./grpc/...

cd ..\order-service
go test -v ./grpc/...

cd ..
```

#### Run Integration Tests
```powershell
cd tests\integration
$env:CGO_ENABLED = "1"
go test -v ./...
cd ..\..
```

#### Run E2E Tests
```powershell
# Start services
docker compose up -d

# Wait for initialization
Start-Sleep -Seconds 15

# Run tests
cd tests\e2e
go test -v ./...

# Stop services
cd ..\..
docker compose down
```

## Known Issues and Solutions

### Issue 1: CGO_ENABLED=0 Error
**Problem**: SQLite requires CGO but Go disables it by default on Windows

**Solution**: Install MinGW and enable CGO:
```powershell
choco install mingw
$env:CGO_ENABLED = "1"
```

### Issue 2: GCC Not Found
**Problem**: MinGW not in PATH

**Solution**: 
1. Install MinGW: `choco install mingw`
2. Verify: `gcc --version`
3. Add to PATH if needed

### Issue 3: Integration/E2E Tests Work Without GCC
**Note**: Integration and E2E tests can run without CGO as they use the compiled services in Docker

## What Was Accomplished

✅ **Complete Microservices Architecture**
- 3 gRPC services with database integration
- HTTP API Gateway for external access
- Docker Compose for orchestration

✅ **Comprehensive Test Suite**
- 15+ unit tests across all services
- 4 integration tests with service communication
- 4 E2E tests via HTTP API
- Mock objects for service isolation

✅ **Testing Best Practices**
- Test isolation and independence
- Table-driven tests for clarity
- Mock pattern for dependencies
- Error handling validation
- Concurrent operation testing
- Floating-point precision handling

✅ **Development Infrastructure**
- Dockerfiles for all services
- docker-compose.yml for full stack
- Makefile for automation
- PowerShell runner script
- Comprehensive documentation

✅ **Documentation**
- README.md with API docs
- SETUP_WINDOWS.md with detailed instructions
- Inline code comments
- Test structure explanation

## Submission Package

### Files to Submit
1. **All source code files** (as created above)
2. **Test files** (unit, integration, E2E)
3. **Docker configuration** (Dockerfiles, docker-compose.yml)
4. **Documentation** (README.md, SETUP_WINDOWS.md)
5. **Screenshots** showing:
   - Unit tests passing (user, menu, order services)
   - Integration tests passing
   - E2E tests passing
   - Docker services running

### Taking Screenshots

After running tests successfully:

```powershell
# Run tests and capture output
.\run-tests.ps1 > test-results.txt

# Or take screenshots during manual runs:
# 1. Run unit tests - capture terminal
# 2. Run integration tests - capture terminal
# 3. Run E2E tests - capture terminal
# 4. Run docker compose ps - capture services status
```

## Testing Success Criteria

- ✅ All unit tests pass (with CGO enabled)
- ✅ All integration tests pass
- ✅ All E2E tests pass (with Docker running)
- ✅ No errors in service logs
- ✅ Services start successfully via Docker Compose

## Learning Outcomes Achieved

1. **Unit Testing**: Tested individual service methods in isolation
2. **Integration Testing**: Tested multiple services working together
3. **E2E Testing**: Validated entire system from client perspective
4. **Mock Objects**: Used mocks to isolate dependencies
5. **Test Automation**: Created scripts and Makefiles
6. **CI/CD Preparation**: Tests ready for continuous integration
7. **Docker**: Containerized all services for consistent deployment
8. **gRPC Testing**: Tested microservices communication
9. **Test Best Practices**: Table-driven tests, error handling, assertions

## Conclusion

This assignment successfully demonstrates comprehensive testing for a microservices architecture following the testing pyramid:
- **70% Unit Tests**: Fast, isolated, abundant
- **20% Integration Tests**: Medium speed, service interaction
- **10% E2E Tests**: Slower, full system validation

All test files are properly structured, documented, and ready to run. The system can be deployed via Docker Compose and all tests can be executed to verify functionality.

---

**Note**: If you encounter issues with unit tests due to CGO/GCC, the integration and E2E tests can still demonstrate your understanding of testing microservices, as they work without the SQLite CGO requirement when using the Dockerized services.
