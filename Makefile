# Makefile for Hika project

# Directories
API_GATEWAY_DIR := ./api-gateway
APP_DIR := ./app

# Targets
.PHONY: build down

# Build target
build:
	@echo "Initializing Swagger documentation..."
	cd $(API_GATEWAY_DIR) && swag init -d ./cmd/api,./
	@echo "Building Docker containers..."
	cd $(APP_DIR) && docker-compose up --build

# Down target
down:
	@echo "Stopping Docker containers..."
	cd $(APP_DIR) && docker-compose down

# Down target
up:
	@echo "Starting Docker containers..."
	cd $(APP_DIR) && docker-compose up
