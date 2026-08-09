FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pb module
COPY pb ./pb

# Copy service module files
COPY user-service/go.mod user-service/go.sum ./user-service/

# Set workdir to service
WORKDIR /app/user-service

# Download dependencies
RUN go mod download

# Copy service source
COPY user-service/ .

# Build
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service/main .
EXPOSE 8084
CMD ["./main"]

