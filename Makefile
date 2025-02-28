.PHONY: build run clean test docker lint

build:
	@echo Building ChronoServe...
	go build -o bin/chronoserve ./client/main.go

start:
	@echo Starting ChronoServe Linux in Production...
	./bin/chronoserve

start-windows:
	@echo Starting ChronoServe Windows in Production...
	./bin/chronoserve.exe

dev:
	go run ./client/main.go

clean:
	rm -rf bin/
	go clean

test:
	go test -v -race ./...

lint:
	golangci-lint run

docker:
	docker build -t chronoserve:$(VERSION) .