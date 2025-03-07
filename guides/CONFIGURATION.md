# SysManix Configuration Guide

This guide explains how to configure SysManix to meet your specific requirements.

## Configuration File

SysManix uses a YAML configuration file (`config.yaml`) to control its behavior. The configuration file is automatically created on first run if it doesn't exist, populated with default values.

### Configuration File Location

By default, SysManix looks for its configuration in the following locations:

- **Windows**: In the same directory as the executable
- **Linux**: In the same directory as the executable, or at `/etc/sysmanix/config.yaml` if run as a system service

You can specify a custom configuration file location using the `-config` flag:

```bash
./sysmanix -config /path/to/custom/config.yaml
```

## Configuration Sections

### Server Configuration

Controls the HTTP server settings:

```yaml
server:
  host: "localhost"   # Interface to listen on (use "0.0.0.0" for all interfaces)
  port: 40200         # Port to listen on
  readTimeout: "15s"  # HTTP read timeout
  writeTimeout: "15s" # HTTP write timeout
  maxHeaderBytes: 1048576  # Max header size in bytes (1MB)
```

### Authentication Configuration

Controls the authentication system:

```yaml
auth:
  secretKey: "your-secure-random-string-here"  # JWT signing key (CHANGE THIS!)
  tokenDuration: 24h                          # Token validity period
  issuedBy: "SysManix"                        # Token issuer name
  allowedRoles:                               # Available roles in the system
    - admin
    - viewer
  users:                                      # User definitions
    admin:
      username: "admin"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - admin
    viewer:
      username: "viewer"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - viewer
    custom_user:
      username: "custom_user"
      password: "plain_password"  # Plain password (will be hashed after restart)
      roles:
        - viewer
```

Important notes:
- `secretKey` should be a strong random string (at least 32 characters)
- `tokenDuration` uses Go's duration format (e.g., "24h", "30m", "1h30m")
- You can provide either `password_hash` or `password` (which will be automatically hashed)
- User IDs will be automatically generated if not provided
- Additional user fields such as `avatarUrl`, `bannerUrl`, `displayName`, and `bio` can be added as needed

### Additional User Fields

You can add the following optional fields to each user:

- `avatarUrl`: URL to the user's avatar image
- `bannerUrl`: URL to the user's banner image
- `displayName`: Display name for the user
- `bio`: Short biography or description for the user

Example:

```yaml
auth:
  users:
    custom_user:
      username: "custom_user"
      password: "plain_password"  # Plain password (will be hashed after restart)
      roles:
        - viewer
      avatarUrl: "https://example.com/avatar.jpg"
      bannerUrl: "https://example.com/banner.jpg"
      displayName: "Custom User"
      bio: "This is a custom user."
```

### Operating System Specific Configuration

Settings specific to each supported operating system:

```yaml
windows:
  serviceCommand: "sc"                         # Windows service command
  logDirectory: "C:\\ProgramData\\SysManix\\logs"  # Log directory

linux:
  serviceCommand: "systemctl"                  # Linux service command
  logDirectory: "/var/log/SysManix"            # Log directory
```

### Logging Configuration

Controls the logging behavior:

```yaml
logging:
  level: "info"      # Log level (debug, info, warn, error)
  directory: "logs"  # Log directory (relative to executable)
  maxSize: 10        # Max size per log file in MB
  maxBackups: 5      # Number of old log files to keep
  maxAge: 30         # Max age of log files in days
  compress: true     # Compress old log files
```

### API Documentation Configuration

Controls the Swagger API documentation:

```yaml
api:
  enableSwagger: true         # Enable/disable Swagger UI
  swaggerPath: "/swagger/"    # Path to Swagger UI
  version: "1.0.0"            # API version
  title: "SysManix API"       # API title for documentation
  description: "Cross-platform service management API"  # API description
```

## Advanced Configuration

### Extended User Roles

Create custom roles for fine-grained permission control:

```yaml
auth:
  allowedRoles:
    - admin       # Full access
    - viewer      # Read-only access
    - operator    # Can start/stop services but not view tokens
    - db_admin    # Can only manage database services
    - web_admin   # Can only manage web server services

  users:
    web_user:
      username: "web_user"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - web_admin

    db_user:
      username: "db_user"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - db_admin
```

