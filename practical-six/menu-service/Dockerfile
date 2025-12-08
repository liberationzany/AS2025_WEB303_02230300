# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy proto files and generate code
COPY proto/ ./proto/
RUN apk add --no-cache protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate proto files
RUN cd proto/menu/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           menu.proto

# Copy service files
COPY menu-service/go.mod menu-service/go.sum ./menu-service/
COPY go.mod go.sum ./

WORKDIR /app/menu-service

# Download dependencies
RUN go mod download

# Copy source code
COPY menu-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /menu-service main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /menu-service .

EXPOSE 50052

CMD ["./menu-service"]
