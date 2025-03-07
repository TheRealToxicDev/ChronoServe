# SysManix Troubleshooting Guide

## Table of Contents
- [Common Issues](#common-issues)
  - [Authentication Issues](#authentication-issues)
    - [Invalid Credentials Error](#invalid-credentials-error)
    - [Token Validation Failed](#token-validation-failed)
  - [Service Access Issues](#service-access-issues)
    - [Permission Denied](#permission-denied)
    - [Service Not Found](#service-not-found)
    - [Protected System Service](#protected-system-service)
    - [Service Operation Timeout](#service-operation-timeout)
  - [Configuration Issues](#configuration-issues)
    - [Default Security Values](#default-security-values)
    - [Port Already in Use](#port-already-in-use)
  - [Logging Issues](#logging-issues)
    - [Missing Log Files](#missing-log-files)
- [Debugging Tips](#debugging-tips)
  - [Enable Debug Logging](#enable-debug-logging)
  - [Check Log Files](#check-log-files)
  - [Test Health Endpoint](#test-health-endpoint)
  - [Verify Service Status](#verify-service-status)
- [Common Error Messages](#common-error-messages)
- [Network Troubleshooting](#network-troubleshooting)
  - [Connection Refused](#connection-refused)
  - [CORS Errors](#cors-errors)
- [Performance Issues](#performance-issues)
  - [Slow API Response](#slow-api-response)
  - [High Memory Usage](#high-memory-usage)
- [Installation Problems](#installation-problems)
  - [Missing Dependencies](#missing-dependencies)
  - [Permission Issues](#permission-issues)
- [Getting Help](#getting-help)

## Common Issues

### Authentication Issues

#### Invalid Credentials Error

**Symptom:**
```
{
    "success": false,
    "error": "Invalid credentials"
}
```

**Common Causes:**
1. Username mismatch in config
2. Incorrect map key in users configuration
3. Password not matching stored hash

**Solution:**
1. Check your config.yaml structure. The user map key should match the login username:

```yaml
auth:
    users:
        toxicdev:            # This key must match login username
            username: toxicdev
            password_hash: "$argon2id$v=19$..."
            roles: ["admin"]
```

2. If setting up a new password:
   - Add password field to config
   - Remove password_hash
   - Restart application
   - Let it generate new hash

#### Token Validation Failed

**Symptom:**
```
{
    "success": false,
    "error": "Invalid or expired token"
}
```

**Solution:**
1. Check token expiration
2. Re-authenticate to get new token
3. Verify secret key hasn't changed

### Service Access Issues

#### Permission Denied

**Symptom:**
```
{
    "success": false,
    "error": "Insufficient permissions"
}
```

**Solution:**
1. Verify user has correct roles
2. Check role requirements:
   - `admin`: Required for start/stop
   - `viewer`: Can only view logs/status

#### Service Not Found

**Symptom:**
```
"Failed to retrieve logs for service XYZ: exit status 1"
```

**Solution:**
1. List available services:
```powershell
$headers = @{
    Authorization = "Bearer $token"
}
Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

2. Try these common service names:
- `wuauserv` (Windows Update)
- `EventLog` (Event Logging)
- `WinRM` (Remote Management)

#### Protected System Service

**Symptom:**
```
{
    "success": false,
    "error": "operation not allowed on protected system service: wininit",
    "code": 403
}
```

**Solution:**
This is a security feature. Critical system services are protected from modifications to prevent system damage:

1. Windows protected services include:
   - `wininit`, `csrss`, `services`, `lsass`, `winlogon`, and other critical Windows components

2. Linux protected services include:
   - `systemd`, `systemd-journald`, `sshd`, `dbus`, and other core system services

3. Use alternate services that can be safely modified

#### Service Operation Timeout

**Symptom:**
```
"timeout waiting for service to start (status: starting)"
```

**Solution:**
1. Service operation may be in progress but taking longer than the timeout (10 seconds)
2. Check service status manually:
```powershell
# For Windows
Get-Service servicename

# For Linux
systemctl status servicename
```
3. Some services have dependencies that need to start first
4. Try starting the service with elevated privileges directly on the system

### Configuration Issues

#### Default Security Values

**Symptom:**
```
"security risk: default secret key must be changed"
```

**Solution:**
1. Update config.yaml:
```yaml
auth:
    secretKey: "your-secure-key"  # Change this
    users:
        admin:
            password: "new-password"  # Add this, remove password_hash
```

2. Restart application

#### Port Already in Use

**Symptom:**
```
"listen tcp :40200: bind: address already in use"
```

**Solution:**
1. Find process using port:
```powershell
netstat -ano | findstr :40200
```

2. Change port in config:
```yaml
server:
    port: 40201  # Use different port
```

### Logging Issues

#### Missing Log Files

**Symptom:**
Log files not appearing in logs directory

**Solution:**
1. Check permissions
2. Verify log config:
```yaml
logging:
    level: debug
    directory: logs
    maxSize: 10
```

3. Create logs directory manually:
```powershell
mkdir logs
```

## Debugging Tips

### Enable Debug Logging
```yaml
logging:
    level: debug  # Set to debug for more info
```

### Check Log Files
```powershell
# View authentication logs
Get-Content .\logs\auth.log -Tail 20

# View application logs
Get-Content .\logs\app.log -Tail 20
```

### Test Health Endpoint
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/health"
```

### Verify Service Status
```powershell
# PowerShell
Get-Service wuauserv  # Check actual service status

# SysManix
$headers = @{
    Authorization = "Bearer $token"
}
Invoke-RestMethod -Uri "http://localhost:40200/services/status/wuauserv" -Headers $headers
```

## Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "Invalid credentials" | Username/password mismatch | Check config user map key |
| "Invalid token" | Expired/malformed JWT | Re-authenticate |
| "Insufficient permissions" | Missing required role | Check user roles |
| "Service not found" | Invalid service name | Use correct service identifier |
| "Default secret key" | Security risk | Update auth.secretKey |

### Network Troubleshooting

#### Connection Refused

**Symptom:**
```
"Failed to connect to localhost:40200: Connection refused"
```

**Common Causes:**
1. SysManix service not running
2. Incorrect port configuration
3. Firewall blocking connection

**Solution:**
1. Check if service is running:
```powershell
# For Windows
Get-Process -Name SysManix* -ErrorAction SilentlyContinue
```

2. Verify listening port:
```powershell
netstat -ano | findstr :40200
```

3. Check firewall settings:
```powershell
# For Windows
New-NetFirewallRule -DisplayName "Allow SysManix" -Direction Inbound -LocalPort 40200 -Protocol TCP -Action Allow
```

#### CORS Errors

**Symptom:**
Browser console shows:
```
Access to fetch at 'http://localhost:40200/auth' from origin 'http://localhost:3000' has been blocked by CORS policy
```

**Solution:**
1. Update config.yaml to allow your frontend origin:
```yaml
server:
  cors:
    allowed_origins: ["http://localhost:3000"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["Content-Type", "Authorization"]
```

2. Restart the application

### Performance Issues

#### Slow API Response

**Symptom:**
API calls taking more than 1-2 seconds to respond

**Solutions:**
1. Check system resource usage:
```powershell
# Check CPU/Memory
Get-Process -Name SysManix* | Select-Object CPU, WorkingSet, ID
```

2. Review log files for slow operations
3. Enable profiling in config:
```yaml
debug:
  profiling: true
  profiling_endpoint: "/debug/pprof"
```

4. Access profiles at http://localhost:40200/debug/pprof/ (requires admin login)

#### High Memory Usage

**Symptom:**
SysManix process consuming excessive memory

**Solutions:**
1. Check for memory leaks in logs
2. Adjust cache settings in config:
```yaml
cache:
  max_size: 100  # Reduce cache size
  ttl: 300       # Shorter time-to-live (seconds)
```

3. Restart application periodically using a scheduled task

### Installation Problems

#### Missing Dependencies

**Symptom:**
```
"Failed to start: exec: xyz: executable file not found in %PATH%"
```

**Solutions:**
1. Install required dependencies:
   - For Windows: PowerShell 5.0+ and .NET Framework 4.5+
   - For Linux: systemd and required libraries

2. For Docker installations, ensure correct base image:
```dockerfile
# Use correct base image
FROM mcr.microsoft.com/windows/servercore:ltsc2019
# OR
FROM ubuntu:20.04
```

#### Permission Issues

**Symptom:**
```
"Failed to access service control manager: Access is denied"
```

**Solutions:**
1. Run application with administrative privileges
2. Check service user permissions:
   - Windows: Run as administrator
   - Linux: Use sudo or proper systemd service user

## Getting Help

1. Check the logs:
   - auth.log for authentication issues
   - app.log for application errors
   - access.log for request history

2. Documentation:
   - [API Reference](API_REFERENCE.md)
   - [Authentication Guide](AUTHENTICATION.md)
   - [Configuration Guide](DOCUMENTATION.md#configuration)

3. Report Issues:
   - GitHub Issues: [SysManix Repository](https://github.com/toxic-development/SysManix)
   - Include logs and configuration (remove sensitive data)
