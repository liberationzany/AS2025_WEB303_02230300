# Student Cafe Microservices Project

This project demonstrates the systematic refactoring of a monolithic application into microservices, following domain-driven design principles.

## Project Structure

```
practical-five/
├── student-cafe-monolith/    # Original monolithic application
├── user-service/             # User management microservice
├── menu-service/             # Menu catalog microservice
├── order-service/            # Order processing microservice
├── api-gateway/              # API Gateway with service discovery
├── docker-compose.yml        # Orchestration for all services
└── README.md                 # This file
```

## Architecture Overview

### Monolithic Architecture
- Single application with all features
- Shared database
- Tight coupling between components
- Port: 8090

### Microservices Architecture
- **User Service** (Port 8081): Manages user profiles and authentication
- **Menu Service** (Port 8082): Handles food catalog and pricing
- **Order Service** (Port 8083): Processes orders with inter-service communication
- **API Gateway** (Port 8080): Single entry point for all client requests
- **Consul** (Port 8500): Service discovery and health checking

Each service has its own database following the database-per-service pattern.

## Service Boundaries

### User Service
- **Bounded Context**: User Management
- **Entities**: User
- **Responsibilities**: User registration, profile management
- **Database**: user_db
- **Scales When**: High user registration traffic

### Menu Service
- **Bounded Context**: Menu Management
- **Entities**: MenuItem
- **Responsibilities**: Menu catalog, pricing updates
- **Database**: menu_db
- **Scales When**: High browsing traffic

### Order Service
- **Bounded Context**: Order Management
- **Entities**: Order, OrderItem
- **Responsibilities**: Order creation, order history
- **Database**: order_db
- **Scales When**: Peak ordering periods (lunch rush)
- **Dependencies**: Calls user-service and menu-service via HTTP

## Prerequisites

- Docker and Docker Compose installed
- Go 1.23+ (for local development)
- PowerShell (for Windows)

## Getting Started

### 1. Start All Services

```powershell
docker-compose up --build
```

This will start:
- 4 PostgreSQL databases
- Consul for service discovery
- 3 microservices (user, menu, order)
- API Gateway
- Monolith (for comparison)

### 2. Verify Services

Wait for all services to start (about 30-60 seconds). Check Consul UI:

```
http://localhost:8500
```

All services should show as healthy (green).

### 3. Test the Microservices

#### Create a User
```powershell
curl -X POST http://localhost:8080/api/users `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"John Doe\", \"email\": \"john@example.com\"}'
```

#### Create Menu Items
```powershell
curl -X POST http://localhost:8080/api/menu `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"Coffee\", \"description\": \"Hot coffee\", \"price\": 2.50}'

curl -X POST http://localhost:8080/api/menu `
  -H "Content-Type: application/json" `
  -d '{\"name\": \"Sandwich\", \"description\": \"Fresh sandwich\", \"price\": 5.00}'
```

#### Get Menu
```powershell
curl http://localhost:8080/api/menu
```

#### Create an Order
```powershell
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 2}, {\"menu_item_id\": 2, \"quantity\": 1}]}'
```

#### Get All Orders
```powershell
curl http://localhost:8080/api/orders
```

### 4. Test Inter-Service Communication

Watch the logs to see inter-service communication:

```powershell
docker-compose logs -f order-service
```

When you create an order, you'll see:
1. Order service discovers user-service via Consul
2. Order service calls user-service to validate user
3. Order service discovers menu-service via Consul
4. Order service calls menu-service to get prices
5. Order service saves the order with snapshot prices

## Key Features Demonstrated

### 1. Database-Per-Service Pattern
Each microservice has its own database:
- Ensures loose coupling
- Allows independent scaling
- Prevents direct database access between services

### 2. Service Discovery with Consul
- Services register themselves on startup
- Health checks ensure only healthy instances are used
- API Gateway and Order Service discover services dynamically
- No hardcoded URLs

### 3. Inter-Service Communication
Order service demonstrates:
- HTTP-based service-to-service calls
- Dynamic service discovery
- Error handling when services are unavailable

### 4. API Gateway Pattern
- Single entry point for clients
- Routes requests to appropriate services
- Hides internal service structure
- Can add authentication, rate limiting, etc.

### 5. Price Snapshotting
Orders store menu prices at order time:
- Historical orders aren't affected by price changes
- Demonstrates temporal data handling
- Common pattern in e-commerce systems

## Comparing Monolith vs Microservices

