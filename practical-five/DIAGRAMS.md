# Architecture Diagrams

## Overall System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                             │
│                    (Web Browser, Mobile App)                     │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │ HTTP/REST
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                        API GATEWAY                               │
│                      (Port 8080)                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ • Route: /api/users/*  → user-service                    │  │
│  │ • Route: /api/menu/*   → menu-service                    │  │
│  │ • Route: /api/orders/* → order-service                   │  │
│  │ • Service Discovery via Consul                           │  │
│  │ • Health Check Aware Routing                             │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────┬─────────────┬─────────────┬────────────────────┘
                 │             │             │
        ┌────────▼──┐    ┌────▼─────┐   ┌──▼────────┐
        │  Discover │    │ Discover │   │  Discover │
        │  Service  │    │ Service  │   │  Service  │
        └────┬──────┘    └────┬─────┘   └──┬────────┘
             │                │             │
             │           ┌────▼─────────────▼───┐
             │           │   CONSUL (8500)      │
             │           │ Service Discovery     │
             │           │ Health Monitoring     │
             │           └────┬─────────────┬───┘
             │                │             │
             │         ┌──────▼──┐   ┌─────▼──────┐   ┌─────────┐
             └─────────► Register│   │  Register  │   │Register │
                       └────┬────┘   └─────┬──────┘   └────┬────┘
                            │              │               │
         ┌──────────────────▼───┐  ┌──────▼──────┐  ┌─────▼──────────┐
         │   USER SERVICE       │  │MENU SERVICE │  │ ORDER SERVICE  │
         │   (Port 8081)        │  │(Port 8082)  │  │ (Port 8083)    │
         │                      │  │             │  │                │
         │ Endpoints:           │  │ Endpoints:  │  │ Endpoints:     │
         │ • POST /users        │  │ • GET /menu │  │ • POST /orders │
         │ • GET /users/{id}    │  │ • POST /menu│  │ • GET /orders  │
         │ • GET /health        │  │ • GET /menu│  │ • GET /health  │
         │                      │  │   /{id}     │  │                │
         │                      │  │ • GET /health│ │ Inter-Service: │
         │                      │  │             │  │ • Call User    │
         │                      │  │             │  │ • Call Menu    │
         └──────────┬───────────┘  └──────┬──────┘  └─────┬──────────┘
                    │                     │               │
                    │                     │               │
         ┌──────────▼───────────┐  ┌─────▼───────┐  ┌────▼───────────┐
         │     user_db          │  │   menu_db   │  │   order_db     │
         │  (PostgreSQL)        │  │(PostgreSQL) │  │ (PostgreSQL)   │
         │  Port: 5434          │  │ Port: 5433  │  │  Port: 5435    │
         │                      │  │             │  │                │
         │ Tables:              │  │ Tables:     │  │ Tables:        │
         │ • users              │  │ • menu_items│  │ • orders       │
         │                      │  │             │  │ • order_items  │
         └──────────────────────┘  └─────────────┘  └────────────────┘
```

---

## Request Flow Diagram

### Simple Request (Menu List)

```
┌────────┐
│ Client │
└───┬────┘
    │
    │ 1. GET /api/menu
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 2. Discover menu-service
    │
┌───▼────────┐
│   Consul   │
└───┬────────┘
    │
    │ 3. Return: http://menu-service:8082
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 4. Forward: GET /menu
    │
┌───▼──────────┐
│Menu Service  │
└───┬──────────┘
    │
    │ 5. Query database
    │
┌───▼──────────┐
│   menu_db    │
└───┬──────────┘
    │
    │ 6. Return menu items
    │
┌───▼──────────┐
│Menu Service  │
└───┬──────────┘
    │
    │ 7. Return JSON
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 8. Return to client
    │
┌───▼────┐
│ Client │
└────────┘
```

---

## Complex Request Flow (Create Order)

```
┌────────┐
│ Client │
└───┬────┘
    │
    │ 1. POST /api/orders
    │    {user_id: 1, items: [{menu_item_id: 1, quantity: 2}]}
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 2. Discover order-service
    │
┌───▼────────┐
│   Consul   │─────────────────────────────┐
└───┬────────┘                             │
    │                                      │
    │ 3. Return: http://order-service:8083│
    │                                      │
┌───▼──────────┐                           │
│ API Gateway  │                           │
└───┬──────────┘                           │
    │                                      │
    │ 4. Forward: POST /orders             │
    │                                      │
┌───▼──────────┐                           │
│Order Service │                           │
└───┬──────────┘                           │
    │                                      │
    │ 5. Discover user-service ────────────┘
    │
┌───▼────────┐
│   Consul   │
└───┬────────┘
    │
    │ 6. Return: http://user-service:8081
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 7. GET /users/1
    │
┌───▼──────────┐
│User Service  │
└───┬──────────┘
    │
    │ 8. Query database
    │
┌───▼──────────┐
│   user_db    │
└───┬──────────┘
    │
    │ 9. Return user data
    │
┌───▼──────────┐
│User Service  │
└───┬──────────┘
    │
    │ 10. Return JSON {id: 1, name: "John"}
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 11. Discover menu-service
    │
┌───▼────────┐
│   Consul   │
└───┬────────┘
    │
    │ 12. Return: http://menu-service:8082
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 13. GET /menu/1
    │
┌───▼──────────┐
│Menu Service  │
└───┬──────────┘
    │
    │ 14. Query database
    │
┌───▼──────────┐
│   menu_db    │
└───┬──────────┘
    │
    │ 15. Return menu item
    │
┌───▼──────────┐
│Menu Service  │
└───┬──────────┘
    │
    │ 16. Return JSON {id: 1, price: 2.50}
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 17. Create order with snapshot price
    │
┌───▼──────────┐
│   order_db   │
└───┬──────────┘
    │
    │ 18. Order saved
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 19. Return order JSON
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 20. Return to client
    │
┌───▼────┐
│ Client │
└────────┘
```

**Total: 20 steps vs 3-4 in monolith**

---

## Service Discovery Pattern

```
┌─────────────────────────────────────────────────────┐
│              CONSUL SERVICE REGISTRY                 │
│                                                      │
│  ┌────────────────────────────────────────────┐    │
│  │  Service Name: user-service                │    │
│  │  Address: user-service-abc123              │    │
│  │  Port: 8081                                │    │
│  │  Health: /health (every 10s)               │    │
│  │  Status: ✓ Passing                         │    │
│  └────────────────────────────────────────────┘    │
│                                                      │
│  ┌────────────────────────────────────────────┐    │
│  │  Service Name: menu-service                │    │
│  │  Address: menu-service-def456              │    │
│  │  Port: 8082                                │    │
│  │  Health: /health (every 10s)               │    │
│  │  Status: ✓ Passing                         │    │
│  └────────────────────────────────────────────┘    │
│                                                      │
│  ┌────────────────────────────────────────────┐    │
│  │  Service Name: order-service               │    │
│  │  Address: order-service-ghi789             │    │
│  │  Port: 8083                                │    │
│  │  Health: /health (every 10s)               │    │
│  │  Status: ✓ Passing                         │    │
│  └────────────────────────────────────────────┘    │
│                                                      │
└──────────▲────────────────────────▲─────────────────┘
           │                        │
           │ Register               │ Discover
           │                        │
    ┌──────┴──────┐         ┌──────┴──────┐
    │   Service   │         │   Service   │
    │   Startup   │         │ Calling API │
    └─────────────┘         └─────────────┘
```

---

## Data Isolation Pattern

```
┌─────────────────────────────────────────────────────┐
│               DATABASE PER SERVICE                   │
└─────────────────────────────────────────────────────┘

         ┌──────────┐      ┌──────────┐      ┌──────────┐
         │   User   │      │   Menu   │      │  Order   │
         │ Service  │      │ Service  │      │ Service  │
         └────┬─────┘      └────┬─────┘      └────┬─────┘
              │                 │                  │
              │ Owns            │ Owns             │ Owns
              ▼                 ▼                  ▼
         ┌──────────┐      ┌──────────┐      ┌──────────┐
         │ user_db  │      │ menu_db  │      │ order_db │
         └──────────┘      └──────────┘      └──────────┘
              │                 │                  │
         ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
         │ users   │       │menu_item│       │ orders  │
         │ table   │       │  table  │       │ table   │
         └─────────┘       └─────────┘       └─────────┘
                                              │
                                         ┌────▼────────┐
                                         │order_items  │
                                         │   table     │
                                         └─────────────┘

✗ No foreign keys between databases
✗ No joins across databases  
✓ API calls for data access
✓ Price snapshotting for consistency
```

---

## Health Check Pattern

```
                    ┌─────────────┐
                    │   Consul    │
                    └──────┬──────┘
                           │
                Every 10 seconds
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼──────┐      ┌────▼─────┐      ┌────▼──────┐
    │GET /health│      │GET /health│      │GET /health│
    └───┬──────┘      └────┬─────┘      └────┬──────┘
        │                  │                  │
    ┌───▼──────┐      ┌────▼─────┐      ┌────▼──────┐
    │  User    │      │   Menu   │      │   Order   │
    │ Service  │      │ Service  │      │  Service  │
    └───┬──────┘      └────┬─────┘      └────┬──────┘
        │                  │                  │
    Returns 200       Returns 200        Returns 200
        │                  │                  │
    ┌───▼──────┐      ┌────▼─────┐      ┌────▼──────┐
    │  Status: │      │  Status: │      │  Status:  │
    │ ✓ Passing│      │ ✓ Passing│      │ ✓ Passing│
    └──────────┘      └──────────┘      └───────────┘

If any service fails 3 health checks:
- Status changes to ✗ Failing
- Removed from service discovery
- No traffic routed to service
```

---

## Failure Handling

```
Scenario: Menu Service Goes Down

┌────────┐
│ Client │
└───┬────┘
    │
    │ 1. POST /api/orders
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 2. Discover menu-service
    │
┌───▼────────┐
│   Consul   │
└───┬────────┘
    │
    │ 3. Error: No healthy instances
    │
┌───▼──────────┐
│Order Service │
└───┬──────────┘
    │
    │ 4. Return 503: Menu service unavailable
    │
┌───▼──────────┐
│ API Gateway  │
└───┬──────────┘
    │
    │ 5. Forward error to client
    │
┌───▼────┐
│ Client │ Receives error but user/menu services still work
└────────┘

Result:
- Orders fail ✗
- User operations work ✓
- Menu browsing fails ✗
- Partial degradation instead of total failure
```

---

## Scaling Pattern

```
                    ┌─────────────┐
                    │   Consul    │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    ┌───▼──────┐      ┌────▼─────┐      ┌────▼──────┐
    │  User    │      │   Menu   │      │   Order   │
    │ Service  │      │ Service  │      │  Service  │
    │ x1       │      │ x3       │      │ x1        │
    └──────────┘      └────┬─────┘      └───────────┘
                           │
                    ┌──────▼───────┐
                    │ Menu Service │
                    │ Instance 1   │
                    └──────────────┘
                    │ Menu Service │
                    │ Instance 2   │
                    └──────────────┘
                    │ Menu Service │
                    │ Instance 3   │
                    └──────────────┘

Menu service scaled to 3 instances because:
- High read traffic (browsing)
- Low write traffic
- Independent scaling possible

To scale:
docker-compose up --scale menu-service=3
```

---

## Component Legend

```
┌──────────┐
│ Service  │  = Microservice (stateless)
└──────────┘

┌──────────┐
│ Database │  = PostgreSQL (stateful)
└──────────┘

┌──────────┐
│  Consul  │  = Service Registry
└──────────┘

    │
    ▼         = HTTP Request/Response

    ─────     = Service Discovery Query

    ✓         = Healthy Status

    ✗         = Failed/Unavailable
```
