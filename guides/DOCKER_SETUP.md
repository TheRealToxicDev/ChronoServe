# SysManix Docker Setup Guide

This guide will help you set up and run SysManix using Docker.

## Prerequisites

- Docker installed on your machine
- Docker Compose installed on your machine
- Docker Hub account

## Building the Docker Image

To build the Docker image for SysManix, run the following command:

```sh
make docker-build
```

This command will build the Docker image using the `Dockerfile` in the project root.

## Running the Docker Container

To run the Docker container for SysManix, use Docker Compose:

```sh
make docker-run
```

This command will start the SysManix container using the `docker-compose.yml` file.

## Stopping the Docker Container

To stop the Docker container, run:

```sh
make docker-stop
```

This command will stop and remove the SysManix container.

## Docker Compose Configuration

The `docker-compose.yml` file is configured to build and run the SysManix container. It maps port `40200` and mounts the configuration directory and Docker socket.

```yaml
version: '3'

services:
  sysmanix:
    build: .
    container_name: sysmanix
    ports:
      - "40200:40200"
    volumes:
      - ./config:/app/config
      - /var/run/docker.sock:/var/run/docker.sock
    restart: unless-stopped
```
