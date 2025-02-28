# ChronoServe Documentation

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [Authentication System](#authentication-system)
4. [Service Management](#service-management)
5. [Logging System](#logging-system)
6. [Configuration](#configuration)
7. [API Reference](#api-reference)

## Architecture Overview

ChronoServe follows a modular architecture designed for cross-platform service management:

```plaintext
ChronoServe/
├── api/           # HTTP routes and handlers
├── middleware/    # Authentication and request processing
├── services/      # OS-specific service management
├── utils/         # Shared utilities
└── client/        # Main application entry point
```

## Core Components

### API Layer (`api/`)

- Route definitions and handlers
- Health monitoring
- Service management endpoints

```go
routes := []Route{
    {Path: "health", Handler: utils.HealthCheck, RequireAuth: false},
    {Path: "auth/login", Handler: middleware.HandleLogin, RequireAuth: false},
    {Path: "services", Handler: serviceHandler.ListServices, RequireAuth: true, Roles: []string{"admin", "viewer"}},
    // ...more routes
}
```

### Middleware Layer (`middleware/`)

- JWT authentication
- Role-based access control
- Request logging
- Error handling

### Service Management (`services/`)

Platform-specific implementations:

#### Windows
```powershell
sc query              # List services
sc start servicename  # Start service
sc stop servicename   # Stop service
```

#### Linux
```bash
systemctl list-units --type=service  # List services
systemctl start servicename          # Start service
systemctl stop servicename           # Stop service
```

## Authentication System

### JWT Token Structure

```json
{
  "uid": "user-id",
  "roles": ["admin", "viewer"],
  "exp": 1735689600,
  "iat": 1735603200,
  "iss": "ChronoServe"
}
```

### Role-Based Access

Two primary roles:
- `admin`: Full access to all endpoints
- `viewer`: Read-only access to services

## Service Management

### Service Operations

| Operation | Endpoint | Required Role | Description |
|-----------|----------|---------------|-------------|
| List | GET /services | admin, viewer | List all services |
| Status | GET /services/status/{name} | admin, viewer | Get service status |
| Start | POST /services/start/{name} | admin | Start a service |
| Stop | POST /services/stop/{name} | admin | Stop a service |
| Logs | GET /services/logs/{name} | admin, viewer | View service logs |

## Logging System

### Log Levels

```go
const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)
```

### Log Configuration

```yaml
logging:
  level: "info"
  directory: "logs"
  maxSize: 10        # 10MB
  maxBackups: 5
  maxAge: 30         # 30 days
  compress: true
```

## Configuration

### Structure Overview

```yaml
server:
  host: "localhost"
  port: 8080
  readTimeout: "15s"
  writeTimeout: "15s"

auth:
  secretKey: "your-secure-key"
  tokenDuration: 24h
  issuedBy: "ChronoServe"
  users:
    admin:
      username: "admin"
      password: "secure-password"
      roles: ["admin"]

logging:
  level: "info"
  directory: "logs"
```

## API Reference

### Authentication

#### POST /auth/login
```json
Request:
{
  "username": "admin",
  "password": "your-password"
}

Response:
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJ...",
    "roles": ["admin"]
  }
}
```

### Service Management

#### GET /services
```json
Response:
{
  "status": "success",
  "data": {
    "services": [
      {
        "name": "Service1",
        "status": "running",
        "description": "Example service"
      }
    ]
  }
}
```

### Error Responses

```json
{
  "status": "error",
  "message": "Error description",
  "code": 400
}
```

## Security Considerations

1. JWT Token Security
   - Tokens expire after configured duration
   - Secure secret key required
   - HTTPS recommended for production

2. Password Security
   - Default credentials must be changed
   - Passwords stored securely
   - Rate limiting for login attempts

3. Access Control
   - Role-based authorization
   - Principle of least privilege
   - Audit logging