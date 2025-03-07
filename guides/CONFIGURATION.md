# SysManix Configuration Guide

This guide explains how to configure SysManix for optimal performance, security, and functionality.

## Configuration File Location

SysManix uses a YAML configuration file that can be located at:

- **Windows**: `C:\Program Files\SysManix\config.yaml` (default installation)
- **Linux**: `/etc/sysmanix/config.yaml`

## Configuration Structure

The configuration file is organized into several sections:

```yaml
server:           # API server settings
auth:             # Authentication settings
logging:          # Logging configuration
windows:          # Windows-specific settings
linux:            # Linux-specific settings
updates:          # Update checking settings
```

## Server Configuration

```yaml
server:
  host: "localhost"      # Interface to bind to (use "0.0.0.0" to listen on all interfaces)
  port: 40200            # Port to listen on
  readTimeout: "15s"     # HTTP read timeout
  writeTimeout: "15s"    # HTTP write timeout
  maxHeaderBytes: 1048576 # Maximum size of request headers (1MB)
```

### Binding to Different Interfaces

- `localhost`: Only accept connections from the local machine
- `0.0.0.0`: Accept connections from any network interface (for remote access)
- Specific IP: Bind to a specific network interface (e.g., "192.168.1.100")

### Changing the Port

If port 40200 is already in use, you can change it:

```yaml
server:
  port: 8080  # Use any available port
```

Remember to update your firewall rules and client configurations if you change the port.

## Authentication Configuration

```yaml
auth:
  secretKey: "your-secure-secret-key"  # JWT signing key (keep this secure!)
  tokenDuration: 24h                   # JWT token validity period
  issuedBy: "SysManix"                 # JWT issuer field
  allowedRoles:                        # Available roles in the system
    - admin
    - viewer
  users:                               # User definitions
    admin:
      username: "admin"
      password_hash: "$argon2id$v=19$..." # Will be generated after first run
      roles:
        - admin
    viewer:
      username: "viewer"
      password_hash: "$argon2id$v=19$..." # Will be generated after first run
      roles:
        - viewer
```

### Setting Up Users

For new installations, you can specify plaintext passwords:

```yaml
auth:
  users:
    admin:
      username: "admin"
      password: "secure-admin-password"  # Will be hashed after first run
      roles:
        - admin
```

After the first run, SysManix will convert plaintext passwords to secure Argon2id hashes.

### Token Duration

Adjust token lifetime based on your security requirements:

```yaml
auth:
  tokenDuration: 8h  # 8 hours
```

Common values:
- Development: `24h` (24 hours)
- Standard security: `8h` (8 hours)
- High security: `1h` (1 hour)

### Custom Roles

You can create custom roles beyond the built-in `admin` and `viewer`:

```yaml
auth:
  allowedRoles:
    - admin
    - viewer
    - service_operator
    - log_viewer
  users:
    operator:
      username: "operator"
      password_hash: "$argon2id$v=19$..."
      roles:
        - service_operator
        - log_viewer
```

## Logging Configuration

```yaml
logging:
  level: "info"        # Log level (debug, info, warn, error)
  directory: "logs"    # Directory for log files
  maxSize: 10          # Maximum size in MB before rotation
  maxBackups: 5        # Number of rotated logs to keep
  maxAge: 30           # Days to keep logs
  compress: true       # Compress rotated logs
```

### Log Levels

Available log levels (from most to least verbose):
- `debug`: Detailed debugging information
- `info`: General operational information
- `warn`: Warning conditions
- `error`: Error conditions that should be addressed

### Log File Management

For high-volume environments, adjust:

```yaml
logging:
  maxSize: 50        # 50MB per file before rotation
  maxBackups: 10     # Keep more backups
  maxAge: 90         # Keep logs longer (90 days)
```

### Log Directory

Specify an absolute path for better control:

```yaml
logging:
  directory: "/var/log/sysmanix"  # Linux
  # or
  directory: "C:\\ProgramData\\SysManix\\logs"  # Windows
```

## Platform-Specific Configuration

### Windows Configuration

```yaml
windows:
  serviceCommand: "sc"  # Windows service command tool
  logDirectory: "C:\\ProgramData\\SysManix\\logs"
  services:
    protected:          # Services that cannot be modified
      - wininit
      - csrss
      - lsass
      - spooler
      - EventLog
      - TrustedInstaller
```

### Linux Configuration

```yaml
linux:
  serviceCommand: "systemctl"  # Systemd control command
  logDirectory: "/var/log/sysmanix"
  services:
    protected:          # Services that cannot be modified
      - systemd
      - systemd-journald
      - dbus
      - sshd
      - NetworkManager
```

## Update Configuration