## Environment Variables

SysManix also supports environment variables to override configuration values:

| Environment Variable | Configuration Path | Description |
|---------------------|-------------------|-------------|
| `SYSMANIX_SERVER_PORT` | server.port | Server port |
| `SYSMANIX_SERVER_HOST` | server.host | Server host |
| `SYSMANIX_AUTH_SECRET_KEY` | auth.secretKey | JWT secret key |
| `SYSMANIX_LOG_LEVEL` | logging.level | Log level |

Example usage:

```bash
# Set server port via environment variable
export SYSMANIX_SERVER_PORT=8080
./sysmanix
```

## Configuration Best Practices

### Security Recommendations

1. **Generate a strong secret key**:
   ```bash
   # Linux
   openssl rand -base64 32

   # Windows PowerShell
   [Convert]::ToBase64String((New-Object byte[] 32) | ForEach-Object { $_ = Get-Random -Minimum 0 -Maximum 256 })
   ```

2. **Use environment variables for secrets** in production environments
3. **Restrict access to the configuration file** (chmod 600 on Linux)
4. **Separate permissions** using custom roles instead of giving everyone admin access
5. **Change default passwords** immediately after installation

### Performance Tuning

1. **Adjust log levels** based on environment:
   - Production: `info` or `warn`
   - Development/Debugging: `debug`

2. **Configure appropriate timeouts**:
   ```yaml
   server:
     readTimeout: "30s"  # Increase for slow networks
     writeTimeout: "30s" # Increase for operations that take longer
   ```

3. **Enable log compression** to save disk space:
   ```yaml
   logging:
     compress: true
   ```

### Development Environment

For development environments, use a more permissive configuration:

```yaml
server:
  host: "localhost"
  port: 8080

auth:
  tokenDuration: 24h  # Longer tokens for convenience during development

logging:
  level: "debug"  # More verbose logging

api:
  enableSwagger: true  # Enable API documentation
```

### Production Environment

For production environments, focus on security and stability:

```yaml
server:
  host: "127.0.0.1"  # Only accessible locally (use with reverse proxy)

auth:
  tokenDuration: 8h  # Shorter token lifetime

logging:
  level: "info"
  maxSize: 50  # Larger log files
  maxBackups: 10  # Keep more backups

api:
  enableSwagger: false  # Disable Swagger in production
```

## Full Configuration Example

Below is a complete configuration example with all available options:

```yaml
server:
  host: "localhost"
  port: 40200
  readTimeout: "15s"
  writeTimeout: "15s"
  maxHeaderBytes: 1048576

auth:
  secretKey: "your-secure-random-string-here"
  tokenDuration: 12h
  issuedBy: "SysManix"
  allowedRoles:
    - admin
    - viewer
    - operator
    - db_admin
  users:
    admin:
      username: "admin"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - admin
    viewer:
      username: "viewer"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - viewer
    operator:
      username: "operator"
      password: "change-me"  # Plain password (will be hashed after restart)
      roles:
        - operator

windows:
  serviceCommand: "sc"
  logDirectory: "C:\\ProgramData\\SysManix\\logs"

linux:
  serviceCommand: "systemctl"
  logDirectory: "/var/log/SysManix"

logging:
  level: "info"
  directory: "logs"
  maxSize: 10
  maxBackups: 5
  maxAge: 30
  compress: true

api:
  enableSwagger: true
  swaggerPath: "/swagger/"
  version: "1.0.0"
  title: "SysManix API"
  description: "Cross-platform service management API"
```

## Troubleshooting Configuration Issues

If you're experiencing issues with your configuration:

1. **Validate YAML syntax**
   ```bash
   yamllint config.yaml
   ```

2. **Check permissions**
   ```bash
   # Linux
   ls -la config.yaml

   # Windows
   icacls "C:\path\to\config.yaml"
   ```

3. **Look for parsing errors** in the logs right after startup

4. **Try starting with a minimal configuration** to identify problematic settings

For more help, see the [Troubleshooting Guide](./TROUBLESHOOTING.md).
