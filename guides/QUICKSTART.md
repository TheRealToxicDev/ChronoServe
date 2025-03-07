# SysManix Quick Start Guide

This guide will help you get up and running with SysManix quickly. For more detailed information, refer to the other documentation files.

## Installation

### Windows
```powershell
# Download the latest release
Invoke-WebRequest -Uri "https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_windows_amd64.exe" -OutFile "SysManix.exe"

# Run the executable to generate initial config
.\SysManix.exe
```

### Linux
```bash
# Download the latest release
wget https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_linux_amd64 -O sysmanix
chmod +x sysmanix

# Run to generate initial config (requires root for service management)
sudo ./sysmanix
```

## First-Time Setup

1. **Edit the configuration file** (created in the same directory as the executable by default):

   ```yaml
   # Minimal required changes
   auth:
     secretKey: "your-secure-random-string-here"
     users:
       admin:
         password: "your-secure-admin-password"
       viewer:
         password: "your-secure-viewer-password"

   server:
     port: 40200  # Change if needed
   ```

2. **Start SysManix**:

   ```bash
   # Windows
   .\SysManix.exe

   # Linux
   sudo ./sysmanix
   ```

## Accessing the API

### Authentication

1. **Get a JWT token**:

   ```bash
   # Using curl
   curl -X POST http://localhost:40200/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"your-secure-admin-password"}'
   ```

   The response will contain your JWT token:

   ```json
   {
     "status": "success",
     "message": "Login successful",
     "data": {
       "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
       "roles": ["admin"]
     }
   }
   ```

2. **Save the token** for use in subsequent requests:

   ```bash
   export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   ```

### Basic Service Management

1. **List all services**:

   ```bash
   curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services
   ```

2. **Check a specific service status**:

   ```bash
   # Replace "nginx" with your service name
   curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services/status/nginx
   ```

3. **Start a service**:

   ```bash
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     http://localhost:40200/services/start/nginx
   ```

4. **Stop a service**:

   ```bash
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     http://localhost:40200/services/stop/nginx
   ```

5. **View service logs**:

   ```bash
   # Get the last 50 log lines
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:40200/services/logs/nginx?lines=50
   ```

## Token Management

1. **List your active tokens**:

   ```bash
   curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:40200/auth/tokens
   ```

2. **Refresh your token** (generates a new token and invalidates the current one):

   ```bash
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     http://localhost:40200/auth/tokens/refresh
   ```

   Remember to update your saved token with the new one from the response.

3. **Revoke a specific token**:

   ```bash
   # Replace tokenId with the ID from your tokens list
   curl -X POST -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"tokenId":"01H5RZ4WH5P3Z72SF"}' \
     http://localhost:40200/auth/tokens/revoke
   ```

## Running as a Service

### Windows Service Setup

```powershell
# Using NSSM (https://nssm.cc/)
nssm.exe install SysManix "[path\to\SysManix.exe]"
nssm.exe set SysManix DisplayName "SysManix Service Manager"
nssm.exe set SysManix Start SERVICE_AUTO_START
Start-Service SysManix
```

### Linux Systemd Setup

```bash
# Create a systemd service file
sudo nano /etc/systemd/system/sysmanix.service
```

Add the following content:

```ini
[Unit]
Description=SysManix Service Manager
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/path/to/sysmanix/folder
ExecStart=/path/to/sysmanix/folder/sysmanix
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable sysmanix
sudo systemctl start sysmanix
```

## Health Check

To verify SysManix is running correctly:

```bash
curl http://localhost:40200/health
```

This endpoint is publicly accessible without authentication and returns system status information.

## Where to Go Next

Now that you have SysManix up and running, explore these guides for more information:

- [Configuration Guide](./CONFIGURATION.md) - Customize SysManix to your needs
- [Authentication Guide](./AUTHENTICATION.md) - Advanced authentication options
- [Service Management Guide](./SERVICE_MANAGEMENT.md) - Detailed service operations
- [API Reference](./API_REFERENCE.md) - Complete API documentation

## Troubleshooting

If you encounter issues:

1. Check the logs directory for error messages
2. Ensure you have proper permissions to manage services
3. Verify the configuration file syntax
4. See the [Troubleshooting Guide](./TROUBLESHOOTING.md) for common problems and solutions
