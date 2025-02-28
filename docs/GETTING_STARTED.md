# Getting Started with ChronoServe

## Prerequisites

- Go 1.23.1 or higher
- Windows or Linux operating system
- Administrative privileges
- Git (for cloning the repository)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/therealtoxicdev/chronoserve.git
cd chronoserve
```

2. Build the project:
```bash
make build
```

## Configuration

On first run, ChronoServe will automatically create a `config.yaml` in the project root with default values. You **must** update the security-sensitive values before running the application.

### Default Configuration Structure

```yaml
server:
  host: "localhost"
  port: 8080
  readTimeout: "15s"
  writeTimeout: "15s"
  maxHeaderBytes: 1048576  # 1MB

auth:
  secretKey: "change-me"      # Must be changed
  tokenDuration: 24h
  issuedBy: "ChronoServe"
  allowedRoles: ["admin", "viewer"]
  users:
    admin:
      username: "admin"
      password: "change-me"   # Must be changed
      roles: ["admin"]
    viewer:
      username: "viewer"
      password: "change-me"   # Must be changed
      roles: ["viewer"]

logging:
  level: "info"
  directory: "logs"
  maxSize: 10        # 10MB
  maxBackups: 5
  maxAge: 30         # 30 days
  compress: true
```

### Platform-Specific Settings

#### Windows
```yaml
windows:
  serviceCommand: "sc"
  logDirectory: "C:\\ProgramData\\ChronoServe\\logs"
  services: {}  # Will be populated with discovered services
```

#### Linux
```yaml
linux:
  serviceCommand: "systemctl"
  logDirectory: "/var/log/chronoserve"
  services: {}  # Will be populated with discovered services
```

## Running ChronoServe

1. Development mode:
```bash
make dev
```

2. Production mode:
```bash
make start
```

## Verifying Installation

1. Check the server health:
```bash
curl http://localhost:8080/health
```

2. Try logging in:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your-password"}'
```

## Security Notes

- The application will refuse to start if default credentials are detected
- All passwords should be changed from their default values
- The JWT secret key must be changed from the default value
- Use secure passwords that meet your organization's requirements

## Next Steps

- Review the [API Documentation](./API.md) for available endpoints
- Configure your [Services](./SERVICES.md)
- Set up [Logging](./LOGGING.md)
- Review [Security Best Practices](./SECURITY.md)

## Troubleshooting

### Common Issues

1. "Security Risk Detected":
   - This means you haven't changed the default security values
   - Update the `secretKey` and user passwords in `config.yaml`

2. Permission Issues:
   - Windows: Run as Administrator
   - Linux: Use sudo or appropriate privileges

3. Port Already in Use:
   - Change the port in `config.yaml`
   - Default is 8080