#!/bin/bash

# Build the gateway
echo "Building API Gateway..."
go build -o api-gateway cmd/main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Start the gateway in background
echo "Starting API Gateway..."
./api-gateway &
GATEWAY_PID=$!

# Wait for it to start
sleep 2

echo "Running Tests..."

# 1. Test Public Route (Auth) - Expecting connection refused/error from downstream, but Gateway should handle it
echo "Testing Public Route (/auth/user/signin)..."
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/auth/user/signin)
echo "Response Code: $RESPONSE"
# We expect 502 because auth service is not running, or 404/200 if it is. The key is we get a response.
if [ "$RESPONSE" == "000" ]; then
    echo "FAILED: Gateway not reachable"
else
    echo "PASSED: Gateway reachable"
fi

# 2. Test Protected Route (User Profile) WITHOUT Token
echo "Testing Protected Route (/users/profile) WITHOUT Token..."
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/users/profile)
echo "Response Code: $RESPONSE"
if [ "$RESPONSE" == "401" ]; then
    echo "PASSED: Correctly returned 401"
else
    echo "FAILED: Expected 401, got $RESPONSE"
fi

# 3. Test Protected Route WITH Invalid Token
echo "Testing Protected Route (/users/profile) WITH Invalid Token..."
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer invalidtoken" http://localhost:8080/users/profile)
echo "Response Code: $RESPONSE"
if [ "$RESPONSE" == "401" ]; then
    echo "PASSED: Correctly returned 401 for invalid token"
else
    echo "FAILED: Expected 401, got $RESPONSE"
fi

# Cleanup
echo "Stopping API Gateway..."
kill $GATEWAY_PID
