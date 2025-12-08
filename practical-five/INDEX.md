# Documentation Index

Welcome to the Student Cafe Microservices Project! This index will help you navigate the documentation.

## Quick Access

| What do you want to do? | Read this file |
|-------------------------|----------------|
| 🚀 **Get started quickly** | [QUICK_START.md](QUICK_START.md) |
| 📋 **Follow step-by-step build process** | [BUILD_CHECKLIST.md](BUILD_CHECKLIST.md) |
| 📖 **Understand the architecture** | [ARCHITECTURE.md](ARCHITECTURE.md) |
| 🔍 **See visual diagrams** | [DIAGRAMS.md](DIAGRAMS.md) |
| ❓ **Troubleshoot problems** | [TROUBLESHOOTING.md](TROUBLESHOOTING.md) |
| 📝 **Get project overview** | [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) |
| 📚 **Read full documentation** | [README.md](README.md) |
| 📄 **See original practical guide** | [practical5.md](practical5.md) |

---

## Documentation Files Overview

### 1. [QUICK_START.md](QUICK_START.md)
**Purpose**: Get up and running in 5 minutes  
**Read if**: You want to quickly test the application  
**Contains**:
- Build and start commands
- Test commands
- Verification steps
- Quick troubleshooting

### 2. [BUILD_CHECKLIST.md](BUILD_CHECKLIST.md)
**Purpose**: Comprehensive build and verification checklist  
**Read if**: You want to ensure everything is set up correctly  
**Contains**:
- Pre-flight system checks
- Step-by-step build process
- Verification procedures
- Testing procedures
- Success criteria

### 3. [README.md](README.md)
**Purpose**: Complete project documentation  
**Read if**: You want to understand everything about the project  
**Contains**:
- Project structure
- Architecture overview
- Service descriptions
- Getting started guide
- Testing instructions
- Comparison with monolith
- Next steps

### 4. [ARCHITECTURE.md](ARCHITECTURE.md)
**Purpose**: Detailed architecture decisions and rationale  
**Read if**: You want to understand WHY things are designed this way  
**Contains**:
- Service boundaries justification
- Inter-service communication patterns
- Database-per-service pattern explanation
- API Gateway pattern
- Monolith vs Microservices comparison
- When to use each approach
- Future improvements

### 5. [DIAGRAMS.md](DIAGRAMS.md)
**Purpose**: Visual representation of the system  
**Read if**: You're a visual learner or need to present the architecture  
**Contains**:
- Overall system architecture diagram
- Simple request flow (menu list)
- Complex request flow (create order)
- Service discovery pattern
- Data isolation pattern
- Health check pattern
- Failure handling
- Scaling pattern

### 6. [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
**Purpose**: High-level project overview  
**Read if**: You need a quick summary of what was built  
**Contains**:
- Project structure
- Services overview
- Key features implemented
- Technical stack
- Getting started
- Learning outcomes
- Success criteria

### 7. [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
**Purpose**: Solutions to common problems  
**Read if**: Something isn't working  
**Contains**:
- 14 common issues with solutions
- Diagnostic commands
- Quick reference table
- Full reset instructions
- Getting help section

### 8. [practical5.md](practical5.md)
**Purpose**: Original practical assignment guide  
**Read if**: You want to understand the academic requirements  
**Contains**:
- Learning outcomes
- Theoretical background
- Step-by-step instructions
- Submission requirements
- Grading criteria

---

## Recommended Reading Order

### For First-Time Users
1. **PROJECT_SUMMARY.md** - Get overview (5 minutes)
2. **QUICK_START.md** - Run the project (10 minutes)
3. **DIAGRAMS.md** - Understand visually (10 minutes)
4. **README.md** - Deep dive (30 minutes)

### For Systematic Build
1. **BUILD_CHECKLIST.md** - Follow step by step (30 minutes)
2. **TROUBLESHOOTING.md** - If issues arise
3. **ARCHITECTURE.md** - Understand design decisions (20 minutes)

### For Academic Submission
1. **practical5.md** - Understand requirements
2. **BUILD_CHECKLIST.md** - Build the project
3. **ARCHITECTURE.md** - Understand design
4. **README.md** - Complete documentation
5. Write reflection essay (requirements in practical5.md)

### For Troubleshooting
1. **TROUBLESHOOTING.md** - Find your issue
2. **QUICK_START.md** - Verify basic setup
3. **BUILD_CHECKLIST.md** - Systematic verification

---

## Code Structure

### Services

#### User Service
```
user-service/
├── models/user.go           # User data model
├── handlers/user_handlers.go # HTTP handlers
├── database/db.go           # Database connection
├── main.go                  # Service entry point
├── Dockerfile               # Container definition
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

#### Menu Service
```
menu-service/
├── models/menu.go           # Menu item model
├── handlers/menu_handlers.go # HTTP handlers
├── database/db.go           # Database connection
├── main.go                  # Service entry point
├── Dockerfile               # Container definition
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

#### Order Service
```
order-service/
├── models/order.go          # Order models
├── handlers/order_handlers.go # HTTP handlers with inter-service calls
├── database/db.go           # Database connection
├── main.go                  # Service entry point
├── Dockerfile               # Container definition
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

#### API Gateway
```
api-gateway/
├── main.go                  # Gateway with routing and discovery
├── Dockerfile               # Container definition
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

#### Monolith
```
student-cafe-monolith/
├── models/                  # All models together
│   ├── user.go
│   ├── menu.go
│   └── order.go
├── handlers/                # All handlers together
│   ├── user_handlers.go
│   ├── menu_handlers.go
│   └── order_handlers.go
├── database/db.go           # Single database
├── main.go                  # Monolithic entry point
├── Dockerfile               # Container definition
├── docker-compose.yml       # Standalone compose file
├── go.mod                   # Go dependencies
└── go.sum                   # Dependency checksums
```

