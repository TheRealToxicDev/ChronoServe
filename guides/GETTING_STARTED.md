# Getting Started with SysManix

## Prerequisites

- Go 1.23.1 or higher (only for building from source)
- Windows or Linux operating system
- Administrative privileges
- Git (for cloning the repository)

## Installation

### Option 1: Prebuilt Binaries (Recommended)

Download the latest prebuilt binary for your platform from the [Releases](https://github.com/toxic-development/SysManix/releases) page:

#### Windows
```powershell
# PowerShell (as Administrator)
Invoke-WebRequest -Uri "https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_windows_amd64.exe" -OutFile "SysManix.exe"
```

#### Linux
```bash
wget https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_linux_amd64
chmod +x SysManix_linux_amd64
```

### Option 2: Build from Source

1. Clone the repository:
```bash
git clone https://github.com/toxic-development/SysManix.git
cd SysManix
```

2. Build the project:
```bash
make build
```

## Configuration

On first run, SysManix will create a `config.yaml` file with default values. You **must** update the security-sensitive values before using the application in production.

### Default Configuration Structure

```yaml
server:
    host: "localhost"
    port: 40200          # Default port
    readTimeout: "15s"
    writeTimeout: "15s"
    maxHeaderBytes: 1048576

auth:
    secretKey: "change-me"      # Must be changed
    tokenDuration: 24h
    issuedBy: "SysManix"
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            username: "admin"
            password: "change-me"   # Will be hashed after first run
            roles:
                - admin
        viewer:
            username: "viewer"
            password: "change-me"   # Will be hashed after first run
            roles:
                - viewer

logging:
    level: "debug"
    directory: "logs"
    maxSize: 10        # 10MB
    maxBackups: 5
    maxAge: 30         # 30 days
    compress: true

# Platform-specific settings
windows:
    serviceCommand: "sc"
    logDirectory: "C:\\ProgramData\\SysManix\\logs"
    services: {}

linux:
    serviceCommand: "systemctl"
    logDirectory: "/var/log/SysManix"
    services: {}
```

## Running SysManix

1. Development mode:
```powershell
make dev
```

2. Production mode:
```powershell
make start
```

## Initial Setup

1. Start the application for the first time:
```powershell
make dev
```

2. The application will:
   - Create default config.yaml
   - Prompt you to update security values
   - Exit for you to make changes

3. Update the config values:
   - Change admin and viewer passwords
   - Set a secure secret key
   - Adjust other settings as needed

4. Restart the application:
```powershell
make dev
```

## Verifying Installation

1. Check the server health:
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/health"
```

2. Try logging in:
```powershell
$body = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

$token = $response.data.token
```

## Security Notes

- Default credentials will be rejected
- Passwords are automatically hashed using Argon2id
- Plain text passwords are removed from config after hashing
- JWT secret key must be changed from default
- Use strong passwords that meet your security requirements

## Common Services

### Windows Services
- Windows Update (`wuauserv`)
- Event Log (`EventLog`)
- Windows Remote Management (`WinRM`)
- Background Intelligent Transfer (`BITS`)

Example viewing service logs:
```powershell
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv?lines=50" -Headers $headers
```

## Troubleshooting

### Common Issues

1. "Security Risk Detected":
   - Update `secretKey` and user passwords in `config.yaml`
   - Restart application

2. "Invalid credentials":
   - Ensure you're using the correct password
   - Check if password was properly hashed after first run

3. Permission Issues:
   - Run as Administrator on Windows
   - Use sudo on Linux

4. Port Already in Use:
   - Change the port in `config.yaml` (default: 40200)
   - Check for other applications using the port

## Next Steps

- Review the [API Documentation](API_REFERENCE.md)
- Read [Complete Documentation](DOCUMENTATION.md)
- Learn about [Security Best Practices](../SECURITY.md)

## Support

For issues and feature requests, please visit the [GitHub repository](https://github.com/toxic-development/SysManix).