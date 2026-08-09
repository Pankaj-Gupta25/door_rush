FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pb module
COPY pb ./pb

# Copy service module files
COPY tracking-service/go.mod tracking-service/go.sum ./tracking-service/

# Set workdir to service
WORKDIR /app/tracking-service

# Download dependencies
RUN go mod download

# Copy service source
COPY tracking-service/ .

# Build
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tracking-service/main .
EXPOSE 8083
CMD ["./main"]