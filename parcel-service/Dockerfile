FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pb module
COPY pb ./pb

# Copy service module files
COPY parcel-service/go.mod parcel-service/go.sum ./parcel-service/

# Set workdir to service
WORKDIR /app/parcel-service

# Download dependencies
RUN go mod download

# Copy service source
COPY parcel-service/ .

# Build
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/parcel-service/main .
EXPOSE 8082
CMD ["./main"]
