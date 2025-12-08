# Build and Run Checklist

## Pre-Flight Checklist

### System Requirements
- [ ] Docker Desktop installed and running
- [ ] Docker Compose available (comes with Docker Desktop)
- [ ] Go 1.23+ installed (optional, for local development)
- [ ] At least 4GB RAM available for Docker
- [ ] At least 10GB disk space available
- [ ] PowerShell available (comes with Windows)

### Verify Installation
```powershell
# Check Docker
docker --version
# Expected: Docker version 20.10+ or higher

# Check Docker Compose
docker-compose --version
# Expected: docker-compose version 1.29+ or higher

# Check Go (optional)
go version
# Expected: go version go1.23+ or higher

# Check Docker is running
docker ps
# Should list containers (may be empty)
```

---

## Build Process

### Step 1: Navigate to Project Directory
```powershell
cd "c:\Users\zeroe\OneDrive\Desktop\practicals Y3S1\practical-five"
```

- [ ] Confirmed in correct directory

### Step 2: Verify Project Structure
```powershell
Get-ChildItem
```

Should see:
- [ ] api-gateway/
- [ ] menu-service/
- [ ] order-service/
- [ ] student-cafe-monolith/
- [ ] user-service/
- [ ] docker-compose.yml
- [ ] README.md
- [ ] test-microservices.ps1

### Step 3: Initialize Go Modules (Optional)
```powershell
# Only if you modified code
cd user-service; go mod tidy; cd ..
cd menu-service; go mod tidy; cd ..
cd order-service; go mod tidy; cd ..
cd api-gateway; go mod tidy; cd ..
cd student-cafe-monolith; go mod tidy; cd ..
```

- [ ] All go.sum files generated

### Step 4: Build and Start Services
```powershell
docker-compose up --build
```

This will:
1. Build Docker images for all services (5-10 minutes first time)
2. Create networks
3. Create volumes
4. Start all containers

- [ ] Build started without errors

### Step 5: Wait for Services to Initialize
Watch the logs for these messages:

```
consul_1           | ==> Consul agent running!
postgres_1         | database system is ready to accept connections
user-db_1          | database system is ready to accept connections
menu-db_1          | database system is ready to accept connections
order-db_1         | database system is ready to accept connections
user-service_1     | User database connected
user-service_1     | User service starting on :8081
menu-service_1     | Menu database connected
menu-service_1     | Menu service starting on :8082
order-service_1    | Order database connected
order-service_1    | Order service starting on :8083
api-gateway_1      | API Gateway starting on :8080
monolith_1         | Monolith server starting on :8080
```

- [ ] All databases started
- [ ] All services started
- [ ] No error messages in logs

**Expected wait time:** 30-60 seconds

---

## Verification Checklist

### Step 6: Verify Consul
Open in browser:
```
http://localhost:8500
```

- [ ] Consul UI loads
- [ ] "Services" page shows services
- [ ] user-service is listed
- [ ] menu-service is listed
- [ ] order-service is listed
- [ ] All services show green (healthy)
- [ ] Health checks are passing

### Step 7: Verify API Gateway
```powershell
curl http://localhost:8080/api/menu
```

Expected response:
```json
[]
```
(Empty array is correct - no menu items yet)

- [ ] Gateway responds
- [ ] No error message

### Step 8: Verify Individual Services
```powershell
# Test user service
curl http://localhost:8081/health

# Test menu service
curl http://localhost:8082/health

# Test order service
curl http://localhost:8083/health
```

- [ ] All health endpoints return 200 OK

### Step 9: Verify Databases
```powershell
docker-compose ps
```

Should show:
- [ ] postgres (port 5432) - Up
- [ ] user-db (port 5434) - Up
- [ ] menu-db (port 5433) - Up
- [ ] order-db (port 5435) - Up

---

## Functional Testing Checklist

