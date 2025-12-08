# Project Summary

## What Was Built

A complete microservices architecture demonstrating the refactoring of a monolithic Student Cafe application into independent services with service discovery, API gateway, and inter-service communication.

## Project Structure

```
practical-five/
├── student-cafe-monolith/        # Baseline monolithic application
│   ├── models/                   # User, Menu, Order models
│   ├── handlers/                 # HTTP handlers
│   ├── database/                 # Database connection
│   ├── main.go                   # Application entry point
│   ├── Dockerfile                # Container definition
│   └── docker-compose.yml        # Local deployment
│
├── user-service/                 # User management microservice
│   ├── models/user.go            # User model
│   ├── handlers/user_handlers.go # User endpoints
│   ├── database/db.go            # Database connection
│   ├── main.go                   # Service with Consul registration
│   └── Dockerfile                # Container definition
│
├── menu-service/                 # Menu catalog microservice
│   ├── models/menu.go            # MenuItem model
│   ├── handlers/menu_handlers.go # Menu endpoints
│   ├── database/db.go            # Database connection
│   ├── main.go                   # Service with Consul registration
│   └── Dockerfile                # Container definition
│
├── order-service/                # Order processing microservice
│   ├── models/order.go           # Order and OrderItem models
│   ├── handlers/order_handlers.go # Order endpoints with inter-service calls
│   ├── database/db.go            # Database connection
│   ├── main.go                   # Service with Consul registration
│   └── Dockerfile                # Container definition
│
├── api-gateway/                  # API Gateway with service discovery
│   ├── main.go                   # Gateway with Consul integration
│   └── Dockerfile                # Container definition
│
├── docker-compose.yml            # Orchestration for all services
├── README.md                     # Complete documentation
├── QUICK_START.md               # Quick start guide
├── ARCHITECTURE.md              # Architecture decisions
└── test-microservices.ps1       # Automated test script
```

## Services Overview

### 1. Monolith (Port 8090)
- **Purpose**: Baseline comparison
- **Database**: student_cafe (single database)
- **Features**: All functionality in one application
- **Use Case**: Demonstrates what we're refactoring from

### 2. User Service (Port 8081)
- **Purpose**: Manages user profiles and authentication
- **Database**: user_db (dedicated)
- **Endpoints**:
  - POST /users - Create user
  - GET /users/{id} - Get user by ID
- **Consul Integration**: ✓ Registered with health checks

### 3. Menu Service (Port 8082)
- **Purpose**: Manages food catalog and pricing
- **Database**: menu_db (dedicated)
- **Endpoints**:
  - GET /menu - List all menu items
  - POST /menu - Create menu item
  - GET /menu/{id} - Get menu item by ID
- **Consul Integration**: ✓ Registered with health checks

### 4. Order Service (Port 8083)
- **Purpose**: Processes orders with validation
- **Database**: order_db (dedicated)
- **Endpoints**:
  - POST /orders - Create order (validates with user/menu services)
  - GET /orders - List all orders
- **Inter-Service Communication**:
  - Calls user-service to validate user_id
  - Calls menu-service to get current prices
  - Uses Consul for service discovery
- **Consul Integration**: ✓ Registered with health checks

