# ChronoServe

A secure, cross-platform service management API that provides controlled access to system services through a RESTful interface.

## Features

- 🔐 JWT-based authentication and role-based access control
- 🖥️ Cross-platform support (Windows and Linux)
- 📝 Detailed logging with rotation
- ⚙️ Flexible configuration system
- 🚦 Health monitoring endpoints
- 🔄 Graceful shutdown handling
- 🛡️ Security-first design

## Quick Links

- [Getting Started](./docs/GETTING_STARTED.md)
- [Documentation](./docs/BREAKDOWN.md)
- [API Reference](./docs/API_USAGE.md)
- [Contributing](./CONTRIBUTING.md)
- [Updating](./docs/VERSIONS.md)

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

## Quick Example

```bash
# Get a token
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your-password"}'

# Use the token to list services
curl -X GET http://localhost:8080/services \
  -H "Authorization: Bearer your-token-here"
```

## License

GNU AFFERO GENERAL PUBLIC LICENSE - See [LICENSE](./LICENSE) file for details

## Security

Found a security issue? Please report it privately via our [Security Policy](./SECURITY.md).