.PHONY: build run clean test docker lint

build:
	@echo Building ChronoServe...
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

clean:
	rm -rf bin/
	go clean

test:
	go test -v -race ./...

lint:
	golangci-lint run

docker:
	docker build -t sysmanix:$(VERSION) .