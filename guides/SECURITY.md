# SysManix Security Guide

This guide covers security best practices for deploying and maintaining SysManix in production environments.

## Security Architecture

SysManix was designed with security as a foundational principle:

- **Authentication**: JWT-based authentication with Argon2id password hashing
- **Authorization**: Role-based access control for endpoint protection
- **Service Protection**: Built-in safeguards for critical system services
- **Secure Defaults**: Conservative default settings that prioritize security
- **Audit Logging**: Comprehensive logging of security-relevant events

## Secure Installation

### Initial Setup Security

During installation and initial setup:

1. **Change Default Credentials**: Always change the default admin and viewer passwords
2. **Generate Strong JWT Secret**: Use a cryptographically secure random secret key
3. **Secure Configuration Files**: Apply appropriate file permissions to config.yaml
4. **Setup Firewall Rules**: Restrict network access to the SysManix API port
5. **Use HTTPS**: Deploy with TLS/SSL encryption using Nginx or another reverse proxy

### Permissions & File Security

#### Windows
```powershell
# Secure the configuration file
$acl = Get-Acl "C:\Program Files\SysManix\config.yaml"
$acl.SetAccessRuleProtection($true, $false)  # Disable inheritance
$adminRule = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrators", "FullControl", "Allow")
$systemRule = New-Object System.Security.AccessControl.FileSystemAccessRule("SYSTEM", "FullControl", "Allow")
$acl.AddAccessRule($adminRule)
$acl.AddAccessRule($systemRule)
Set-Acl "C:\Program Files\SysManix\config.yaml" $acl
```

#### Linux
```bash
# Secure the configuration file
sudo chown root:sysmanix /etc/sysmanix/config.yaml
sudo chmod 640 /etc/sysmanix/config.yaml
```

## Authentication Security

### Password Policy

Implement a strong password policy:

1. **Minimum Length**: 12 characters
2. **Complexity**: Require uppercase, lowercase, numbers, and special characters
3. **Rotation Policy**: Regularly update passwords (but avoid too-frequent rotations)
4. **No Common Passwords**: Block known weak or compromised passwords

### JWT Configuration

Secure your JWT implementation:

```yaml
auth:
  secretKey: "[generate-random-64-character-string]"
  tokenDuration: 8h   # Reduce from 24h default for production
  issuedBy: "SysManix-Prod"
```

Best practices for JWT:

1. **Token Expiration**: Set appropriate token lifetimes based on your security requirements
2. **Rotation**: Periodically rotate the JWT secret key
3. **Validation**: Always validate token signature, expiration, and audience
4. **Transport Security**: Only transmit tokens over HTTPS

## Network Security

### Firewall Configuration

Restrict access to the SysManix API:

#### Windows Firewall
```powershell
New-NetFirewallRule -DisplayName "SysManix API" `
                   -Direction Inbound `
                   -Protocol TCP `
                   -LocalPort 40200 `
                   -RemoteAddress 10.0.0.0/8,192.168.1.0/24 `
                   -Action Allow
```

#### Linux (UFW)
```bash
sudo ufw allow from 10.0.0.0/8 to any port 40200 proto tcp
sudo ufw allow from 192.168.1.0/24 to any port 40200 proto tcp
```

### TLS/SSL Configuration

When using Nginx as a TLS termination proxy, follow these best practices:

```nginx
# Modern SSL configuration
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305;
ssl_prefer_server_ciphers off;
ssl_session_timeout 1d;
ssl_session_cache shared:SSL:10m;
ssl_session_tickets off;

# OCSP Stapling
ssl_stapling on;
ssl_stapling_verify on;

# Security headers
add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";
add_header X-Content-Type-Options nosniff;
add_header X-Frame-Options DENY;
add_header X-XSS-Protection "1; mode=block";
```

## Role-Based Access Control

### Principle of Least Privilege

Follow the principle of least privilege when assigning roles:

1. Use `viewer` role for users who only need to view service status
2. Reserve `admin` role for users who need to start/stop services
3. Consider creating custom roles for specific needs

### Implementation Example

```yaml
auth:
    allowedRoles:
        - admin
        - viewer
        - service_operator
        - log_viewer
    users:
        admin:
            username: "admin"
            password_hash: "$argon2id$v=19$..."
            roles: ["admin"]
        operator:
            username: "operator"
            password_hash: "$argon2id$v=19$..."
            roles: ["service_operator", "log_viewer"]
        logviewer:
            username: "logviewer"
            password_hash: "$argon2id$v=19$..."
            roles: ["log_viewer"]
```

## Protecting Sensitive Services

SysManix includes built-in protection for system-critical services. Extend this protection to additional services as needed:

### Windows Protected Services Configuration

```yaml
windows:
  services:
    protected:
      - wininit
      - csrss
      - lsass
      - spooler
      - EventLog
      - TrustedInstaller
      - YourCriticalService
