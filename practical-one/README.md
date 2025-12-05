**Hands-On Exercise: Building gRPC Microservices with Docker Containers**

This hands-on exercise demonstrates the development and deployment of two interconnected microservices using gRPC technology. The services are packaged as Docker containers and coordinated using Docker Compose for seamless orchestration.

**Project Highlights**
This implementation illustrates several key modern development practices:

- Inter-service communication via gRPC for efficient binary data transfer
- Microservices architecture pattern with independent, scalable components
- Containerization of services using Docker for consistent deployment
- Multi-service management and networking through Docker Compose
- Protocol Buffer schemas ensuring type-safe service contracts

**Documentation Images**

**Image 1: gRPC Client Request**
Demonstration of a successful gRPC call using grpcurl, returning a personalized greeting message that incorporates the current server timestamp.
![grpcurl Test Output](assets/Screenshot%202025-12-04%20124655.png)

**Image 2: Containerized Services Running**
Visual confirmation showing both microservices operating within isolated Docker containers, managed as a cohesive application unit.
![Docker Containers Running](assets/Screenshot%202025-12-04%20131821.png)