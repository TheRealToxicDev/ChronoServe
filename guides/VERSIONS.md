# ChronoServe Version Management

## Checking Your Version

You can check your current ChronoServe version in several ways:

### Command Line
```bash
chronoserve --version
```

### From Running Server
The version is displayed in:
1. Health endpoint response
2. Log file headers

## Version Format

ChronoServe follows semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: New features, backwards compatible
- PATCH: Bug fixes, backwards compatible

Example: `v1.2.3`

## Update Checking

ChronoServe automatically checks for updates by:
1. Comparing your version with latest GitHub release
2. Displaying notifications if updates are available
3. Providing update instructions based on your OS

### Manual Update Check
```bash
chronoserve --check-updates
```

## Downloading Updates

### Windows
```powershell
# PowerShell (as Administrator)
1. Stop ChronoServe service
Stop-Service ChronoServe

2. Download latest release
Invoke-WebRequest -Uri "https://github.com/therealtoxicdev/chronoserve/releases/latest/download/chronoserve_windows_amd64.exe" -OutFile ".\bin\chronoserve.exe"

3. Restart service
Start-Service ChronoServe
```

### Linux
```bash
# Terminal (with sudo)
1. Stop ChronoServe service
sudo systemctl stop chronoserve

2. Download latest release
sudo wget https://github.com/therealtoxicdev/chronoserve/releases/latest/download/chronoserve_linux_amd64 -O /usr/local/bin/chronoserve

3. Make executable
sudo chmod +x /usr/local/bin/chronoserve

4. Restart service
sudo systemctl start chronoserve
```

## Version History

### Latest Version (v0.1.0)
- Initial release
- Basic service management
- JWT authentication
- Role-based access control

### Update Notifications

Example update available:
```
ChronoServe v0.1.0 (update available: v0.2.0)
Visit https://github.com/therealtoxicdev/chronoserve/releases for the latest version
```

Example up-to-date:
```
ChronoServe v0.1.0 (latest)
```

## Configuration

Version checking can be configured in `config.yaml`:

```yaml
updates:
  checkOnStartup: true    # Check for updates when server starts
  notifyInLogs: true      # Show update notifications in logs
  checkInterval: 24h      # How often to check for updates
  githubTimeout: 10s      # Timeout for GitHub API requests
```

## Troubleshooting

### Common Issues

1. "Failed to check for updates"
   - Check your internet connection
   - Verify GitHub API access
   - Check firewall settings

2. "Version check timed out"
   - Increase `githubTimeout` in config
   - Check network latency

3. "Invalid version format"
   - Clear cached version data
   - Reinstall from official release

### Support

For version-related issues:
1. Check [GitHub Issues](https://github.com/therealtoxicdev/chronoserve/issues)
2. Open a new issue with label 'version'
3. Include your current version and OS details