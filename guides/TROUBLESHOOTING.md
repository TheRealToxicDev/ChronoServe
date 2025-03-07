# SysManix Troubleshooting Guide

This guide provides solutions to common issues encountered while using SysManix.

## Table of Contents
- [Common Issues and Solutions](#common-issues-and-solutions)
  - [Installation Issues](#installation-issues)
    - [Unable to start SysManix](#unable-to-start-sysmanix)
    - [Permission Denied Errors](#permission-denied-errors)
  - [Authentication Issues](#authentication-issues)
    - ["Invalid credentials" error](#invalid-credentials-error)
    - ["Token expired" or "Invalid token" errors](#token-expired-or-invalid-token-errors)
  - [Service Management Issues](#service-management-issues)
    - [Service operations timing out](#service-operations-timing-out)
    - ["Protected system service" errors](#protected-system-service-errors)
  - [Configuration Issues](#configuration-issues)
    - [Configuration not being applied](#configuration-not-being-applied)
    - [Default credentials security warning](#default-credentials-security-warning)
  - [Platform-Specific Issues](#platform-specific-issues)
    - [Windows Service Issues](#windows-service-issues)
    - [Linux Service Issues](#linux-service-issues)
  - [Network Issues](#network-issues)
    - [API inaccessible from remote hosts](#api-inaccessible-from-remote-hosts)
    - [CORS-related errors](#cors-related-errors)
  - [Performance Issues](#performance-issues)
    - [Slow service listing](#slow-service-listing)
    - [High CPU or memory usage](#high-cpu-or-memory-usage)
  - [Logging and Diagnostics](#logging-and-diagnostics)
    - [Insufficient logging information](#insufficient-logging-information)
    - [Log files growing too large](#log-files-growing-too-large)
- [Diagnostic Procedures](#diagnostic-procedures)
  - [Health Check Endpoint](#health-check-endpoint)
  - [API Testing](#api-testing)
  - [Configuration Verification](#configuration-verification)
- [Getting Help](#getting-help)
- [Common Error Codes](#common-error-codes)

## Common Issues and Solutions

### Installation Issues

#### Unable to start SysManix

**Symptoms:**
- Application crashes immediately after starting
- Error message about port already in use

**Potential Solutions:**
1. Check if another application is using port 40200
   ```bash
   # Linux
   sudo netstat -tulpn | grep 40200

   # Windows (PowerShell)
   netstat -ano | findstr 40200
   ```
2. Change the port in `config.yaml`:
   ```yaml
   server:
     port: 8080  # Change to an available port
   ```

#### Permission Denied Errors

**Symptoms:**
- Error messages containing "permission denied" or "access denied"
- Unable to start/stop services

**Potential Solutions:**
1. Ensure SysManix is running with administrative privileges
   ```bash
   # Linux
   sudo ./sysmanix

   # Windows (PowerShell as Administrator)
   Start-Process .\SysManix.exe -Verb RunAs
   ```
2. Check file permissions on the config file
   ```bash
   # Linux
   sudo chmod 640 /etc/sysmanix/config.yaml
   sudo chown root:sysmanix /etc/sysmanix/config.yaml
   ```

### Authentication Issues

#### "Invalid credentials" error

**Symptoms:**
- Login attempts fail with "Invalid credentials" message
- Unable to obtain JWT token

**Potential Solutions:**
1. Verify username and password in your request
2. Check if the user exists in the config.yaml file
3. Reset the password by editing config.yaml:
   ```yaml
   auth:
     users:
       admin:
         username: "admin"
         password: "new-password"  # Will be hashed after restart
         # password_hash will be regenerated
   ```
4. Restart SysManix to hash the new password

#### "Token expired" or "Invalid token" errors

**Symptoms:**
- API requests fail with 401 Unauthorized
- Error messages about expired or invalid tokens

**Potential Solutions:**
1. Obtain a new token through the login endpoint
2. Use the token refresh endpoint:
   ```bash
   curl -X POST http://localhost:40200/auth/tokens/refresh \
     -H "Authorization: Bearer YOUR_EXPIRED_TOKEN"
   ```
3. Increase token duration in config.yaml (default is 24h)
   ```yaml
   auth:
     tokenDuration: 72h  # 3 days
   ```

### Service Management Issues

#### Service operations timing out

**Symptoms:**
- Start/stop operations taking too long
- 408 Request Timeout errors

**Potential Solutions:**
1. Check service status in OS directly to see if it's responsive
   ```bash
   # Linux
   systemctl status service-name

   # Windows (PowerShell)
   Get-Service service-name
   ```
2. Increase operation timeout in code or config
3. Restart the problematic service from OS native tools
4. Check service logs for underlying issues:
   ```bash
   # Linux
   journalctl -u service-name

   # Windows (PowerShell)
   Get-EventLog -LogName System | Where-Object {$_.Source -eq "service-name"}
   ```

#### "Protected system service" errors

**Symptoms:**
- 403 Forbidden responses
- Error message about protected services

**Explanation:**
- SysManix prevents operations on critical system services by design
- This is a safety feature to prevent system instability

**Potential Solutions:**
1. Use a different, non-protected service
2. If necessary, modify the protected services list in your code:
   ```go
   // In windows/critical.go or linux/critical.go
   // Modify the protection list carefully and at your own risk
   ```

### Configuration Issues

#### Configuration not being applied

**Symptoms:**
- Changes to config.yaml don't seem to take effect
- Default values used instead of configured ones

**Potential Solutions:**
1. Verify the config file path is correct
2. Make sure the yaml syntax is valid
3. Restart SysManix after configuration changes
4. Check logs for configuration parsing errors
5. Try running with explicit config path:
   ```bash
   ./sysmanix -config /path/to/config.yaml
   ```

#### Default credentials security warning

**Symptoms:**
- Warning about default credentials on startup
- Security risks reported in logs

**Potential Solutions:**
1. Change the default secret key:
   ```yaml
   auth:
     secretKey: "generate-a-long-random-string-here"
   ```
2. Update default user passwords
3. Restart SysManix to apply changes

### Platform-Specific Issues

#### Windows Service Issues

**Symptoms:**
- Unable to retrieve Windows service information
- PowerShell execution errors

**Potential Solutions:**
1. Ensure PowerShell execution policy allows script execution:
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope LocalMachine
   ```
2. Run SysManix as Administrator
3. Check Windows Event Viewer for PowerShell errors
4. Verify service names match exactly (case-sensitive)

#### Linux Service Issues

**Symptoms:**
- Unable to control systemd services
- Permission errors when accessing systemctl

**Potential Solutions:**
1. Run SysManix as root or with sudo
2. Ensure systemd is available and running:
   ```bash
   systemctl --version
   ```
3. Check if user has proper sudo permissions for systemctl
4. Try adding user to systemd-journal group:
   ```bash
   sudo usermod -a -G systemd-journal username
   ```

### Network Issues

#### API inaccessible from remote hosts

**Symptoms:**
- Can access API locally but not from other machines
- Connection refused or timeout errors

**Potential Solutions:**
1. Check host configuration in config.yaml:
   ```yaml
   server:
     host: "0.0.0.0"  # Listen on all interfaces, not just localhost
   ```
2. Verify firewall rules allow access to the configured port
3. Check if any security software is blocking connections
4. Verify network routing between client and server

#### CORS-related errors

**Symptoms:**
- Browser console shows CORS errors
- Web clients can't access API

**Potential Solutions:**
1. Verify CORS middleware is properly configured
2. Ensure API responses include proper CORS headers
3. For testing, try a browser extension that disables CORS restrictions

### Performance Issues

#### Slow service listing

**Symptoms:**
- `/services` endpoint takes a long time to respond
- Timeouts when listing many services

**Potential Solutions:**
1. Implement pagination for service listing
2. Reduce timeout threshold in client applications
3. Monitor system resource usage during operations
4. Optimize service status retrieval to use more efficient commands

#### High CPU or memory usage

**Symptoms:**
- SysManix consuming excessive system resources
- System performance degradation

**Potential Solutions:**
1. Check for memory leaks in long-running instances
2. Implement request throttling for busy API endpoints
3. Add more granular logging to identify resource-heavy operations
4. Consider implementing caching for frequent service status requests

### Logging and Diagnostics

#### Insufficient logging information

**Symptoms:**
- Unable to troubleshoot issues due to limited log data
- Missing context in error messages

**Potential Solutions:**
1. Increase log verbosity in config.yaml:
   ```yaml
   logging:
     level: "debug"  # Options: debug, info, warn, error
   ```
2. Enable request/response logging for API endpoints
3. Add specific debug flags for problematic components
4. Implement structured logging for better parsing

#### Log files growing too large

**Symptoms:**
- Disk space filling up quickly
- Log files becoming unmanageable

**Potential Solutions:**
1. Configure log rotation in config.yaml:
   ```yaml
   logging:
     maxSize: 10    # MB per file
     maxBackups: 5  # Number of rotated files to keep
     compress: true # Compress old log files
   ```
2. Reduce logging verbosity in production
3. Implement external log management (e.g., syslog, ELK stack)

### Configuration File Not Found

Ensure the `config.yaml` file is in the correct location and properly formatted.

### Incorrect Credentials

Verify the credentials in your `config.yaml` file:

```yaml
auth:
  users:
    admin:
      username: "admin"
      password_hash: "$argon2id$v=19$m=65536,t=1,p=4$..."
      roles:
        - admin
```

### Service Control Fails

Ensure SysManix is running with the necessary permissions to manage services.

### Logs Not Appearing

Check the log directory specified in your `config.yaml` file:

```yaml
logging:
  directory: "/var/log/sysmanix"
```

## Viewing Logs

### Linux

```bash
sudo journalctl -u sysmanix
```

### Windows

```powershell
Get-EventLog -LogName Application | Where-Object {$_.Source -eq "SysManix"}
```

## Diagnostic Procedures

### Health Check Endpoint

Use the health check endpoint to verify system status:

```bash
curl http://localhost:40200/health
```

The response will include:
- Current version
- Uptime
- Memory usage
- Go runtime information

### API Testing

Use curl or Postman to test API endpoints independently:

```bash
# Test authentication
curl -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'

# Test service listing
curl -X GET http://localhost:40200/services \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Configuration Verification

Validate your configuration file syntax:

```bash
# Using yaml-lint (install first if needed)
yaml-lint config.yaml

# Or online at https://www.yamllint.com/
```

## Getting Help

If you're still experiencing issues after trying the solutions in this guide:

1. Check the [GitHub Issues](https://github.com/toxic-development/sysmanix/issues) for similar problems and solutions
2. Search the project documentation for more specific guidance
3. Open a new GitHub issue with:
   - Detailed description of the problem
   - Steps to reproduce
   - Relevant logs and configuration (with sensitive data redacted)
   - Operating system and environment details
   - SysManix version information

## Common Error Codes

| Error Code | Description | Possible Causes |
|------------|-------------|----------------|
| 400 | Bad Request | Invalid request body, missing parameters |
| 401 | Unauthorized | Missing or invalid token, expired token |
| 403 | Forbidden | Insufficient permissions, protected service |
| 404 | Not Found | Service doesn't exist, invalid endpoint |
| 408 | Request Timeout | Service operation took too long |
| 500 | Internal Server Error | Server-side errors, system issues |
| 503 | Service Unavailable | API temporarily overloaded or down |
