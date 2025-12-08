# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy proto files and generate code
COPY proto/ ./proto/
RUN apk add --no-cache protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate all proto files
RUN cd proto/user/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           user.proto && \
    cd ../../menu/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           menu.proto && \
    cd ../../order/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           order.proto

# Copy service files
COPY order-service/go.mod order-service/go.sum ./order-service/
COPY go.mod go.sum ./

WORKDIR /app/order-service

# Download dependencies
RUN go mod download

# Copy source code
COPY order-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /order-service .

EXPOSE 50053

CMD ["./order-service"]
