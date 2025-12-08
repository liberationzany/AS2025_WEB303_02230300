# Quick Test Script - No MinGW Required
# This runs E2E tests only (works with Docker)

Write-Host "`n=== Practical 6 - Quick E2E Test ===" -ForegroundColor Cyan
Write-Host "This script runs E2E tests via Docker (no MinGW/GCC needed)`n" -ForegroundColor Yellow

# Check Docker
Write-Host "Checking Docker..." -ForegroundColor White
try {
    $dockerVersion = docker --version
    Write-Host "✓ Docker found: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker not found! Please install Docker Desktop." -ForegroundColor Red
    exit 1
}

Write-Host "`n=== Step 1: Building Docker Images ===" -ForegroundColor Green
Write-Host "This may take a few minutes..." -ForegroundColor Yellow
docker compose build

if ($LASTEXITCODE -ne 0) {
    Write-Host "`n✗ Docker build failed!" -ForegroundColor Red
    Write-Host "Check the error messages above." -ForegroundColor Yellow
    exit 1
}

Write-Host "`n=== Step 2: Starting Services ===" -ForegroundColor Green
docker compose up -d

Write-Host "`nWaiting for services to initialize (20 seconds)..." -ForegroundColor Yellow
Start-Sleep -Seconds 20

Write-Host "`n=== Step 3: Checking Service Status ===" -ForegroundColor Green
docker compose ps

Write-Host "`n=== Step 4: Running E2E Tests ===" -ForegroundColor Green
Set-Location tests\e2e
go test -v ./...
$testResult = $LASTEXITCODE
Set-Location ..\..

Write-Host "`n=== Step 5: Viewing Service Logs ===" -ForegroundColor Green
Write-Host "(Last 20 lines from each service)`n" -ForegroundColor Yellow
docker compose logs --tail=20

Write-Host "`n=== Step 6: Stopping Services ===" -ForegroundColor Green
docker compose down

Write-Host "`n=== Test Results ===" -ForegroundColor Cyan
if ($testResult -eq 0) {
    Write-Host "✓ E2E Tests PASSED!" -ForegroundColor Green
    Write-Host "`nYou can now take screenshots of the test output above." -ForegroundColor White
} else {
    Write-Host "✗ Some E2E tests failed." -ForegroundColor Red
    Write-Host "Review the error messages above." -ForegroundColor Yellow
}

Write-Host "`n=== Next Steps ===" -ForegroundColor Cyan
Write-Host "For FULL testing (including unit & integration tests):" -ForegroundColor White
Write-Host "1. Install MinGW: choco install mingw" -ForegroundColor Yellow
Write-Host "2. Run: .\run-tests.ps1" -ForegroundColor Yellow
Write-Host "`nFor submission:" -ForegroundColor White
Write-Host "- Screenshot the E2E test output above" -ForegroundColor Yellow
Write-Host "- Include all source code files" -ForegroundColor Yellow
Write-Host "- Submit via your assignment portal" -ForegroundColor Yellow
