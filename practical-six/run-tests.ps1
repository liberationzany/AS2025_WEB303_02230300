# Set CGO for Windows
$env:CGO_ENABLED = "1"

Write-Host "=== Practical 6 Test Runner ===" -ForegroundColor Cyan
Write-Host ""

# Check for GCC
try {
    $gccVersion = gcc --version 2>&1 | Select-Object -First 1
    Write-Host "GCC Found: $gccVersion" -ForegroundColor Green
} catch {
    Write-Host "WARNING: GCC not found. Unit tests will fail." -ForegroundColor Red
    Write-Host "Install MinGW with: choco install mingw" -ForegroundColor Yellow
    Write-Host "Continuing with integration and E2E tests only..." -ForegroundColor Yellow
    $skipUnit = $true
}

Write-Host ""

if (-not $skipUnit) {
    Write-Host "=== Running User Service Unit Tests ===" -ForegroundColor Green
    Set-Location user-service
    go test -v ./grpc/...
    if ($LASTEXITCODE -ne 0) { Write-Host "User service tests failed!" -ForegroundColor Red }
    Set-Location ..

    Write-Host "`n=== Running Menu Service Unit Tests ===" -ForegroundColor Green
    Set-Location menu-service
    go test -v ./grpc/...
    if ($LASTEXITCODE -ne 0) { Write-Host "Menu service tests failed!" -ForegroundColor Red }
    Set-Location ..

    Write-Host "`n=== Running Order Service Unit Tests ===" -ForegroundColor Green
    Set-Location order-service
    go test -v ./grpc/...
    if ($LASTEXITCODE -ne 0) { Write-Host "Order service tests failed!" -ForegroundColor Red }
    Set-Location ..
}

Write-Host "`n=== Running Integration Tests ===" -ForegroundColor Green
Set-Location tests\integration
go test -v ./...
if ($LASTEXITCODE -ne 0) { Write-Host "Integration tests failed!" -ForegroundColor Red }
Set-Location ..\..

Write-Host "`n=== Starting Docker Services ===" -ForegroundColor Green
docker compose up -d

Write-Host "Waiting for services to initialize (15 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

Write-Host "`n=== Running E2E Tests ===" -ForegroundColor Green
Set-Location tests\e2e
go test -v ./...
$e2eResult = $LASTEXITCODE
Set-Location ..\..

Write-Host "`n=== Stopping Docker Services ===" -ForegroundColor Green
docker compose down

Write-Host ""
if ($e2eResult -eq 0) {
    Write-Host "=== All Tests Passed! ===" -ForegroundColor Green
} else {
    Write-Host "=== Some Tests Failed ===" -ForegroundColor Red
}

Write-Host ""
Write-Host "Next steps for submission:" -ForegroundColor Cyan
Write-Host "1. Take screenshots of test output" -ForegroundColor White
Write-Host "2. Generate coverage reports (if unit tests passed)" -ForegroundColor White
Write-Host "3. Submit all files and screenshots" -ForegroundColor White
