.PHONY: run-auth run-user run-parcel run-tracking run-gateway tidy-all help gen-grpc


help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy-all: ## Run go mod tidy for all services
	@echo "Tidying all modules..."
	cd auth-service && go mod tidy
	cd user-service && go mod tidy
	cd parcel-service && go mod tidy
	cd tracking-service && go mod tidy
	cd api-gateway && go mod tidy
	@echo "Done."

run-auth: ## Run the Auth Service
	cd auth-service && go run cmd/main.go

run-user: ## Run the User Service
	cd user-service && go run cmd/main.go

run-parcel: ## Run the Parcel Service
	cd parcel-service && go run cmd/main.go

run-tracking: ## Run the Tracking Service
	cd tracking-service && go run cmd/main.go

run-gateway: ## Run the API Gateway
	cd api-gateway && go run cmd/main.go

gen-grpc: ## Generate gRPC code from proto files
	@echo "Generating gRPC code..."
	protoc --proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/**/*.proto
	@echo "Done."
