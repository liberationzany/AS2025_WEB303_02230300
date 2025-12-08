# Troubleshooting Guide

## Common Issues and Solutions

### 1. Services Won't Start

#### Symptom
```
Error response from daemon: Ports are not available
```

#### Solution
Another application is using the required ports. Check which ports are in use:

```powershell
# Check what's using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

Or modify ports in `docker-compose.yml`:
```yaml
api-gateway:
  ports:
    - "9080:8080"  # Changed from 8080 to 9080
```

---

### 2. Go Module Checksum Errors

#### Symptom
```
SECURITY ERROR: checksum mismatch for github.com/go-chi/chi/v5
```

#### Solution
Regenerate go.sum files:

```powershell
# For each service
cd user-service
Remove-Item go.sum
go mod tidy

cd ..\menu-service
Remove-Item go.sum
go mod tidy

cd ..\order-service
Remove-Item go.sum
go mod tidy

cd ..\api-gateway
Remove-Item go.sum
go mod tidy

cd ..\student-cafe-monolith
Remove-Item go.sum
go mod tidy
```

---

### 3. Services Not Appearing in Consul

#### Symptom
Consul UI shows no services or services are failing health checks.

#### Solution

**Check if Consul is running:**
```powershell
docker-compose ps consul
```

**Wait for initialization:**
Services need 10-15 seconds to:
1. Start container
2. Connect to database
3. Register with Consul

**Check service logs:**
```powershell
docker-compose logs user-service
```

Look for:
```
User database connected
User service starting on :8081
```

**Verify health endpoint:**
```powershell
curl http://localhost:8081/health
```

Should return: 200 OK

---

### 4. Order Creation Fails

#### Symptom
```json
{"error": "User not found"}
```
or
```json
{"error": "Menu item not found"}
```

#### Solution

**Ensure user exists:**
```powershell
# Create user first
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{\"name\": \"Test User\", \"email\": \"test@example.com\"}'

# Verify it was created
curl http://localhost:8080/api/users/1
```

**Ensure menu items exist:**
```powershell
# Create menu item
curl -X POST http://localhost:8080/api/menu -H "Content-Type: application/json" -d '{\"name\": \"Coffee\", \"description\": \"Hot coffee\", \"price\": 2.50}'

# Verify it was created
curl http://localhost:8080/api/menu
```

**Check service availability:**
```powershell
# Check Consul UI
Start-Process http://localhost:8500

# All services should be green (healthy)
```

---

### 5. Service Unavailable Errors

#### Symptom
```json
{"error": "User service unavailable"}
```
or
```json
{"error": "Menu service unavailable"}
```

#### Solution

**Check service is running:**
```powershell
docker-compose ps
```

All services should show "Up" status.

**Restart the service:**
```powershell
docker-compose restart user-service
```

**Check service logs:**
```powershell
docker-compose logs -f user-service
```

**Verify Consul registration:**
1. Open http://localhost:8500
2. Check if service appears
3. Verify health check is passing (green)

---

### 6. Database Connection Errors

#### Symptom
```
Failed to connect to database: connection refused
```

#### Solution

**Check database is running:**
```powershell
docker-compose ps user-db
```

**Restart database:**
```powershell
docker-compose restart user-db
```

**Check logs:**
```powershell
docker-compose logs user-db
```

Look for:
```
database system is ready to accept connections
```

**Verify connection string:**
Check environment variables in `docker-compose.yml`:
```yaml
environment:
  DATABASE_URL: "host=user-db user=postgres password=postgres dbname=user_db port=5432 sslmode=disable"
```

---

### 7. Docker Build Fails

#### Symptom
```
failed to solve: failed to copy files
```

#### Solution

**Ensure go.sum exists:**
```powershell
cd user-service
go mod tidy
```

**Clean Docker cache:**
```powershell
docker-compose down
docker system prune -f
docker-compose up --build
```

**Check Dockerfile syntax:**
Ensure no typos in Dockerfile.

---

### 8. API Gateway Returns 404

#### Symptom
```
404 page not found
```

#### Solution

**Check URL format:**
All requests must go through `/api` prefix:
```powershell
# Correct
curl http://localhost:8080/api/menu

# Incorrect
curl http://localhost:8080/menu
```

**Verify gateway is running:**
```powershell
docker-compose ps api-gateway
```

**Check gateway logs:**
```powershell
docker-compose logs api-gateway
```

---

### 9. Services Can't Reach Each Other

#### Symptom
Order service logs show:
```
Get "http://user-service:8081/users/1": dial tcp: lookup user-service: no such host
```

#### Solution

**Ensure services are on same network:**
Docker Compose creates a default network automatically.

**Restart all services:**
```powershell
docker-compose down
docker-compose up
```

**Check network:**
```powershell
docker network ls
docker network inspect practical-five_default
```

---

### 10. Curl Commands Don't Work

#### Symptom
```
curl: command not found
```
or
```
Invoke-WebRequest : A parameter cannot be found that matches parameter name 'd'.
```

#### Solution

**For PowerShell, use backticks for line continuation:**
```powershell
curl -X POST http://localhost:8080/api/users `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"John\", \"email\": \"john@example.com\"}'
```

**Or use Invoke-RestMethod:**
```powershell
$body = @{
    name = "John Doe"
    email = "john@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Method POST -Uri http://localhost:8080/api/users -Body $body -ContentType "application/json"
```

**Or use the test script:**
```powershell
.\test-microservices.ps1
```

---

### 11. Slow Performance

#### Symptom
Requests take 5+ seconds.

#### Solution

**Check if services are healthy:**
```powershell
# Open Consul UI
Start-Process http://localhost:8500
```

**Restart services:**
```powershell
docker-compose restart
```

**Check system resources:**
```powershell
docker stats
```

Ensure CPU and memory are not maxed out.

**Rebuild without cache:**
```powershell
docker-compose down
docker-compose build --no-cache
docker-compose up
```

---

### 12. Data Not Persisting

#### Symptom
After restarting, all users/menu items/orders are gone.

#### Solution

**Check volumes:**
```powershell
docker volume ls
```

Should see:
- practical-five_postgres_data
- practical-five_user_data
- practical-five_menu_data
- practical-five_order_data

**Don't use `-v` flag when stopping:**
```powershell
# This removes volumes (data loss)
docker-compose down -v

# This keeps volumes (data persists)
docker-compose down
```

---

### 13. Can't Access Consul UI

#### Symptom
Browser shows "Can't reach this page" for http://localhost:8500

#### Solution

**Check Consul is running:**
```powershell
docker-compose ps consul
```

**Check port mapping:**
```powershell
docker port $(docker-compose ps -q consul)
```

Should show:
```
8500/tcp -> 0.0.0.0:8500
```

**Restart Consul:**
```powershell
docker-compose restart consul
```

---

### 14. Test Script Fails

#### Symptom
```powershell
.\test-microservices.ps1
# Error: Cannot retrieve menu
```

#### Solution

**Ensure all services are running:**
```powershell
docker-compose ps
```

**Wait for services to be healthy:**
Services need 30-60 seconds to fully initialize.

**Check Consul UI:**
All services should be green before running tests.

**Run services first:**
```powershell
docker-compose up -d
Start-Sleep -Seconds 30
.\test-microservices.ps1
```

---

## Diagnostic Commands

### Check All Services Status
```powershell
docker-compose ps
```

### View All Logs
```powershell
docker-compose logs
```

### View Specific Service Logs
```powershell
docker-compose logs -f user-service
```

### Test Service Connectivity
```powershell
# Test user service
curl http://localhost:8081/health

# Test menu service
curl http://localhost:8082/health

# Test order service
curl http://localhost:8083/health

# Test API gateway
curl http://localhost:8080/api/menu
```

### Check Network
```powershell
docker network ls
docker network inspect practical-five_default
```

### Check Volumes
```powershell
docker volume ls
docker volume inspect practical-five_user_data
```

### Full Reset
```powershell
# Stop everything
docker-compose down -v

# Remove all containers
docker-compose rm -f

# Rebuild and start
docker-compose up --build
```

---

## Getting Help

If none of these solutions work:

1. **Check logs carefully:**
   ```powershell
   docker-compose logs > debug.log
   ```

2. **Verify Go installation:**
   ```powershell
   go version
   ```

3. **Verify Docker installation:**
   ```powershell
   docker --version
   docker-compose --version
   ```

4. **Check system resources:**
   - Docker Desktop needs 4GB+ RAM
   - 10GB+ disk space

5. **Review documentation:**
   - README.md
   - QUICK_START.md
   - ARCHITECTURE.md