### Step 10: Create Test Data
```powershell
# Create a user
curl -X POST http://localhost:8080/api/users `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"Test User\", \"email\": \"test@example.com\"}'
```

- [ ] User created successfully
- [ ] Received JSON response with user ID

```powershell
# Create menu items
curl -X POST http://localhost:8080/api/menu `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"Coffee\", \"description\": \"Hot coffee\", \"price\": 2.50}'
```

- [ ] Menu item created successfully

### Step 11: Test Order Creation (Inter-Service Communication)
```powershell
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 2}]}'
```

- [ ] Order created successfully
- [ ] Received JSON response with order details

### Step 12: Verify Inter-Service Communication
In a separate terminal:
```powershell
docker-compose logs -f order-service
```

Look for:
- [ ] "Proxying to user-service"
- [ ] "Proxying to menu-service"
- [ ] Successful API calls

### Step 13: Run Automated Tests
```powershell
.\test-microservices.ps1
```

- [ ] All tests pass
- [ ] Users created
- [ ] Menu items created
- [ ] Orders created
- [ ] No errors in script output

---

## Service-Specific Checks

### User Service
- [ ] Can create user
- [ ] Can retrieve user by ID
- [ ] Database persists data

### Menu Service
- [ ] Can create menu item
- [ ] Can list all menu items
- [ ] Can retrieve menu item by ID
- [ ] Database persists data

### Order Service
- [ ] Can create order with valid user and menu items
- [ ] Fails with invalid user ID (returns 400)
- [ ] Fails with invalid menu item ID (returns 400)
- [ ] Successfully calls user-service for validation
- [ ] Successfully calls menu-service for prices
- [ ] Stores snapshot prices in order

### API Gateway
- [ ] Routes /api/users/* to user-service
- [ ] Routes /api/menu/* to menu-service
- [ ] Routes /api/orders/* to order-service
- [ ] Returns 503 when service is unavailable

### Consul
- [ ] All services registered
- [ ] Health checks passing
- [ ] Services can discover each other
- [ ] UI accessible

---

## Resilience Testing Checklist

### Step 14: Test Service Failure
```powershell
# Stop menu service
docker-compose stop menu-service
```

- [ ] Consul shows menu-service as unhealthy

```powershell
# Try to create order
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}]}'
```

- [ ] Order creation fails gracefully
- [ ] Error message indicates service unavailable

```powershell
# Restart menu service
docker-compose start menu-service
```

- [ ] Service restarts successfully
- [ ] Consul shows menu-service as healthy
- [ ] Order creation works again

---

## Performance Checklist

### Step 15: Monitor Resource Usage
```powershell
docker stats
```

Check:
- [ ] CPU usage reasonable (<50% per container)
- [ ] Memory usage reasonable (<200MB per service)
- [ ] No container constantly restarting

### Step 16: Test Response Times
```powershell
Measure-Command { curl http://localhost:8080/api/menu }
```

- [ ] Response time < 1 second for simple requests
- [ ] Response time < 3 seconds for order creation

---

## Documentation Checklist

### Project Documentation
- [ ] README.md exists and is complete
- [ ] QUICK_START.md provides clear instructions
- [ ] ARCHITECTURE.md explains design decisions
- [ ] TROUBLESHOOTING.md covers common issues
- [ ] PROJECT_SUMMARY.md provides overview

### Code Documentation
- [ ] Each service has clear structure
- [ ] Dockerfiles are properly configured
- [ ] docker-compose.yml is well-organized
- [ ] Go code is readable and commented

---

## Cleanup Checklist

### Stop Services
```powershell
# Stop but keep data
docker-compose down

# Stop and remove data
docker-compose down -v
```

- [ ] Services stopped cleanly

### Verify Cleanup
```powershell
docker-compose ps
```

- [ ] No containers running

---

## Common Issues Quick Reference

| Issue | Quick Fix |
|-------|-----------|
| Port conflict | Change ports in docker-compose.yml |
| Service won't start | Check logs: `docker-compose logs <service>` |
| Not in Consul | Wait 15 seconds, check health endpoint |
| Order fails | Ensure user and menu items exist |
| Slow build | First build takes 5-10 minutes (normal) |
| Can't curl | Use PowerShell backticks or test script |

---

## Success Criteria

✅ All services running (docker-compose ps shows "Up")  
✅ All services healthy in Consul (green status)  
✅ Can create users, menu items, and orders  
✅ Order service successfully calls other services  
✅ Test script completes without errors  
✅ No error messages in logs  
✅ API Gateway routes correctly  
✅ Services can restart and recover  

---

## Next Steps After Successful Build

1. [ ] Explore Consul UI to understand service discovery
2. [ ] Review service logs to see inter-service calls
3. [ ] Test failure scenarios (stop services)
4. [ ] Review code to understand implementation
5. [ ] Read ARCHITECTURE.md to understand design decisions
6. [ ] Consider enhancements (monitoring, caching, etc.)

---

## Support

If any checklist item fails:
1. Check TROUBLESHOOTING.md
2. Review service logs: `docker-compose logs`
3. Verify system requirements met
4. Try full rebuild: `docker-compose down -v && docker-compose up --build`

---

## Final Verification Command

Run this to verify everything is working:

```powershell
# Quick health check
curl http://localhost:8081/health; `
curl http://localhost:8082/health; `
curl http://localhost:8083/health; `
curl http://localhost:8080/api/menu

# Or use the test script
.\test-microservices.ps1
```

If all commands succeed, your microservices architecture is fully operational! 🎉
