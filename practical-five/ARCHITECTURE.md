# Architecture Documentation

## Service Boundaries Justification

### 1. User Service

**Why it's a separate service:**
- **Bounded Context**: User management has distinct business rules separate from orders and menu
- **Change Frequency**: User authentication and profile features evolve independently
- **Scaling Requirements**: User registration/login patterns differ from menu browsing
- **Team Ownership**: Can be owned by identity/auth team

**What it owns:**
- User entities (ID, name, email)
- User database (user_db)
- User validation logic

**What it doesn't own:**
- Order history (belongs to order service)
- Menu preferences (could be separate if needed)

### 2. Menu Service

**Why it's a separate service:**
- **Bounded Context**: Product catalog is a distinct domain
- **Read-Heavy**: Menu browsing is read-intensive and benefits from independent caching/scaling
- **Change Frequency**: Menu updates (prices, items) happen independently of user or order changes
- **Business Logic**: Menu availability, pricing rules are self-contained

**What it owns:**
- MenuItem entities (ID, name, description, price)
- Menu database (menu_db)
- Pricing logic

**What it doesn't own:**
- Order quantities (belongs to orders)
- User preferences (could be separate)

### 3. Order Service

**Why it's a separate service:**
- **Bounded Context**: Order processing is the core transaction domain
- **Complex Dependencies**: Needs to coordinate with users and menu, but doesn't own them
- **Scaling Requirements**: Peak ordering times (lunch rush) need independent scaling
- **Transaction Boundaries**: Orders are the aggregate root

**What it owns:**
- Order and OrderItem entities
- Order database (order_db)
- Order workflow logic (pending → completed)

**What it doesn't own:**
- User details (references user_id only)
- Current menu prices (but snapshots prices at order time)

## Inter-Service Communication Patterns

### Synchronous HTTP Calls
Used by Order Service to validate:
1. **User validation**: GET /users/{id}
2. **Menu validation**: GET /menu/{id}

**Why synchronous:**
- Need immediate validation before order creation
- Simple request/response pattern
- Order creation must fail if user/menu invalid

**Trade-offs:**
- **Pros**: Simple, easy to debug, strong consistency
- **Cons**: Coupling, latency, cascading failures

### Service Discovery with Consul
All services register with Consul for:
- Dynamic service location
- Health checking
- Load balancing (if scaled)

**Benefits:**
- No hardcoded URLs
- Automatic failover
- Easy to add instances

## Database-Per-Service Pattern

### Implementation
Each service has its own PostgreSQL database:
- `user_db` - User service
- `menu_db` - Menu service
- `order_db` - Order service

### Benefits
1. **Loose Coupling**: Services can't bypass APIs to access data directly
2. **Independent Scaling**: Each database can scale independently
3. **Technology Freedom**: Could use different databases per service
4. **Schema Changes**: Can change schema without coordinating with other services

### Challenges
1. **No Foreign Keys**: Can't use database foreign keys across services
2. **Distributed Transactions**: Need eventual consistency patterns
3. **Data Duplication**: Order service duplicates prices from menu

### How We Handle Challenges

#### No Foreign Keys
- **Solution**: Validate via API calls
- Order service calls user-service to validate user_id exists
- Order service calls menu-service to validate menu_item_id exists

#### Data Consistency
- **Solution**: Price snapshotting
- Orders store menu prices at order time
- Historical orders remain accurate even if menu prices change

#### Transaction Management
- **Current**: Simple fail-fast approach
- **Future**: Could implement Saga pattern or compensating transactions

## API Gateway Pattern

### Purpose
Single entry point for all client requests:
- Clients call `localhost:8080` only
- Gateway routes to appropriate service
- Hides internal service structure

### Routing Rules
```
/api/users/*  → user-service:8081
/api/menu/*   → menu-service:8082
/api/orders/* → order-service:8083
```

### Benefits
1. **Simplified Client**: Clients don't need to know about multiple services
2. **Security**: Can add authentication at gateway
3. **Cross-Cutting Concerns**: Rate limiting, logging, caching
4. **API Versioning**: Can route v1/v2 to different services

### Trade-offs
- **Single Point of Failure**: Gateway down = all services unavailable
- **Network Hop**: Adds latency
- **Bottleneck**: Can become performance bottleneck

## Comparison: Monolith vs Microservices

### Deployment

**Monolith:**
```
1 container
1 database
1 deployment process
```

**Microservices:**
```
5 containers (3 services + gateway + consul)
4 databases (3 services + monolith)
Multiple deployment processes
```

### Data Access

**Monolith:**
```go
// Direct database access
database.DB.First(&user, userID)
database.DB.First(&menuItem, menuItemID)
```

**Microservices:**
```go
// HTTP calls
http.Get("http://user-service/users/" + userID)
http.Get("http://menu-service/menu/" + menuItemID)
```

### Scaling

**Monolith:**
- Scale entire application as a unit
- All resources scaled together
- Simpler orchestration

**Microservices:**
- Scale services independently
- Menu service can have more instances than user service
- Requires orchestration (Consul, Kubernetes)

### Failure Modes

**Monolith:**
- One bug can crash entire app
- Database connection issue affects everything
- All-or-nothing availability

**Microservices:**
- Menu service down → orders fail, but users still work
- Isolated failures
- Graceful degradation possible

## When to Choose Microservices

### Choose Microservices When:
1. **Large Teams**: Multiple teams need to work independently
2. **Different Scaling**: Services have different load patterns
3. **Technology Diversity**: Need different tech stacks
4. **Fault Isolation**: Want to contain failures
5. **Frequent Deployments**: Need to deploy services independently

### Stick with Monolith When:
1. **Small Team**: 2-5 developers can manage monolith
2. **Uniform Load**: All features have similar traffic
3. **Tight Coupling**: Features are highly interconnected
4. **Starting New**: MVP/prototype stage
5. **Simple Deployment**: Want minimal operational complexity

## Future Improvements

### 1. Async Communication
- Add message queue (RabbitMQ/Kafka)
- Order creation → OrderCreated event
- Notification service subscribes to events
- Reduces coupling

### 2. Circuit Breakers
- Detect when service is down
- Fail fast instead of waiting for timeout
- Automatically retry when service recovers

### 3. Caching
- Add Redis for menu service
- Cache menu items (read-heavy)
- Reduce database load

### 4. API Versioning
- Support multiple API versions
- `/api/v1/users`, `/api/v2/users`
- Backward compatibility

### 5. Monitoring
- Add Prometheus for metrics
- Add Grafana for dashboards
- Track request rates, latencies, errors

### 6. Kubernetes Deployment
- Replace Docker Compose
- Production-grade orchestration
- Auto-scaling, rolling updates

### 7. gRPC Communication
- Replace HTTP with gRPC
- Faster, type-safe
- Better for inter-service calls

## Key Takeaways

1. **Service boundaries** are based on business capabilities, not technical layers
2. **Database-per-service** enforces loose coupling but requires careful design
3. **Inter-service communication** should be minimized and well-defined
4. **Microservices aren't free** - they trade simplicity for flexibility
5. **Start with monolith** - only split when you need to scale teams or services independently
