.PHONY: build run clean test docker lint

build:
	@echo Building SysManix...
	go build -o bin/sysmanix ./client/main.go

start:
	@echo Starting SysManix Linux in Production...
	./bin/sysmanix

start-windows:
	@echo Starting SysManix Windows in Production...
	./bin/sysmanix.exe

dev:
	go run ./client/main.go

# Build Documentation
docs:
	@echo "Generating OpenAPI documentation"
	@swag i -g docs.go --dir api
	@echo "Swagger docs generated in $(SWAGGER_DIR)"

# Clean up the project
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	CGO_ENABLED=1 go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# For Windows PowerShell - separate targets as PowerShell has different syntax
test-windows:
	powershell -Command "$$env:CGO_ENABLED=1; go test -v -race ./..."

test-coverage-windows:
	powershell -Command "$$env:CGO_ENABLED=1; go test -v -race -coverprofile=coverage.out ./..."
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

docker:
	docker build -t sysmanix:$(VERSION) .
