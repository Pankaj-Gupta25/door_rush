FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pb module
COPY pb ./pb

# Copy service module files
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/

# Set workdir to service
WORKDIR /app/api-gateway

# Download dependencies
RUN go mod download

# Copy service source
COPY api-gateway/ .

# Build
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway/main .
CMD ["./main"]