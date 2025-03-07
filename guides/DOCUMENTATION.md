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

### Service Protection

SysManix implements a protection mechanism for critical system services:

```go
// Windows protected services
var criticalWindowsServices = []string{
    "wininit",           // Windows Start-Up Application
    "csrss",             // Client Server Runtime Process
    "services",          // Services and Controller app
    "lsass",             // Local Security Authority Process
    "winlogon",          // Windows Logon
    "smss",              // Windows Session Manager
    // ...other critical services
}

// Linux protected services
var criticalLinuxServices = []string{
    "systemd",             // Core system daemon
    "systemd-journald",    // Journal logging service
    "systemd-logind",      // Login service
    "systemd-udevd",       // udev management daemon
    "sshd",                // SSH daemon
    "dbus",                // D-Bus system message bus
    // ...other critical services
}
```

When a user attempts to modify a protected service, the system returns a 403 Forbidden response with an appropriate error message.

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

# SysManix Documentation Guide

This document serves as an index for all SysManix documentation and provides an overview of available resources.

## Documentation Structure

SysManix documentation is organized into several guides, each focused on specific aspects of the system:

### Core Documentation

| Guide | Description |
|-------|-------------|
| [Introduction](./INTRODUCTION.md) | Overview of SysManix, its features, and use cases |
| [Getting Started](./GETTING_STARTED.md) | Quick start guide to get up and running |
| [Installation](./INSTALLATION.md) | Detailed installation instructions for all platforms |
| [Quick Start](./QUICKSTART.md) | Fast setup and basic usage examples |
| [Configuration](./CONFIGURATION.md) | Comprehensive configuration options |
| [API Reference](./API_REFERENCE.md) | Complete API endpoint documentation |

### Feature-Specific Guides

| Guide | Description |
|-------|-------------|
| [Authentication](./AUTHENTICATION.md) | User management, JWT tokens, and security |
| [Service Management](./SERVICE_MANAGEMENT.md) | Detailed guide on managing system services |
| [Permissions](./PERMISSIONS.md) | Role-based access control and permission model |
| [Security](./SECURITY.md) | Security best practices and hardening guide |

### Platform-Specific Guides

| Guide | Description |
|-------|-------------|
| [Windows Setup](./WINDOWS_SETUP.md) | Windows-specific installation and configuration |
| [Linux Setup](./LINUX_SETUP.md) | Linux-specific installation and configuration |
| [Systemd Setup](./SYSTEMD_SETUP.md) | Using SysManix with systemd services |
| [Nginx Setup](./NGINX_SETUP.md) | Setting up Nginx as a reverse proxy |

### Advanced Topics

| Guide | Description |
|-------|-------------|
| [Troubleshooting](./TROUBLESHOOTING.md) | Common issues and solutions |
| [Versions](./VERSIONS.md) | Version history and compatibility information |

## Recommended Reading Path

### For New Users

1. [Introduction](./INTRODUCTION.md) - Understand what SysManix is
2. [Getting Started](./GETTING_STARTED.md) - Basic installation and usage
3. [Configuration](./CONFIGURATION.md) - Configure for your environment
4. Platform-specific guide: [Windows Setup](./WINDOWS_SETUP.md) or [Linux Setup](./LINUX_SETUP.md)
5. [Service Management](./SERVICE_MANAGEMENT.md) - Learn how to manage services

### For Administrators

1. [Security](./SECURITY.md) - Security best practices
2. [Authentication](./AUTHENTICATION.md) - User management and authentication
3. [Permissions](./PERMISSIONS.md) - Configure role-based access
4. [Nginx Setup](./NGINX_SETUP.md) - Production deployment with HTTPS
5. [Troubleshooting](./TROUBLESHOOTING.md) - Solve common issues

### For Developers

1. [API Reference](./API_REFERENCE.md) - Complete API documentation
2. [Authentication](./AUTHENTICATION.md) - JWT token usage
3. Code examples in [Service Management](./SERVICE_MANAGEMENT.md)

## Documentation Format

All guides follow a consistent format:

- **Introductory section**: Brief overview of the topic
- **Main content**: Detailed information with examples
- **Platform-specific sections**: Differences between Windows and Linux
- **Examples**: Code samples and usage scenarios
- **Troubleshooting**: Common issues related to the topic

## Contributing to Documentation

We welcome contributions to improve SysManix documentation. Please follow these guidelines:

1. **Clarity** - Write clear, concise explanations
2. **Examples** - Include practical examples for each concept
3. **Accuracy** - Ensure all information is correct and up-to-date
4. **Consistency** - Follow the established format and style
5. **Cross-reference** - Link to other relevant guides where appropriate

To contribute documentation improvements:

1. Fork the repository
2. Make your changes
3. Submit a pull request with a description of your changes

## Finding Information

You can find information in the SysManix documentation in several ways:

- **Navigation**: Use the links in this guide to navigate directly to specific topics
- **Search**: Use your browser's search function (Ctrl+F) within guides
- **Cross-references**: Follow links between guides for related topics
- **API Reference**: For endpoint-specific details, refer to the [API Reference](./API_REFERENCE.md)

If you can't find the information you need, please check the [Troubleshooting](./TROUBLESHOOTING.md) guide for common issues or report a documentation issue on our GitHub repository.

## Additional Resources

In addition to these guides, SysManix provides:

- **CLI Help**: Run `./sysmanix --help` for command-line options
- **Swagger Documentation**: Available at `/swagger/` when enabled in config
- **Health Endpoint**: Check `/health` for system status information
- **GitHub Repository**: Visit our [repository](https://github.com/toxic-development/sysmanix) for latest updates and issues

## Version Information

The current documentation applies to SysManix version 0.1.0 and later. For version-specific information, see the [Versions](./VERSIONS.md) guide.
