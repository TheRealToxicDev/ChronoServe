# SysManix Version History and Compatibility

This guide provides information about SysManix versions, changes between versions, and compatibility considerations.

## Version Numbering

SysManix follows semantic versioning (SemVer) with the format `MAJOR.MINOR.PATCH`:

- **MAJOR**: Incremented for incompatible API changes
- **MINOR**: Incremented for new functionality in a backward-compatible manner
- **PATCH**: Incremented for backward-compatible bug fixes

## Current Version

The current stable version of SysManix is **0.1.0**.

## Version History

### v0.1.0 (Initial Release)

**Release Date:** 2023-06-01

**New Features:**
- Initial release with core functionality
- Cross-platform service management (Windows/Linux)
- JWT-based authentication system
- Role-based access control
- Service operations (list, status, start, stop, logs)
- Protected system service safety mechanism
- Configuration system with defaults

**Supported Operating Systems:**
- Windows 10/11, Windows Server 2016/2019/2022
- Ubuntu 20.04+
- Debian 11+
- CentOS/RHEL 8+
- Fedora 34+

**Dependencies:**
- Go 1.23.1
- Third-party libraries as listed in go.mod

## Upgrade Guide

### Upgrading to v0.1.0

As this is the initial release, no upgrade path is necessary.

For future versions, this section will include:
- Step-by-step upgrade instructions
- Configuration changes required
- Database migrations (if applicable)
- Breaking changes and mitigations

## API Compatibility

### v0.1.0 API Endpoints

The initial SysManix API includes the following stable endpoints:

#### Authentication
- **POST** `/auth/login`: Authenticate and get JWT token
- **GET** `/auth/tokens`: List your active tokens
- **POST** `/auth/tokens/revoke`: Revoke a specific token
- **POST** `/auth/tokens/refresh`: Refresh your current token

#### Service Management
- **GET** `/services`: List all services
- **GET** `/services/status/{service}`: Get service status
- **POST** `/services/start/{service}`: Start a service
- **POST** `/services/stop/{service}`: Stop a service
- **GET** `/services/logs/{service}`: View service logs

#### System
- **GET** `/health`: System health information

All endpoints are considered stable and will follow semantic versioning rules for any changes.

## Configuration Compatibility

### v0.1.0 Configuration

The initial configuration structure includes:

```yaml
server:
  host: "localhost"
  port: 40200
  readTimeout: "15s"
  writeTimeout: "15s"

auth:
  secretKey: "your-secure-random-string-here"
  tokenDuration: 24h
  issuedBy: "SysManix"
  allowedRoles:
    - admin
    - viewer
  users:
    admin:
      username: "admin"
      password_hash: "$argon2id$v=19$m=65536,t=1,p=4$..."
      roles:
        - admin
    viewer:
      username: "viewer"
      password_hash: "$argon2id$v=19$m=65536,t=1,p=4$..."
      roles:
        - viewer

windows:
  serviceCommand: "sc"
  logDirectory: "C:\\ProgramData\\SysManix\\logs"
  services:
    protected:
      - wininit
      - csrss
      - lsass

linux:
  serviceCommand: "systemctl"
  logDirectory: "/var/log/SysManix"
  services:
    protected:
      - systemd
      - systemd-journald
      - dbus

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
```

For future versions, any changes to this structure will be documented in the release notes.

## Known Issues

### v0.1.0 Known Issues

1. **Windows Token Storage**: On Windows, token storage uses in-memory database which resets on application restart
2. **Long service operations**: Service operations exceeding the timeout (10s default) may return a timeout error while the operation continues in the background
3. **Large service lists**: Listing services on systems with many services can be slow and memory-intensive
4. **PowerShell execution policy**: Restrictive PowerShell execution policies may block some service operations on Windows

## Platform-Specific Compatibility

### Windows Compatibility

| Windows Version | Compatibility | Notes |
|-----------------|---------------|-------|
| Windows 11 | ✅ Full | Tested and fully supported |
| Windows 10 | ✅ Full | Tested and fully supported |
| Windows Server 2022 | ✅ Full | Tested and fully supported |
| Windows Server 2019 | ✅ Full | Tested and fully supported |
| Windows Server 2016 | ✅ Full | Tested and fully supported |
| Windows 8.1 | ⚠️ Partial | Not officially supported, may work with limitations |
| Windows 7 | ❌ Not supported | PowerShell version requirements not met |

### Linux Compatibility

| Distribution | Compatibility | Notes |
|--------------|---------------|-------|
| Ubuntu 22.04 LTS | ✅ Full | Tested and fully supported |
| Ubuntu 20.04 LTS | ✅ Full | Tested and fully supported |
| Debian 12 | ✅ Full | Tested and fully supported |
| Debian 11 | ✅ Full | Tested and fully supported |
| CentOS/RHEL 9 | ✅ Full | Tested and fully supported |
| CentOS/RHEL 8 | ✅ Full | Tested and fully supported |
| Fedora 38+ | ✅ Full | Tested and fully supported |
| Alpine Linux | ⚠️ Partial | May require manual configuration |
| Arch Linux | ✅ Full | Tested and working as expected |

## Future Development

### Planned for v0.2.0

- User management API endpoints
- Persistent token storage
- Enhanced service filtering and search
- Service dependency visualization
- Docker integration and container management
- Improved performance for large service lists
- Additional authentication methods

### Under Consideration

Features being considered for future versions:

- Web-based user interface
- Service metrics collection and visualization
- Event-driven service monitoring
- Service auto-recovery rules
- Scheduled service operations
- Clustered deployment for high availability

## Backward Compatibility Policy

SysManix follows these backward compatibility principles:

1. **API Stability**: Existing API endpoints will maintain backward compatibility within the same major version
2. **Configuration Format**: Configuration file format changes will be backward compatible within the same major version
3. **Database Schema**: Future database schema changes will include automatic migrations
4. **Deprecation Policy**: Features will be marked deprecated before removal with at least one minor version cycle notice

## Reporting Issues

If you encounter any issues or incompatibilities:

1. Check this version document and the [Troubleshooting Guide](./TROUBLESHOOTING.md)
2. Search existing issues on the [GitHub repository](https://github.com/toxic-development/sysmanix/issues)
3. Submit a new issue with:
   - SysManix version details
   - Operating system and version
   - Steps to reproduce
   - Expected vs. actual behavior
   - Any relevant logs or error messages

## Security Updates

Security updates may be released as patch versions and should be applied promptly. Critical security issues will be highlighted in release notes and may receive expedited releases.
