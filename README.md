# SysManix

SysManix is a secure, cross-platform service management API that provides controlled access to system services through a RESTful interface. It bridges the gap between system administration and application management by offering a standardized way to interact with system services across different operating systems.

<div align="center">
  <img src="https://user-images.githubusercontent.com/1234567/example.png" alt="SysManix Logo" width="300" />
  <h3>Cross-Platform Service Management API</h3>
  <p>Secure, efficient, and consistent system service management across Windows and Linux platforms.</p>
</div>

<div align="center">

![License](https://img.shields.io/github/license/toxic-development/sysmanix)
![Version](https://img.shields.io/badge/version-0.1.0-blue)
![Go Version](https://img.shields.io/badge/go-1.23.1-00ADD8)

</div>

## Overview

SysManix provides a unified RESTful API for managing system services across different operating systems. It bridges the gap between system administration and application management with a secure, standardized interface.

### Key Features

- **Cross-Platform**: Works seamlessly on both Windows and Linux
- **Secure Authentication**: JWT-based auth with role-based access control
- **Protected Services**: Built-in safeguards for critical system services
- **RESTful API**: Clean, consistent endpoints for service management
- **Comprehensive Logs**: Detailed logging of service operations
- **Swagger Documentation**: Interactive API documentation included

## Quick Start

### Installation

#### Windows
```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_windows_amd64.exe" -OutFile "SysManix.exe"

# Run the executable
.\SysManix.exe
```

#### Linux
```bash
# Download the latest release
wget https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_linux_amd64 -O sysmanix

# Make it executable
chmod +x sysmanix

# Run the application
./sysmanix
```

### First Steps

1. Access the API at `http://localhost:40200`
2. Get an authentication token:
   ```bash
   curl -X POST http://localhost:40200/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"change-me"}'
   ```
3. Use the token for subsequent requests:
   ```bash
   curl -X GET http://localhost:40200/services \
     -H "Authorization: Bearer your-token-here"
   ```

## Getting Started

To get started with SysManix, check out the [Getting Started Guide](./guides/GETTING_STARTED.md).

## Documentation

For detailed usage instructions and examples, see the guides in the [guides](./guides) directory:

- [Introduction](./guides/INTRODUCTION.md): Overview and concept explanation
- [Installation Guide](./guides/INSTALLATION.md): Detailed installation instructions
- [Quick Start Guide](./guides/QUICKSTART.md): Get up and running quickly
- [Configuration Guide](./guides/CONFIGURATION.md): Configuration options and examples
- [Authentication Guide](./guides/AUTHENTICATION.md): Authentication system explained
- [Service Management Guide](./guides/SERVICE_MANAGEMENT.md): Managing services with SysManix
- [API Reference](./guides/API_REFERENCE.md): Detailed API documentation

## Platform-Specific Guides

- [Windows Setup](./guides/WINDOWS_SETUP.md): Windows-specific configuration
- [Linux Setup](./guides/LINUX_SETUP.md): Linux-specific configuration
- [Systemd Integration](./guides/SYSTEMD_SETUP.md): Running as a systemd service
- [Nginx Configuration](./guides/NGINX_SETUP.md): Setting up with Nginx as a reverse proxy

## API Endpoints

### Authentication
- **POST** `/auth/login`: Authenticate and get JWT token
- **GET** `/auth/tokens`: List your active tokens
- **POST** `/auth/tokens/revoke`: Revoke a specific token
- **POST** `/auth/tokens/refresh`: Refresh your current token

### Service Management
- **GET** `/services`: List all services
- **GET** `/services/status/{service}`: Get service status
- **POST** `/services/start/{service}`: Start a service
- **POST** `/services/stop/{service}`: Stop a service
- **GET** `/services/logs/{service}`: View service logs

### System
- **GET** `/health`: System health information

## Development

### Prerequisites
- Go 1.23.1 or newer
- Access to system service management (admin/root privileges)

### Building from Source
```bash
git clone https://github.com/toxic-development/sysmanix.git
cd sysmanix
go build -o sysmanix ./client
```

### Running Tests
```bash
go test -v ./...
```

## Security

SysManix takes security seriously. Key security features:

- JWT tokens with configurable expiration
- Argon2id password hashing
- Role-based access control
- Protected critical system services
- Comprehensive audit logging

For detailed security information, see our [Security Guide](./guides/SECURITY.md).

## License

SysManix is licensed under the MIT License. See [LICENSE](./LICENSE) for details.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Support

- [Troubleshooting Guide](./guides/TROUBLESHOOTING.md)
- [API Documentation](./guides/API_REFERENCE.md)
- [GitHub Issues](https://github.com/toxic-development/SysManix/issues)
