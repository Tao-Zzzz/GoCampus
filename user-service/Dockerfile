# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
COPY config/config.yaml ./config/
EXPOSE 8080 9090
CMD ["./user-service"]