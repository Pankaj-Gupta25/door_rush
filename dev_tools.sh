#!/bin/bash

# Source this file to get quick navigation aliases
# Usage: source dev_tools.sh

echo "Loading PDTS-Go Navigation Tools..."

# Root
alias pdts="cd $(pwd)"

# Services
alias cdauth="cd $(pwd)/auth-service"
alias cduser="cd $(pwd)/user-service"
alias cdparcel="cd $(pwd)/parcel-service"
alias cdtrack="cd $(pwd)/tracking-service"
alias cdgate="cd $(pwd)/api-gateway"

echo "Aliases loaded!"
echo "   - cdauth   : Go to Auth Service"
echo "   - cduser   : Go to User Service"
echo "   - cdparcel : Go to Parcel Service"
echo "   - cdtrack  : Go to Tracking Service"
echo "   - cdgate   : Go to API Gateway"