---

## Testing Resources

### Manual Testing
- Use commands in **QUICK_START.md**
- Follow test scenarios in **README.md**

### Automated Testing
- Run **test-microservices.ps1** script
- Check **BUILD_CHECKLIST.md** for verification steps

### Visual Testing
- Open Consul UI: http://localhost:8500
- Monitor service health
- View service registrations

---

## Key Concepts Explained

### Service Discovery (Consul)
- Explained in: **ARCHITECTURE.md** (page 3)
- Diagram in: **DIAGRAMS.md** (Service Discovery Pattern)
- Example in: **README.md** (Part 7)

### Inter-Service Communication
- Explained in: **ARCHITECTURE.md** (page 2)
- Diagram in: **DIAGRAMS.md** (Complex Request Flow)
- Code in: `order-service/handlers/order_handlers.go`

### Database-Per-Service
- Explained in: **ARCHITECTURE.md** (page 2-3)
- Diagram in: **DIAGRAMS.md** (Data Isolation Pattern)
- Implementation in: Each service's `database/db.go`

### API Gateway
- Explained in: **ARCHITECTURE.md** (page 3)
- Diagram in: **DIAGRAMS.md** (Overall System Architecture)
- Code in: `api-gateway/main.go`

### Health Checks
- Explained in: **README.md** (Part 7)
- Diagram in: **DIAGRAMS.md** (Health Check Pattern)
- Implementation in: Each service's `main.go`

---

## Common Tasks

### Start the Project
```powershell
docker-compose up --build
```
See: **QUICK_START.md**

### Test the Project
```powershell
.\test-microservices.ps1
```
See: **BUILD_CHECKLIST.md** (Step 13)

### View Logs
```powershell
docker-compose logs -f order-service
```
See: **TROUBLESHOOTING.md** (Diagnostic Commands)

### Stop the Project
```powershell
docker-compose down
```
See: **QUICK_START.md** (Stop Services)

### Clean Everything
```powershell
docker-compose down -v
```
See: **BUILD_CHECKLIST.md** (Cleanup)

---

## URLs Reference

| Service | URL | Purpose |
|---------|-----|---------|
| API Gateway | http://localhost:8080 | Main entry point |
| User Service | http://localhost:8081 | Direct access (bypass gateway) |
| Menu Service | http://localhost:8082 | Direct access (bypass gateway) |
| Order Service | http://localhost:8083 | Direct access (bypass gateway) |
| Consul UI | http://localhost:8500 | Service discovery dashboard |
| Monolith | http://localhost:8090 | Monolithic comparison |

---

## API Endpoints

### Via API Gateway (Recommended)
```
POST   /api/users              Create user
GET    /api/users/{id}         Get user
GET    /api/menu               List menu items
POST   /api/menu               Create menu item
GET    /api/menu/{id}          Get menu item
POST   /api/orders             Create order
GET    /api/orders             List orders
```

See: **README.md** (Testing Through Gateway)

---

## Learning Path

### Beginner Level
1. Read **PROJECT_SUMMARY.md**
2. Run through **QUICK_START.md**
3. View **DIAGRAMS.md**

### Intermediate Level
1. Read **README.md** completely
2. Follow **BUILD_CHECKLIST.md**
3. Experiment with test scenarios
4. Read **ARCHITECTURE.md**

### Advanced Level
1. Review all code in each service
2. Understand **ARCHITECTURE.md** deeply
3. Modify code and test changes
4. Implement suggested improvements
5. Deploy to Kubernetes (future work)

---

## Support

### If You're Stuck
1. Check **TROUBLESHOOTING.md** first
2. Verify with **BUILD_CHECKLIST.md**
3. Review error messages in logs: `docker-compose logs`
4. Check service health: http://localhost:8500

### If You Want to Learn More
1. Read **ARCHITECTURE.md** for design decisions
2. Review code comments in each service
3. Read original **practical5.md** for theory
4. Explore suggested next steps in **README.md**

---

## File Sizes (Approximate)

- **README.md**: 15 KB - Comprehensive documentation
- **QUICK_START.md**: 5 KB - Quick reference
- **BUILD_CHECKLIST.md**: 12 KB - Detailed checklist
- **ARCHITECTURE.md**: 10 KB - Design decisions
- **DIAGRAMS.md**: 8 KB - Visual diagrams
- **PROJECT_SUMMARY.md**: 7 KB - Overview
- **TROUBLESHOOTING.md**: 10 KB - Problem solving
- **practical5.md**: 45 KB - Original assignment

**Total Documentation**: ~112 KB of comprehensive guides

---

## Version Information

- **Project Version**: 1.0
- **Go Version**: 1.23
- **Docker Compose Version**: 3.8+
- **PostgreSQL Version**: 13
- **Consul Version**: Latest

---

## Contributing

If you enhance this project:
1. Update relevant documentation files
2. Add new diagrams to **DIAGRAMS.md** if needed
3. Update **PROJECT_SUMMARY.md** with new features
4. Add troubleshooting entries if you encounter new issues

---

## License

This is an academic project for educational purposes.

---

## Acknowledgments

- Original practical guide: **practical5.md**
- Architecture patterns: Microservices.io
- Service discovery: HashiCorp Consul
- Web framework: go-chi

---

**Ready to start?** → Head to [QUICK_START.md](QUICK_START.md)

**Need help?** → Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

**Want to understand?** → Read [ARCHITECTURE.md](ARCHITECTURE.md)
