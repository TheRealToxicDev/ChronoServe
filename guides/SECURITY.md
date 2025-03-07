# SysManix Security Guide

This guide outlines security best practices for deploying and operating SysManix in production environments.

## Security Architecture

SysManix is designed with security in mind, implementing several layers of protection:

1. **Authentication**: JWT-based authentication with Argon2id password hashing
2. **Authorization**: Role-based access control for all operations
3. **Protection**: Critical system services are protected from modification
4. **Logging**: Comprehensive audit logging of all operations
5. **Input Validation**: Strict validation of all user inputs

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

### Password Storage

SysManix uses Argon2id for password hashing, which is:
- Memory-hard, making hardware attacks expensive
- Configurable for time/memory trade-offs
- Resistant to side-channel attacks
- Winner of the Password Hashing Competition

Default Argon2id parameters in SysManix:
- Memory: 64 MB
- Iterations: 1
- Parallelism: 4
- Key length: 32 bytes

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

### JWT Token Security

To enhance JWT security:

1. **Secret Key**: Use a strong, randomly generated secret key:
   ```bash
   # Linux
   openssl rand -base64 64

   # Windows PowerShell
   [Convert]::ToBase64String((New-Object byte[] 64) | ForEach-Object { $_ = Get-Random -Minimum 0 -Maximum 256 })
   ```

2. **Token Lifetime**: Configure an appropriate token lifetime based on security requirements:
   ```yaml
   auth:
     tokenDuration: 8h  # 8 hours is a reasonable default
   ```

3. **Secret Rotation**: Regularly rotate the JWT secret key in production environments

4. **Token Revocation**: SysManix implements token revocation capabilities to invalidate tokens when needed

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

SysManix should always be deployed behind a TLS-enabled reverse proxy (like Nginx) in production. Configure TLS with modern, secure settings:

```nginx
server {
    # Modern SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_session_tickets off;

    # HSTS (optional, but recommended - 6 months)
    add_header Strict-Transport-Security "max-age=15768000; includeSubDomains" always;
}
```

See the [NGINX Setup Guide](./NGINX_SETUP.md) for complete configuration.

### Firewall Configuration

Restrict access to SysManix API port at the network level:

#### Linux Firewall Rules

```bash
# UFW (Ubuntu/Debian)
sudo ufw allow from 10.0.0.0/8 to any port 40200 proto tcp
sudo ufw allow from 192.168.0.0/16 to any port 40200 proto tcp

# firewalld (CentOS/RHEL/Fedora)
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/8" port port="40200" protocol="tcp" accept'
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.0.0/16" port port="40200" protocol="tcp" accept'
sudo firewall-cmd --reload
```

#### Windows Firewall Rules

```powershell
# Allow specific IP ranges
New-NetFirewallRule -DisplayName "SysManix API - Internal" `
                   -Direction Inbound `
                   -Action Allow `
                   -Protocol TCP `
                   -LocalPort 40200 `
                   -RemoteAddress 10.0.0.0/8,192.168.0.0/16
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

## System Security

### File Permissions

Set restrictive permissions on configuration files:

```bash
# Linux
sudo chmod 640 /etc/sysmanix/config.yaml
sudo chown root:sysmanix /etc/sysmanix/config.yaml
```

```powershell
# Windows
$acl = Get-Acl "C:\Program Files\SysManix\config.yaml"
$acl.SetAccessRuleProtection($true, $false)  # Disable inheritance
$adminRule = New-Object System.Security.AccessControl.FileSystemAccessRule("Administrators", "FullControl", "Allow")
$systemRule = New-Object System.Security.AccessControl.FileSystemAccessRule("SYSTEM", "FullControl", "Allow")
$acl.AddAccessRule($adminRule)
$acl.AddAccessRule($systemRule)
Set-Acl "C:\Program Files\SysManix\config.yaml" $acl
```

### Protected System Services

SysManix includes built-in protection for critical system services to prevent accidental modification that could compromise system stability or security.

You can extend the protected services list in the configuration:

```yaml
windows:
  services:
    protected:
      - wininit
      - csrss
      - lsass
      # Add additional services as needed

linux:
  services:
    protected:
      - systemd
      - systemd-journald
      - dbus
      - sshd
      # Add additional services as needed
```

### Principle of Least Privilege

Run SysManix with the minimum required privileges:

- On **Windows**: Run as a service with a dedicated service account rather than Local System
- On **Linux**: Use systemd security features to limit permissions (see [Systemd Setup](./SYSTEMD_SETUP.md))

## Secure Deployment

### Containers

When deploying in containers, follow these security practices:

