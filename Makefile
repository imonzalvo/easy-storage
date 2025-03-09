.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build docker image
docker-build:
	docker build -t easy-storage .

# Run with docker-compose
docker-up:
	docker-compose up -d

# Stop docker containers
docker-down:
	docker-compose down

# Get project dependencies
deps:
	go mod download

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs/api