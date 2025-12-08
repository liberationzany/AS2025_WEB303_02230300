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

# Copy gateway files
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/

WORKDIR /app/api-gateway

# Download dependencies
RUN go mod download

# Copy source code
COPY api-gateway/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
