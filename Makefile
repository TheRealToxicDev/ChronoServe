.PHONY: build run clean test docker lint

build:
	@echo Building ChronoServe...
	go build -o bin/chronoserve ./client/main.go

build-all: build-windows build-linux

build-windows:
	@echo Building ChronoServe for Windows...
	set "GOOS=windows" && set "GOARCH=amd64" && go build -o bin/chronoserve_windows_amd64.exe ./client/main.go

build-linux:
	@echo Building ChronoServe for Linux...
	set "GOOS=linux" && set "GOARCH=amd64" && go build -o bin/chronoserve_linux_amd64 ./client/main.go

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