```yaml
updates:
  checkOnStartup: true    # Check for updates when server starts
  notifyInLogs: true      # Log update notifications
  checkInterval: 24h      # How often to check for updates
  githubTimeout: 10s      # Timeout for GitHub API requests
```

### Disabling Update Checks

For air-gapped environments:

```yaml
updates:
  checkOnStartup: false
  notifyInLogs: false
```

## Advanced Configuration

### Debug Features

For troubleshooting, you can enable debugging features:

```yaml
debug:
  profiling: true
  profiling_endpoint: "/debug/pprof"
  trace: true
  memory_profile: true
```

These features should be disabled in production environments.

### CORS Configuration

Configure Cross-Origin Resource Sharing for browser-based clients:

```yaml
server:
  cors:
    enabled: true
    allowed_origins:
      - "https://admin.example.com"
    allowed_methods:
      - "GET"
      - "POST"
    allowed_headers:
      - "Authorization"
      - "Content-Type"
    max_age: 86400  # 24 hours
```

### Rate Limiting

Protect against abuse with rate limiting:

```yaml
server:
  rate_limiting:
    enabled: true
    requests_per_minute: 60
    burst: 10
```

## Configuration Examples

### Minimal Production Configuration

```yaml
server:
  host: "0.0.0.0"
  port: 40200

auth:
  secretKey: "your-secure-random-string"
  tokenDuration: 8h
  users:
    admin:
      username: "admin"
      password_hash: "$argon2id$v=19$..."
      roles:
        - admin

logging:
  level: "info"
  directory: "/var/log/sysmanix"
```

### High-Security Configuration

```yaml
server:
  host: "localhost"  # Only local access
  port: 40200
  readTimeout: "5s"
  writeTimeout: "5s"

auth:
  secretKey: "very-long-random-string-at-least-64-characters-for-enhanced-security"
  tokenDuration: 1h  # Short-lived tokens
  users:
    admin:
      username: "admin"
      password_hash: "$argon2id$v=19$..."
      roles:
        - admin

logging:
  level: "info"
  directory: "/var/log/sysmanix"
  maxBackups: 30
  maxAge: 90

linux:
  services:
    protected:
      - systemd
      - systemd-journald
      - dbus
      - sshd
      - NetworkManager
      - firewalld
      - ufw
```

### Development Configuration

```yaml
server:
  host: "localhost"
  port: 40200

auth:
  secretKey: "dev-secret-do-not-use-in-production"
  tokenDuration: 24h
  users:
    admin:
      username: "admin"
      password: "admin"  # Will be hashed after first run
      roles:
        - admin
    viewer:
      username: "viewer"
      password: "viewer"  # Will be hashed after first run
      roles:
        - viewer

logging:
  level: "debug"
  directory: "logs"

debug:
  profiling: true
  profiling_endpoint: "/debug/pprof"
```

## Configuration Validation

SysManix validates your configuration at startup. If there are issues, you'll see errors in the logs.

Common validation errors:
- Missing required fields
- Invalid values (e.g., invalid duration format)
- Security issues (e.g., weak default passwords in production mode)

## Environment Variable Overrides

You can override configuration settings with environment variables:

```bash
# Linux/macOS
export SYSMANIX_SERVER_PORT=8080
export SYSMANIX_AUTH_TOKEN_DURATION=12h

# Windows PowerShell
$env:SYSMANIX_SERVER_PORT = 8080
$env:SYSMANIX_AUTH_TOKEN_DURATION = "12h"
```

## Reloading Configuration

After changing the configuration file, restart SysManix to apply changes:

### Windows
```powershell
Restart-Service -Name SysManix
```

### Linux
```bash
sudo systemctl restart sysmanix
```

Some settings (if implemented) may be reloadable without a restart:

```bash
# Send SIGHUP to reload configuration (Linux)
sudo kill -HUP $(pidof sysmanix)
```

## Troubleshooting

### Invalid Configuration

If SysManix fails to start due to configuration issues:

1. Check log files for specific error messages
2. Verify YAML syntax (no tabs, proper indentation)
3. Validate file permissions
4. Try the minimal configuration example above

### Password Issues

If you're unable to log in:

1. Ensure the username in the request matches the configuration exactly (case-sensitive)
2. For fresh installs, use plaintext passwords and let SysManix hash them
3. To reset a password, you can temporarily replace a password_hash with a plaintext password field

## Further Reading

- [Authentication Guide](./AUTHENTICATION.md): More details on authentication options
- [Security Guide](./SECURITY.md): Security best practices
- [Troubleshooting Guide](./TROUBLESHOOTING.md): Solving common problems
- [Advanced Configuration](./ADVANCED_CONFIG.md): Advanced configuration options
