# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy proto files and generate code
COPY proto/ ./proto/
RUN apk add --no-cache protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate proto files
RUN cd proto/user/v1 && \
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           user.proto

# Copy service files
COPY user-service/go.mod user-service/go.sum ./user-service/
COPY go.mod go.sum ./

WORKDIR /app/user-service

# Download dependencies
RUN go mod download

# Copy source code
COPY user-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /user-service .

EXPOSE 50051

CMD ["./user-service"]
