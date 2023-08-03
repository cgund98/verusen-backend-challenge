# Lint the package
lint:
	@echo "Linting go files..."
	@revive -config revive.toml ./...

# Format the go package
format:
	@echo "Formatting go files..."
	@gofmt -w ./internal ./api ./cmd

# Run unit tests
test:
	@echo "Running unit tests..."
	@go test ./internal/...

# Build application
build:
	@echo "Building containers..."
	@docker compose build

# Start the service
start:
	@echo "Starting service..."
	@docker compose up