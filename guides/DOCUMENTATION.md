# SysManix Documentation

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [Authentication System](#authentication-system)
4. [Service Management](#service-management)
5. [Logging System](#logging-system)
6. [Configuration](#configuration)
7. [API Reference](#api-reference)

## Architecture Overview

SysManix follows a modular architecture designed for cross-platform service management:

```plaintext
SysManix/
├── api/           # HTTP routes and handlers
├── middleware/    # Authentication and request processing
├── services/      # OS-specific service management
├── utils/         # Shared utilities
└── client/        # Main application entry point
```

## Core Components

### API Layer (`/api`)

- Route definitions and handlers
- Health monitoring
- Service management endpoints

```go
routes := []Route{
    {Path: "health", Handler: utils.HealthCheck, RequireAuth: false},
    {Path: "auth/login", Handler: middleware.HandleLogin, RequireAuth: false},
    {Path: "services", Handler: serviceHandler.ListServices, RequireAuth: true, Roles: []string{"admin", "viewer"}},
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
Get-WinEvent         # View service logs
```

#### Linux
```bash
systemctl list-units --type=service  # List services
systemctl start servicename          # Start service
systemctl stop servicename           # Stop service
journalctl                          # View service logs
```

## Authentication System

### Password Security
Passwords are secured using Argon2id hashing with the following configuration:

```go
type PasswordConfig struct {
    time    uint32 // 1
    memory  uint32 // 64 * 1024
    threads uint8  // 4
    keyLen  uint32 // 32
}
```

### JWT Token Structure

```json
{
  "uid": "user-id",
  "roles": ["admin", "viewer"],
  "exp": 1735689600,
  "iat": 1735603200,
  "iss": "SysManix"
}
```

### Role-Based Access

Two primary roles:
- `admin`: Full access to all endpoints
- `viewer`: Read-only access to services

### Password Storage Flow
1. Initial configuration:
```yaml
users:
    admin:
        username: "admin"
        password: "your-password"  # Plain text (temporary)
        roles: ["admin"]
```

2. After first run:
```yaml
users:
    admin:
        username: "admin"
        password_hash: "$argon2id$v=19$m=65536,t=1,p=4$..."  # Secure hash
        roles: ["admin"]
```

## Service Management

### Service Operations

| Operation | Endpoint | Required Role | Description |
|-----------|----------|---------------|-------------|
| List | GET /services | admin, viewer | List all services |
| Status | GET /services/status/{name} | admin, viewer | Get service status |
| Start | POST /services/start/{name} | admin | Start a service |
| Stop | POST /services/stop/{name} | admin | Stop a service |
| Logs | GET /services/logs/{name} | admin, viewer | View service logs |

### Common Windows Services
- Windows Update (`wuauserv`)
- Event Log (`EventLog`)
- Windows Remote Management (`WinRM`)
- Background Intelligent Transfer (`BITS`)

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
    level: debug
    directory: logs
    maxSize: 10        # 10MB
    maxBackups: 5
    maxAge: 30         # 30 days
    compress: true
```

### Log Files
- `app.log`: Application logs
- `auth.log`: Authentication and authorization logs
- `access.log`: HTTP request logs

## Configuration

### Complete Structure

```yaml
server:
    host: localhost
    port: 40200
    readTimeout: 15s
    writeTimeout: 15s
    maxHeaderBytes: 1048576

auth:
    secretKey: your-secure-key
    tokenDuration: 24h0m0s
    issuedBy: SysManix
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            username: admin
            password_hash: "$argon2id$v=19$..."
            roles:
                - admin

linux:
    serviceCommand: systemctl
    logDirectory: /var/log/SysManix
    services: {}

windows:
    serviceCommand: sc
    logDirectory: C:\ProgramData\SysManix\logs
    services: {}

logging:
    level: debug
    directory: logs
    maxSize: 10
    maxBackups: 5
    maxAge: 30
    compress: true
```

## Security Considerations

1. Password Security
   - Argon2id hashing for all passwords
   - Automatic hash generation on first run
   - No plain text passwords stored
   - Secure password verification using constant-time comparison

2. Authentication
   - JWT tokens with expiration
   - Role-based access control
   - Secure token validation

3. Configuration
   - Secure storage of sensitive values
   - Automatic removal of plain text passwords
   - Required changes for default credentials

## API Reference

See [API_REFERENCE.md](API_REFERENCE.md) for detailed API documentation.