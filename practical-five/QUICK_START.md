# Quick Start Guide

## Build and Start All Services

```powershell
cd "c:\Users\zeroe\OneDrive\Desktop\practicals Y3S1\practical-five"
docker-compose up --build
```

Wait for all services to start (30-60 seconds).

## Verify Services are Running

Open Consul UI in browser:
```
http://localhost:8500
```

You should see all services (user-service, menu-service, order-service) in green (healthy).

## Test the Application

### Step 1: Create a User
```powershell
curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json" -d '{\"name\": \"John Doe\", \"email\": \"john@example.com\"}'
```

Expected response:
```json
{"ID":1,"CreatedAt":"...","UpdatedAt":"...","DeletedAt":null,"name":"John Doe","email":"john@example.com"}
```

### Step 2: Create Menu Items
```powershell
curl -X POST http://localhost:8080/api/menu -H "Content-Type: application/json" -d '{\"name\": \"Coffee\", \"description\": \"Hot coffee\", \"price\": 2.50}'

curl -X POST http://localhost:8080/api/menu -H "Content-Type: application/json" -d '{\"name\": \"Sandwich\", \"description\": \"Fresh sandwich\", \"price\": 5.00}'

curl -X POST http://localhost:8080/api/menu -H "Content-Type: application/json" -d '{\"name\": \"Tea\", \"description\": \"Hot tea\", \"price\": 1.50}'
```

### Step 3: View Menu
```powershell
curl http://localhost:8080/api/menu
```

### Step 4: Create an Order
```powershell
curl -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 2}, {\"menu_item_id\": 2, \"quantity\": 1}]}'
```

This order includes:
- 2x Coffee
- 1x Sandwich

### Step 5: View All Orders
```powershell
curl http://localhost:8080/api/orders
```

## Monitor Inter-Service Communication

In a separate terminal, watch the order service logs:
```powershell
docker-compose logs -f order-service
```

When you create an order, you'll see:
1. Service discovery calls to Consul
2. HTTP calls to user-service
3. HTTP calls to menu-service
4. Database operations

## Test Service Resilience

### Stop Menu Service
```powershell
docker-compose stop menu-service
```

### Try to Create Order (Will Fail)
```powershell
curl -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}]}'
```

Should return: "Menu service unavailable"

### Restart Menu Service
```powershell
docker-compose start menu-service
```

Wait 10 seconds for service to register with Consul, then retry the order creation.

## Compare with Monolith

The monolith is running on port 8090. You can test it the same way:

```powershell
curl http://localhost:8090/api/menu
```

## Stop All Services

```powershell
docker-compose down
```

## Clean Up (Remove All Data)

```powershell
docker-compose down -v
```

## Troubleshooting

### Services Not Appearing in Consul
- Wait 15 seconds after startup
- Check service logs: `docker-compose logs user-service`
- Ensure Consul is running: `docker-compose ps consul`

### Order Creation Fails
- Ensure user exists (create user first)
- Ensure menu items exist (create menu items first)
- Check all services are healthy in Consul

### Port Already in Use
- Stop other applications using ports 8080-8083, 8500, 5432-5435
- Or modify port mappings in docker-compose.yml
