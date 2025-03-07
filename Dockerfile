# Use the official Golang image to build the application
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o sysmanix ./client/main.go

# Start a new stage from scratch
FROM alpine:latest

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/sysmanix /usr/local/bin/sysmanix

# Expose port 40200 to the outside world
EXPOSE 40200

# Command to run the executable
CMD ["sysmanix"]
