# ChronoServe Troubleshooting Guide

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

# ChronoServe
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
   - GitHub Issues: [ChronoServe Repository](https://github.com/therealtoxicdev/chronoserve)
   - Include logs and configuration (remove sensitive data)