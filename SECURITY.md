# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously at SysManix. Please follow these steps to report security issues:

1. **DO NOT** open public issues for security vulnerabilities
2. Report issues - [here](https://github.com/toxic-development/sysmanix/issues)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if available)

## Security Best Practices

### Production Deployment

1. **Always use HTTPS** in production
2. Change all default credentials
3. Use strong passwords
4. Keep the secret key secure
5. Run with minimal required privileges

### Configuration Security

1. Store sensitive config values as environment variables:
```bash
CHRONOSERVE_SECRET_KEY=your-secure-key
CHRONOSERVE_ADMIN_PASSWORD=secure-admin-password
```

2. Use secure file permissions:
```powershell
# Windows (PowerShell)
icacls .\config.yaml /inheritance:r /grant:r "SYSTEM:(R)" "ADMINISTRATORS:(F)"
```

```bash
# Linux
chmod 600 config.yaml
chown root:root config.yaml
```

## Security Features

### Authentication
- JWT-based authentication with configurable expiration
- Argon2id password hashing (memory-hard algorithm)
- Automatic removal of plain text passwords
- Constant-time password comparison

### Access Control
- Role-based access control (RBAC)
- Granular permission system
- Limited service access based on roles
- Token validation on every request
- Protected service list to prevent critical system modifications

### Password Security
- Enforced password complexity
- Secure password storage
- No plain text passwords in logs
- Password validation requirements:
  - Minimum length: 12 characters
  - Mixed case letters
  - Numbers
  - Special characters

### Configuration Protection
- Secure default settings
- Required secret key change
- Protected config file access
- Sensitive value validation

### Logging Security
- No sensitive data in logs
- Rotated log files
- Compressed archived logs
- Configurable log retention

### Service Security
- Minimal required permissions
- Sanitized command input
- Protected service operations
- Validation of service names

## Security Checklist

### Initial Setup
- [ ] Change default secret key
- [ ] Set strong admin password
- [ ] Configure HTTPS
- [ ] Set secure file permissions
- [ ] Enable proper logging
- [ ] Review service access

### Regular Maintenance
- [ ] Update ChronoServe regularly
- [ ] Review access logs
- [ ] Check file permissions
- [ ] Rotate secrets
- [ ] Audit user access
- [ ] Monitor service usage

## Known Security Considerations

### Windows Service Access
Windows services require administrative privileges. Ensure ChronoServe runs with minimal required permissions:

```powershell
# Create dedicated service account
New-LocalUser -Name "chronoserve" -Description "ChronoServe Service Account"
Add-LocalGroupMember -Group "Administrators" -Member "chronoserve"

# Set service to run as dedicated account
sc.exe config ChronoServe obj= ".\chronoserve" password= "secure-password"
```

### Linux Service Management
On Linux systems, configure sudo rules for specific service operations:

```bash
# /etc/sudoers.d/chronoserve
chronoserve ALL=(ALL) NOPASSWD: /bin/systemctl status *, /bin/systemctl start *, /bin/systemctl stop *
```

## Security Updates

Security updates will be released as patch versions. Subscribe to:
- GitHub releases
- Security advisories
- Version update notifications

## Compliance

ChronoServe follows security best practices:
- OWASP Top 10 guidelines
- NIST password recommendations
- Secure coding standards
- Regular security audits