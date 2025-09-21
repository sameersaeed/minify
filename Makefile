.PHONY: build run test clean docker docker-up docker-down deps

APP_NAME := minify

build:
	go build -o $(APP_NAME) .

run:
	go run main.go

test:
	go test ./...

clean:
	go clean
	rm -f $(APP_NAME)

deps:
	go mod download
	go mod tidy

dev: deps
	go run main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run

db-create:
	createdb minify

db-drop:
	dropdb minify

build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(APP_NAME) .

help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  dev        - Run in development mode"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message"
