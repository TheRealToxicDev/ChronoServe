# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously at ChronoServe. Please follow these steps to report security issues:

1. **DO NOT** open public issues for security vulnerabilities
2. Send reports to [security@example.com](mailto:security@example.com)
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
```bash
# Windows (PowerShell)
icacls .\config.yaml /inheritance:r /grant:r "SYSTEM:(R)" "ADMINISTRATORS:(F)"

# Linux
chmod 600 config.yaml
chown root:root config.yaml
```

## Security Features

- JWT-based authentication
- Role-based access control
- Password security validation
- Rate limiting
- Secure default configurations
- Automated security checks