1. **Use non-root user**: Add a user and group in the Dockerfile:
   ```dockerfile
   RUN groupadd -r sysmanix && useradd -r -g sysmanix sysmanix
   USER sysmanix
   ```

2. **Multi-stage builds**: Use multi-stage builds to minimize image size:
   ```dockerfile
   FROM golang:1.23-alpine AS builder
   # Build the application

   FROM alpine:latest
   COPY --from=builder /app/sysmanix /usr/local/bin/
   ```

3. **Read-only filesystem**: Mount the file system as read-only and provide a writable volume for logs:
   ```yaml
   # docker-compose.yml
   services:
     sysmanix:
       read_only: true
       volumes:
         - ./logs:/var/log/sysmanix:rw
         - ./config:/etc/sysmanix:ro
   ```

4. **Scan images**: Scan container images for vulnerabilities:
   ```bash
   # Using Trivy
   trivy image sysmanix:latest
   ```

### API Security

1. **Rate limiting**: Implement rate limiting to prevent abuse:
   ```nginx
   # In Nginx configuration
   limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
   limit_req zone=api burst=20 nodelay;
   ```

2. **Input validation**: SysManix validates all inputs to prevent injection attacks
   - Service names are validated for format and length
   - JSON payloads are strictly validated against expected schemas
   - Query parameters are sanitized and type-checked

3. **Output encoding**: All response data is properly encoded to prevent XSS and other injection attacks

4. **Sanitized error messages**: Error messages don't leak sensitive information

## Audit and Logging

### Security Logging

SysManix maintains comprehensive security logs:

1. **Authentication events**:
   - Login attempts (successful and failed)
   - Token creation, revocation, and refresh
   - Permission violations

2. **Service operations**:
   - All start/stop operations with the username
   - Attempted operations on protected services

3. **Configuration changes**:
   - Password changes
   - Configuration file modifications

### Log Security

Protect your log files:

1. **Secure storage**: Set appropriate permissions on log files:
   ```bash
   # Linux
   sudo chmod 640 /var/log/sysmanix/*
   sudo chown root:sysmanix /var/log/sysmanix/*
   ```

2. **Log forwarding**: Forward logs to a secure central logging system
   ```bash
   # Example rsyslog configuration
   if $programname == 'sysmanix' then @logserver.example.com:514
   ```

3. **Log rotation**: Implement proper log rotation to prevent disk space issues
   ```yaml
   # In SysManix config.yaml
   logging:
     maxSize: 10
     maxBackups: 5
     maxAge: 30
     compress: true
   ```

## Vulnerability Management

### Security Updates

1. **Keep SysManix updated**: Check regularly for new versions:
   ```bash
   # Using the built-in version checker
   curl http://localhost:40200/health
   ```

2. **Update dependencies**: SysManix's dependencies are regularly audited and updated

3. **Operating system updates**: Keep the host system updated with security patches

### Security Reporting

If you discover a security vulnerability in SysManix:

1. **DO NOT** disclose it publicly on GitHub issues
2. Email security details to [security@example.com](mailto:security@example.com)
3. Provide detailed information about the vulnerability and steps to reproduce

## Production Hardening Checklist

Use this checklist when deploying SysManix in production:

- [ ] Changed default `secretKey` to a strong random value
- [ ] Changed default admin and viewer passwords
- [ ] Deployed behind TLS-enabled reverse proxy
- [ ] Limited API access with firewall rules
- [ ] Set appropriate file permissions on config and logs
- [ ] Configured proper log rotation
- [ ] Added necessary services to the protected services list
- [ ] Implemented monitoring for failed authentication attempts
- [ ] Set up regular backups of configuration
- [ ] Created a security incident response plan

## Security Monitoring

Monitor the following security metrics:

1. **Failed authentication attempts**:
   - Watch for brute force attacks
   - Implement escalating timeouts or IP bans

2. **Unusual API access patterns**:
   - High volume of requests
   - Requests outside normal business hours
   - Access attempts to protected services

3. **Error rates**:
   - Sudden increases in errors may indicate attacks
   - Monitor 403/401 error rates separately

4. **System resource usage**:
   - Unusual CPU/memory usage
   - Unexpected network traffic

## Further Security Resources

- [OWASP API Security Top 10](https://owasp.org/API-Security/editions/2023/en/0x00-introduction/)
- [JWT Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-jwt-bcp-07)
- [Secure Configuration Guide for Nginx](https://www.nginx.com/resources/wiki/start/topics/tutorials/config_pitfalls/)
- [Linux Security Basics](https://linuxsecurity.expert/security-basics/linux-security-basics)
- [Windows Security Fundamentals](https://docs.microsoft.com/en-us/windows/security/)
