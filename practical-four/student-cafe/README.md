# Student Cafe - Kubernetes Microservices with Kong Gateway

A complete microservices application demonstrating Kubernetes deployment, service discovery with Consul, and API gateway routing with Kong.

## Project Structure

```
student-cafe/
├── food-catalog-service/          # Go microservice for food items
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── order-service/                 # Go microservice for order management
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── cafe-ui/                       # React frontend
│   ├── src/
│   │   ├── App.js
│   │   ├── App.css
│   │   ├── index.js
│   │   └── index.css
│   ├── public/
│   │   └── index.html
│   ├── package.json
│   └── Dockerfile
├── app-deployment.yaml            # Kubernetes deployment manifests
├── kong-ingress.yaml              # Kong API Gateway ingress config
└── README.md                       # This file
```

## Architecture Overview

- **Frontend (React.js)**: Single-page application with menu and cart functionality
- **food-catalog-service (Go)**: Provides list of available food items
- **order-service (Go)**: Handles order creation and management
- **Service Discovery (Consul)**: Enables service-to-service communication
- **API Gateway (Kong)**: Routes external traffic to appropriate services
- **Orchestration (Kubernetes)**: Container management and deployment

## Prerequisites

Ensure you have installed:
- Go 1.23+
- Node.js & npm (LTS version)
- Docker
- Minikube
- kubectl
- Helm

## Getting Started

### 1. Start Minikube

```bash
minikube start --cpus 4 --memory 4096
eval $(minikube -p minikube docker-env)
```

### 2. Create Kubernetes Namespace

```bash
kubectl create namespace student-cafe
```

### 3. Deploy Consul (Service Discovery)

```bash
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul --set global.name=consul --namespace student-cafe --set server.replicas=1 --set server.bootstrapExpect=1
```

### 4. Deploy Kong (API Gateway)

```bash
helm repo add kong https://charts.konghq.com
helm repo update
helm install kong kong/kong --namespace student-cafe
```

### 5. Build Docker Images

From the project root directory:

```bash
docker build -t food-catalog-service:v1 ./food-catalog-service
docker build -t order-service:v1 ./order-service
docker build -t cafe-ui:v1 ./cafe-ui
```

### 6. Deploy Application Services

```bash
kubectl apply -f app-deployment.yaml
```

### 7. Configure Kong Ingress

```bash
kubectl apply -f kong-ingress.yaml
```

### 8. Access the Application

Get the Kong proxy URL:

```bash
minikube service -n student-cafe kong-kong-proxy --url
```

Open the returned URL in your web browser. You should see the Student Cafe interface with the menu and cart.

## Testing the Application

### Check Pod Status

```bash
kubectl get pods -n student-cafe
kubectl describe pod <pod-name> -n student-cafe
```

### View Service Endpoints

```bash
kubectl get services -n student-cafe
kubectl get endpoints -n student-cafe
```

### Monitor Logs

```bash
# Food catalog service logs
kubectl logs -f deployment/food-catalog-deployment -n student-cafe

# Order service logs
kubectl logs -f deployment/order-deployment -n student-cafe

# Kong gateway logs
kubectl logs -f deployment/kong-kong -n student-cafe
```

### Test API Endpoints Directly

```bash
# Get minikube IP
MINIKUBE_IP=$(minikube ip)
KONG_PORT=$(kubectl get service -n student-cafe kong-kong-proxy -o jsonpath='{.spec.ports[0].nodePort}')

# Test catalog endpoint
curl http://$MINIKUBE_IP:$KONG_PORT/api/catalog/items

# Test order endpoint
curl -X POST http://$MINIKUBE_IP:$KONG_PORT/api/orders/orders \
  -H "Content-Type: application/json" \
  -d '{"item_ids": ["1", "2"]}'
```

## Troubleshooting

### Go Version Issues
If you get `go: go.mod requires go >= 1.25.0`, update Go version in Dockerfiles and go.mod files to 1.23.

### Docker Images Not Found
Ensure you run `eval $(minikube -p minikube docker-env)` before building images.

### Consul Connection Issues
Services are configured to handle Consul gracefully:
- Consul client address: `consul-server:8500`
- Service registration runs asynchronously
- Errors are logged but don't crash the service

### Kong Not Accessible
1. Verify Kong is running: `kubectl get pods -n student-cafe | grep kong`
2. Check ingress: `kubectl get ingress -n student-cafe`
3. Verify Kong logs: `kubectl logs deployment/kong-kong -n student-cafe`

### Pods Stuck in Pending
1. Check pod details: `kubectl describe pod <pod-name> -n student-cafe`
2. Ensure Docker images are built and tagged correctly
3. Increase minikube resources: `minikube start --cpus 4 --memory 4096`

## API Endpoints

### Food Catalog Service
- `GET /health` - Health check
- `GET /items` - Get all food items

### Order Service
- `GET /health` - Health check
- `POST /orders` - Create new order

### Kong Gateway Routes
- `/api/catalog/...` → food-catalog-service
- `/api/orders/...` → order-service
- `/` → cafe-ui (React frontend)

## Features

✅ Multi-service microservices architecture
✅ Service discovery with Consul
✅ API gateway routing with Kong
✅ Kubernetes deployment
✅ React-based frontend with menu and cart
✅ Order management system
✅ Health checks and monitoring
✅ Container-based deployment

## Debugging Commands

```bash
# Get all resources
kubectl get all -n student-cafe

# Check pod logs with tail
kubectl logs -f <pod-name> -n student-cafe

# Describe pod for details
kubectl describe pod <pod-name> -n student-cafe

# Port forward for direct testing
kubectl port-forward svc/food-catalog-service 8080:8080 -n student-cafe

# Exec into pod
kubectl exec -it <pod-name> -n student-cafe -- sh

# Check ingress status
kubectl describe ingress cafe-ingress -n student-cafe

# Get minikube IP
minikube ip

# SSH into minikube
minikube ssh
```

## Cleanup

To remove all resources:

```bash
# Delete application deployments
kubectl delete -f app-deployment.yaml
kubectl delete -f kong-ingress.yaml

# Delete Kong
helm uninstall kong -n student-cafe

# Delete Consul
helm uninstall consul -n student-cafe

# Delete namespace
kubectl delete namespace student-cafe

# Stop minikube
minikube stop
```

## Next Steps: Resilience Patterns

This project provides a foundation for implementing resilience patterns:
- **Timeout**: Limit how long services wait for responses
- **Retry**: Automatically retry failed requests
- **Circuit Breaker**: Stop sending requests to failing services

These patterns will be implemented in Part 2 of the practical.

## Notes

- Images are built locally and used directly by Minikube
- Service discovery happens automatically via Consul
- Kong automatically discovers services through Kubernetes service endpoints
- All data is in-memory; stopping pods will lose order history

## Support

For issues or questions, refer to the troubleshooting section or check pod logs for error details.
