FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pb module
COPY pb ./pb

# Copy service module files
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

# Set workdir to service
WORKDIR /app/auth-service

# Download dependencies
RUN go mod download

# Copy service source
COPY auth-service/ .

# Build
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth-service/main .
EXPOSE 8081
CMD ["./main"]