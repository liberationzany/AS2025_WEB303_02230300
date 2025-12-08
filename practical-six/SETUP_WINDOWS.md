# Setup Instructions for Windows

This guide helps you set up the project on Windows to run all tests successfully.

## Prerequisites Installation

### 1. Install Go
Download and install Go 1.23+ from: https://go.dev/dl/

### 2. Install Protocol Buffers Compiler (protoc)

**Option A: Using Chocolatey (Recommended)**
```powershell
choco install protoc
```

**Option B: Manual Installation**
1. Download from: https://github.com/protocolbuffers/protobuf/releases
2. Extract to `C:\protoc`
3. Add `C:\protoc\bin` to PATH

### 3. Install MinGW (GCC for Windows) - Required for SQLite

**Option A: Using Chocolatey**
```powershell
choco install mingw
```

**Option B: Manual Installation**
1. Download from: https://sourceforge.net/projects/mingw-w64/
2. Install and add to PATH

Verify installation:
```powershell
gcc --version
```

### 4. Install Docker Desktop
Download from: https://www.docker.com/products/docker-desktop/

## Project Setup

### 1. Clone/Extract the Project
```powershell
cd "C:\Users\zeroe\OneDrive\Desktop\practicals Y3S1\practical-six"
```

### 2. Install Go Tools
```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 3. Generate Protobuf Code
```powershell
# User service
cd proto\user\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto

# Menu service
cd ..\..\..\proto\menu\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative menu.proto

# Order service
cd ..\..\..\proto\order\v1
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative order.proto

cd ..\..\..
```

### 4. Download Dependencies
```powershell
cd user-service
go mod tidy

cd ..\menu-service
go mod tidy

cd ..\order-service
go mod tidy

cd ..\api-gateway
go mod tidy

cd ..\tests\integration
go mod tidy

cd ..\e2e
go mod tidy

cd ..\..
```

## Running Tests

### Unit Tests (Require GCC/MinGW)

Make sure GCC is installed and in PATH, then:

```powershell
# Enable CGO (required for SQLite)
$env:CGO_ENABLED = "1"

# User Service
cd user-service
go test -v ./grpc/...

# Menu Service
cd ..\menu-service
go test -v ./grpc/...

# Order Service
cd ..\order-service
go test -v ./grpc/...

cd ..
```

### Integration Tests

```powershell
cd tests\integration
$env:CGO_ENABLED = "1"
go test -v ./...
cd ..\..
```

### E2E Tests

First, start the services:
```powershell
docker compose up -d
```

Wait 15 seconds for services to initialize, then:
```powershell
cd tests\e2e
go test -v ./...
cd ..\..
```

Stop services:
```powershell
docker compose down
```

## Troubleshooting

### CGO Error
**Error**: `Binary was compiled with 'CGO_ENABLED=0'`

**Solution**: Make sure MinGW/GCC is installed and enable CGO:
```powershell
$env:CGO_ENABLED = "1"
gcc --version  # Should show version
```

### GCC Not Found
**Error**: `exec: "gcc": executable file not found`

**Solution**: Install MinGW and add to PATH, or use Chocolatey:
```powershell
choco install mingw
```

### Protoc Not Found
**Error**: `protoc: command not found`

**Solution**: Install Protocol Buffers:
```powershell
choco install protoc
```

### Port Already in Use
**Error**: Port 8080/5005x already in use

**Solution**: Stop existing services:
```powershell
docker compose down
```

Or find and kill the process:
```powershell
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess | Stop-Process -Force
```

### Module Not Found
**Error**: `cannot find module`

**Solution**: Make sure you ran `go mod tidy` in all directories and generated proto files.

## Running Tests - Alternative Without CGO

If you cannot install MinGW/GCC, you can still run integration and E2E tests which don't require CGO:

```powershell
# Skip unit tests, run integration tests
cd tests\integration
go test -v ./...

# Run E2E tests (services must be running)
cd ..\e2e
go test -v ./...
```

## Expected Test Results

### Unit Tests
All services should show:
```
=== RUN   TestCreateUser
=== RUN   TestGetUser
=== RUN   TestGetUsers
--- PASS: TestCreateUser
--- PASS: TestGetUser
--- PASS: TestGetUsers
PASS
```

### Integration Tests
```
=== RUN   TestIntegration_CreateUser
=== RUN   TestIntegration_CompleteOrderFlow
=== RUN   TestIntegration_OrderValidation
=== RUN   TestIntegration_ConcurrentOrders
--- PASS: All tests
PASS
```

### E2E Tests
```
Waiting for services...
Services ready!
=== RUN   TestE2E_CompleteOrderFlow
=== RUN   TestE2E_OrderValidation
--- PASS: All tests
PASS
```

## Coverage Reports

Generate coverage (requires CGO):
```powershell
$env:CGO_ENABLED = "1"

cd user-service
go test -coverprofile=coverage.out ./grpc/...
go tool cover -html=coverage.out -o coverage.html
start coverage.html

cd ..\menu-service
go test -coverprofile=coverage.out ./grpc/...
go tool cover -html=coverage.out -o coverage.html
start coverage.html

cd ..\order-service
go test -coverprofile=coverage.out ./grpc/...
go tool cover -html=coverage.out -o coverage.html
start coverage.html
```

## Submission Checklist

- [ ] MinGW/GCC installed (for unit tests)
- [ ] All unit tests passing
- [ ] Integration tests passing
- [ ] E2E tests passing (with Docker)
- [ ] Screenshots of test results
- [ ] All files committed

## Quick Start Script

Create a file `run-tests.ps1`:

```powershell
# Set CGO
$env:CGO_ENABLED = "1"

Write-Host "=== Running User Service Tests ===" -ForegroundColor Green
cd user-service
go test -v ./grpc/...

Write-Host "`n=== Running Menu Service Tests ===" -ForegroundColor Green
cd ..\menu-service
go test -v ./grpc/...

Write-Host "`n=== Running Order Service Tests ===" -ForegroundColor Green
cd ..\order-service
go test -v ./grpc/...

Write-Host "`n=== Running Integration Tests ===" -ForegroundColor Green
cd ..\tests\integration
go test -v ./...

Write-Host "`n=== Starting Docker Services ===" -ForegroundColor Green
cd ..\..
docker compose up -d

Write-Host "Waiting for services to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

Write-Host "`n=== Running E2E Tests ===" -ForegroundColor Green
cd tests\e2e
go test -v ./...

Write-Host "`n=== Stopping Docker Services ===" -ForegroundColor Green
cd ..\..
docker compose down

Write-Host "`n=== All Tests Complete! ===" -ForegroundColor Green
```

Run with:
```powershell
.\run-tests.ps1
```
