.PHONY: build run clean test docker lint

VERSION := 0.1.0
LDFLAGS := -X main.version=$(VERSION)

build:
	@echo Building ChronoServe v$(VERSION)...
	go build -ldflags "$(LDFLAGS)" -o bin/chronoserve ./client/main.go

start:
	@echo Starting ChronoServe v$(VERSION) in Production...
	./bin/chronoserve

start-windows:
	@echo Starting ChronoServe v$(VERSION) in Production...
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