```

### Linux Protected Services Configuration

```yaml
linux:
  services:
    protected:
      - systemd
      - systemd-journald
      - dbus
      - sshd
      - NetworkManager
      - your-critical-service
```

## Access Control Integration

### LDAP/Active Directory Integration

For enterprise environments, consider integrating with your directory service:

```yaml
auth:
  providers:
    ldap:
      enabled: true
      server: "ldap://dc.example.com"
      bindDN: "CN=sysmanix,OU=ServiceAccounts,DC=example,DC=com"
      bindPassword: "secure-password"
      userSearchBase: "OU=Users,DC=example,DC=com"
      userSearchFilter: "(&(objectClass=user)(sAMAccountName=%s))"
      groupSearchBase: "OU=Groups,DC=example,DC=com"
      groupSearchFilter: "(&(objectClass=group)(member=%s))"
      roleMapping:
        "SysManix Admins": ["admin"]
        "SysManix Viewers": ["viewer"]
```

## Audit Logging

Configure comprehensive security audit logging:

```yaml
logging:
  level: info
  directory: "/var/log/sysmanix"
  maxSize: 10
  maxBackups: 30
  maxAge: 90
  compress: true
  auditLogs:
    enabled: true
    file: "audit.log"
    logAuthAttempts: true
    logServiceOperations: true
```

### Log Review Process

Implement a regular log review process:

1. Check authentication failures for potential brute-force attempts
2. Monitor service start/stop operations for unauthorized activity
3. Review role permission changes
4. Set up alerting for suspicious patterns

## Security Monitoring

### Health Check Monitoring

Regularly monitor the health endpoint:

```powershell
# PowerShell script for monitoring
$response = Invoke-RestMethod -Uri "https://sysmanix.example.com/health" -Headers @{
    "Authorization" = "Bearer $token"
}

if ($response.data.status -ne "healthy") {
    # Send alert
    Send-MailMessage -To "admin@example.com" -Subject "SysManix Health Check Failed"
}
```

### Failed Login Alerts

Set up alerts for multiple failed login attempts:

```bash
# Check for failed logins in the last hour
failed_attempts=$(grep "Failed login attempt" /var/log/sysmanix/auth.log | grep -c "$(date -d '1 hour ago' '+%Y-%m-%d %H:')")

if [ $failed_attempts -gt 5 ]; then
    # Send alert
    mail -s "SysManix Security Alert: Multiple Failed Logins" admin@example.com
fi
```

## Software Updates

### Update Checking

Enable automatic update checking:

```yaml
updates:
  checkOnStartup: true
  notifyInLogs: true
  checkInterval: 24h
  githubTimeout: 10s
```

### Update Procedure

When a new version is available:

1. Review the release notes for security fixes
2. Test the update in a non-production environment
3. Back up the configuration and database
4. Apply the update during a scheduled maintenance window
5. Verify all functionality after updating

## Secure Development Practices

If you're contributing to SysManix or extending it:

1. **Code Reviews**: All changes should undergo security-focused code review
2. **Dependency Scanning**: Regularly scan dependencies for vulnerabilities
3. **Static Analysis**: Use static analysis tools to identify security issues
4. **Input Validation**: Always validate and sanitize all user inputs
5. **Output Encoding**: Properly encode output to prevent injection attacks

## Disaster Recovery

### Backup Configuration

Regularly back up your SysManix configuration:

```bash
# Linux backup script
mkdir -p /backups/sysmanix/$(date +%Y-%m-%d)
cp -r /etc/sysmanix/* /backups/sysmanix/$(date +%Y-%m-%d)/
```

```powershell
# Windows backup script
$date = Get-Date -Format "yyyy-MM-dd"
$backupDir = "C:\Backups\SysManix\$date"
New-Item -ItemType Directory -Force -Path $backupDir
Copy-Item -Path "C:\Program Files\SysManix\config.yaml" -Destination $backupDir
```

### Recovery Plan

Create a documented recovery plan:

1. Reinstall SysManix using the installation guide
2. Restore configuration from backup
3. Verify service operations
4. Reset service states as necessary
5. Update documentation with lessons learned

## Compliance Considerations

### Audit Trail Requirements

For environments with compliance requirements:

1. Enable detailed audit logging
2. Configure log retention to meet compliance requirements
3. Implement log forwarding to a central SIEM system
4. Set up access controls for log files

### Data Privacy

Ensure PII (Personally Identifiable Information) is protected:

1. Don't store sensitive data in logs
2. Sanitize error messages that might contain sensitive data
3. Ensure logs are stored securely with appropriate retention policies

## Further Reading

- [OWASP Security Cheat Sheet](https://cheatsheetseries.owasp.org/)
- [JWT Security Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)
- [Linux Service Hardening](https://www.cisecurity.org/benchmark/distribution_independent_linux/)
- [Windows Service Hardening](https://www.cisecurity.org/benchmark/microsoft_windows_server/)
- [NestJS Security Guide](https://docs.nestjs.com/techniques/security)