### 5. API Gateway (Port 8080)
- **Purpose**: Single entry point for all clients
- **Routes**:
  - /api/users/* → user-service
  - /api/menu/* → menu-service
  - /api/orders/* → order-service
- **Features**:
  - Dynamic service discovery via Consul
  - Reverse proxy to backend services
  - Path-based routing

### 6. Consul (Port 8500)
- **Purpose**: Service discovery and health checking
- **Features**:
  - Service registration
  - Health monitoring (10s intervals)
  - DNS and HTTP API
  - Web UI for monitoring

## Key Features Implemented

### ✅ Domain-Driven Design
- Identified bounded contexts (User, Menu, Order)
- Separated aggregates into services
- Clear service boundaries based on business capabilities

### ✅ Database-Per-Service Pattern
- Each service has dedicated database
- No cross-database queries
- Schema independence
- Prevents tight coupling

### ✅ Service Discovery
- All services register with Consul
- Health checks every 10 seconds
- Dynamic service location
- No hardcoded URLs

### ✅ Inter-Service Communication
- Order service discovers and calls user/menu services
- HTTP-based REST APIs
- Graceful error handling
- Service unavailability detection

### ✅ API Gateway
- Single client entry point
- Path-based routing
- Service discovery integration
- Simplified client code

### ✅ Price Snapshotting
- Orders store menu prices at order time
- Historical data integrity
- Price changes don't affect old orders

### ✅ Health Checks
- All services expose /health endpoint
- Consul monitors service health
- Only routes to healthy instances

## Technical Stack

- **Language**: Go 1.23
- **Web Framework**: Chi (lightweight HTTP router)
- **ORM**: GORM
- **Database**: PostgreSQL 13
- **Service Discovery**: Consul
- **Containerization**: Docker
- **Orchestration**: Docker Compose

## Getting Started

### Prerequisites
```powershell
# Required
- Docker Desktop
- Docker Compose
- Go 1.23+ (optional, for local dev)
```

### Start All Services
```powershell
cd "c:\Users\zeroe\OneDrive\Desktop\practicals Y3S1\practical-five"
docker-compose up --build
```

### Run Tests
```powershell
.\test-microservices.ps1
```

### View Service Health
```
http://localhost:8500
```

## Testing Flow

1. **Create User** → User Service → user_db
2. **Create Menu Items** → Menu Service → menu_db
3. **Create Order** → Order Service:
   - Discovers user-service via Consul
   - Calls user-service to validate user
   - Discovers menu-service via Consul
   - Calls menu-service to get prices
   - Saves order to order_db

## Learning Outcomes Demonstrated

### LO1: Architecture Comparison
- Built both monolith and microservices
- Documented trade-offs
- Identified when to use each approach

### LO2: Domain-Driven Design
- Applied bounded contexts
- Identified aggregates
- Separated based on business capabilities

### LO3: Incremental Refactoring
- Started with working monolith
- Extracted services one by one
- Maintained functionality throughout

### LO4: Service Discovery
- Implemented Consul integration
- Dynamic service location
- Health monitoring

### LO5: Orchestration
- Docker Compose for multi-container deployment
- Service dependencies
- Network configuration
- Volume management

### LO6: Migration Path
- Architecture ready for gRPC
- Can deploy to Kubernetes
- Patterns support production scale

## Challenges Solved

### 1. Service Communication
**Problem**: How do services find each other?
**Solution**: Consul service discovery

### 2. Data Consistency
**Problem**: No foreign keys across databases
**Solution**: API-based validation + price snapshotting

### 3. Service Dependencies
**Problem**: Order service needs user and menu data
**Solution**: HTTP calls with error handling

### 4. Health Monitoring
**Problem**: Knowing which services are available
**Solution**: Consul health checks

### 5. Single Entry Point
**Problem**: Clients need to know multiple service URLs
**Solution**: API Gateway with routing

## Performance Characteristics

### Monolith
- **Latency**: Low (direct database access)
- **Throughput**: Limited by single instance
- **Scalability**: All-or-nothing

### Microservices
- **Latency**: Higher (network hops)
- **Throughput**: Can scale per service
- **Scalability**: Independent scaling

### Measured with Test Script
Order creation path:
```
Client → API Gateway → Order Service
                    ↓
              Consul (discover)
                    ↓
              User Service (validate)
                    ↓
              Consul (discover)
                    ↓
              Menu Service (get prices)
                    ↓
              Order Database (save)
```

## Next Steps

### Immediate Improvements
1. Add error handling for partial failures
2. Implement retry logic
3. Add request logging/tracing
4. Add authentication at gateway

### Advanced Features
1. Replace HTTP with gRPC
2. Add circuit breakers (Hystrix pattern)
3. Implement saga pattern for transactions
4. Add caching layer (Redis)
5. Deploy to Kubernetes
6. Add monitoring (Prometheus + Grafana)

## Documentation Files

- **README.md**: Complete project documentation
- **QUICK_START.md**: Quick start guide with commands
- **ARCHITECTURE.md**: Architecture decisions and rationale
- **test-microservices.ps1**: Automated testing script

## Success Criteria Met

✅ All services run independently  
✅ Inter-service communication works  
✅ Consul service discovery implemented  
✅ API Gateway routes correctly  
✅ Database-per-service pattern implemented  
✅ Health checks functioning  
✅ Comprehensive documentation provided  
✅ Test script for validation  
✅ Docker orchestration working  
✅ Both monolith and microservices deployable  

## Conclusion

This project successfully demonstrates:
- How to refactor a monolith into microservices
- Why service boundaries matter
- How to implement service discovery
- Trade-offs between architectural styles
- Production-ready patterns and practices

The implementation is ready for:
- Academic submission
- Further enhancement (gRPC, Kubernetes)
- Production hardening (monitoring, security)
- Team-based development
