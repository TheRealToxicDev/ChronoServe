<h2 align='center'>
  <img src="https://img.toxicdev.me/api/user/510065483693817867/IPf3aZXe.png" height='150px' width='150px'/>
  <br> 
</h2>

# ChronoServe
A secure, cross-platform service management API that provides controlled access to system services through a RESTful interface.

## Features

- 🔐 JWT-based authentication with Argon2id password hashing
- 🖥️ Cross-platform support (Windows and Linux)
- 📝 Detailed logging with rotation and compression
- ⚙️ YAML-based configuration system
- 🚦 Health monitoring with detailed metrics
- 🔄 Graceful shutdown handling
- 🛡️ Security-first design with RBAC

## Quick Links

- [Getting Started](./docs/GETTING_STARTED.md)
- [Full Documentation](./docs/DOCUMENTATION.md)
- [API Reference](./docs/API_REFERENCE.md)
- [Security](./SECURITY.md)
- [Contributing](./CONTRIBUTING.md)
- [Troubleshooting](./docs/TROUBLESHOOTING.md)

## API Overview

### Public Endpoints
- `GET /health` - Server health check
- `POST /auth/login` - Authentication endpoint

### Protected Endpoints
- `GET /services` - List all services
- `GET /services/status/{name}` - Get service status
- `POST /services/start/{name}` - Start a service (admin only)
- `POST /services/stop/{name}` - Stop a service (admin only)
- `GET /services/logs/{name}` - View service logs

## Quick Start

```powershell
# Install ChronoServe
Invoke-WebRequest -Uri "https://github.com/therealtoxicdev/chronoserve/releases/latest/download/chronoserve_windows_amd64.exe" -OutFile "chronoserve.exe"

# First run (creates config)
.\chronoserve.exe

# Update config.yaml with your credentials
notepad config.yaml

# Start the server
.\chronoserve.exe
```

### Authentication Example

```powershell
# Login and get token
$body = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

$token = $response.data.token

# Use token to list services
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

## Common Windows Services

```powershell
# View Windows Update service logs
Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv" -Headers $headers

# Check Event Log service status
Invoke-RestMethod -Uri "http://localhost:40200/services/status/EventLog" -Headers $headers
```

## Security Features

- Argon2id password hashing
- Role-based access control
- Automatic plain-text password removal
- Secure configuration handling
- Detailed security logging

## Development

```powershell
# Clone repository
git clone https://github.com/therealtoxicdev/chronoserve.git
cd chronoserve

# Install dependencies
go mod download

# Run in development mode
make dev
```

## License

GNU AFFERO GENERAL PUBLIC LICENSE Version 3 - See [LICENSE](./LICENSE) file for details

## Security

Found a security issue? Please report it privately following our [Security Policy](./SECURITY.md).

## Support

- [Troubleshooting Guide](./docs/TROUBLESHOOTING.md)
- [API Documentation](./docs/API_REFERENCE.md)
- [GitHub Issues](https://github.com/therealtoxicdev/chronoserve/issues)