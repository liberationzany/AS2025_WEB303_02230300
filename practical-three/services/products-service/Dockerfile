FROM golang:1.23-bullseye AS builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./services/products-service

FROM gcr.io/distroless/cc-debian11:latest
COPY --from=builder /server /server
EXPOSE 50052
CMD ["/server"]