### Monolith (Port 8090)
```powershell
# Test monolith
curl http://localhost:8090/api/menu
```

**Advantages:**
- Simpler deployment (one container)
- No network latency between components
- Easier to test locally
- Simpler transaction management

**Disadvantages:**
- Scales as a single unit
- Tight coupling
- One failure can bring down entire app
- Difficult to adopt new technologies

### Microservices (Port 8080)
```powershell
# Test via API Gateway
curl http://localhost:8080/api/menu
```

**Advantages:**
- Independent scaling
- Technology flexibility
- Fault isolation
- Team autonomy

**Disadvantages:**
- Complex deployment
- Network overhead
- Distributed transaction challenges
- More moving parts

## Troubleshooting

### Services Not Starting
```powershell
# Check logs
docker-compose logs

# Restart specific service
docker-compose restart user-service
```

### Consul Not Showing Services
- Wait 10-15 seconds after startup
- Services register after connecting to database
- Check service logs: `docker-compose logs user-service`

### Port Conflicts
If ports are already in use:
1. Stop other services using those ports
2. Or modify port mappings in docker-compose.yml

### Order Creation Fails
Common causes:
- User ID doesn't exist (create user first)
- Menu item ID doesn't exist (create menu items first)
- User or menu service is down (check Consul)

## Development

### Run Individual Service Locally
```powershell
cd user-service
go mod tidy
go run main.go
```

### Rebuild Single Service
```powershell
docker-compose up --build user-service
```

### View Service Logs
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f order-service
```

### Stop All Services
```powershell
docker-compose down
```

### Clean Up (Remove Volumes)
```powershell
docker-compose down -v
```

## Architecture Diagram

```
                    ┌─────────────┐
                    │   Clients   │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ API Gateway │
                    │  (Port 8080)│
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────▼──────┐  ┌─────▼──────┐  ┌─────▼──────┐
    │   User     │  │    Menu    │  │   Order    │
    │  Service   │  │  Service   │  │  Service   │
    │ (Port 8081)│  │ (Port 8082)│  │ (Port 8083)│
    └─────┬──────┘  └─────┬──────┘  └─────┬──────┘
          │                │                │
    ┌─────▼──────┐  ┌─────▼──────┐  ┌─────▼──────┐
    │  user_db   │  │  menu_db   │  │  order_db  │
    └────────────┘  └────────────┘  └────────────┘

                    ┌─────────────┐
                    │   Consul    │
                    │ (Port 8500) │
                    │   Service   │
                    │  Discovery  │
                    └─────────────┘
```

## Testing Scenarios

### 1. Service Failure Resilience
```powershell
# Stop menu service
docker-compose stop menu-service

# Try to create order - should fail gracefully
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}]}'

# Restart menu service
docker-compose start menu-service

# Wait for registration, then retry
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}]}'
```

### 2. Service Discovery
```powershell
# Scale order service
docker-compose up --scale order-service=2

# Check Consul - both instances should appear
# Open http://localhost:8500
```

### 3. Price Changes Don't Affect Old Orders
```powershell
# Create an order
curl -X POST http://localhost:8080/api/orders `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": 1, \"items\": [{\"menu_item_id\": 1, \"quantity\": 1}]}'

# Update menu item price (would need update endpoint)
# Old order still has original price
curl http://localhost:8080/api/orders
```

## Next Steps

1. **Add gRPC**: Replace HTTP with gRPC for faster inter-service communication
2. **Kubernetes**: Deploy to Kubernetes for production-grade orchestration
3. **Monitoring**: Add Prometheus and Grafana for metrics
4. **Circuit Breakers**: Implement resilience patterns
5. **Authentication**: Add JWT-based auth at API Gateway
6. **Message Queue**: Use RabbitMQ/Kafka for async communication

## Learning Outcomes Achieved

✅ **LO1**: Identified characteristics and trade-offs of monolith vs microservices  
✅ **LO2**: Applied domain-driven design to identify service boundaries  
✅ **LO3**: Systematically extracted services while maintaining functionality  
✅ **LO4**: Implemented Consul service discovery  
✅ **LO5**: Deployed and orchestrated services with Docker Compose  
✅ **LO6**: Understand path to gRPC and Kubernetes  

## References

- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Microservices Patterns](https://microservices.io/)
- [Consul Service Discovery](https://www.consul.io/)
- [Database per Service Pattern](https://microservices.io/patterns/data/database-per-service.html)
