# SysManix Installation Guide

This guide provides detailed instructions for installing SysManix on different operating systems.

## System Requirements

### Minimum Requirements

- **CPU**: 1 core
- **RAM**: 256 MB
- **Disk Space**: 50 MB
- **Operating System**:
  - Windows 10/11 or Windows Server 2016/2019/2022
  - Linux with systemd (Ubuntu 20.04+, Debian 11+, CentOS 8+, etc.)

### Recommended Requirements

- **CPU**: 2+ cores
- **RAM**: 512 MB+
- **Disk Space**: 100 MB+
- **Network**: Stable connection if managing remote services

### Required Permissions

- Administrative privileges (Windows) or root access (Linux)
- Permission to manage system services
- Network access for remote management (if applicable)

## Installation Methods

### Installing Pre-built Binaries

#### Windows Installation

1. **Download the latest release**:
   - Visit the [GitHub releases page](https://github.com/toxic-development/sysmanix/releases)
   - Download the Windows binary (`SysManix_windows_amd64.exe`)

2. **Place the executable in your preferred location**:
   ```powershell
   # Example: Create program directory and move the file
   New-Item -ItemType Directory -Path "C:\Program Files\SysManix" -Force
   Move-Item -Path .\SysManix_windows_amd64.exe -Destination "C:\Program Files\SysManix\SysManix.exe"
   ```

3. **Create a basic configuration file**:
   ```powershell
   # Run once to generate a default config file
   & "C:\Program Files\SysManix\SysManix.exe"
   ```

4. **Edit the configuration** (located at `config.yaml` in the same directory as the executable):
   - Update the `secretKey` with a secure random string
   - Change default passwords
   - Adjust other settings as needed

5. **Install as a Windows Service (optional but recommended)**:
   ```powershell
   # Using NSSM (Non-Sucking Service Manager)
   # Download NSSM from https://nssm.cc/ first
   .\nssm.exe install SysManix "C:\Program Files\SysManix\SysManix.exe"
   .\nssm.exe set SysManix DisplayName "SysManix Service Manager"
   .\nssm.exe set SysManix Description "Cross-platform service management API"
   .\nssm.exe set SysManix Start SERVICE_AUTO_START
   .\nssm.exe set SysManix AppDirectory "C:\Program Files\SysManix"
   .\nssm.exe set SysManix AppParameters "-config C:\Program Files\SysManix\config.yaml"

   # Start the service
   Start-Service SysManix
   ```

#### Linux Installation

1. **Download the latest release**:
   ```bash
   mkdir -p /opt/sysmanix
   cd /opt/sysmanix
   wget https://github.com/toxic-development/sysmanix/releases/latest/download/SysManix_linux_amd64 -O sysmanix
   chmod +x sysmanix
   ```

2. **Create a basic configuration file**:
   ```bash
   # Run once to generate a default config
   sudo ./sysmanix
   ```

3. **Edit the configuration**:
   ```bash
   sudo nano config.yaml
   # Update secretKey, passwords, and other settings
   ```

4. **Create a systemd service file**:
   ```bash
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
   WorkingDirectory=/opt/sysmanix
   ExecStart=/opt/sysmanix/sysmanix -config /opt/sysmanix/config.yaml
   Restart=on-failure
   RestartSec=5s

   [Install]
   WantedBy=multi-user.target
   ```

5. **Enable and start the service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable sysmanix
   sudo systemctl start sysmanix
   ```

### Building from Source

#### Prerequisites

- Go 1.23.1 or newer
- Git
- Build tools (gcc, etc.)

#### Build Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/toxic-development/sysmanix.git
   cd sysmanix
   ```

2. **Build the application**:
   ```bash
   # For the current platform
   go build -o sysmanix ./client

   # Cross-compile for Windows from Linux
   GOOS=windows GOARCH=amd64 go build -o sysmanix.exe ./client

   # Cross-compile for Linux from Windows
   set GOOS=linux
   set GOARCH=amd64
   go build -o sysmanix ./client
   ```

3. **Run tests** (optional but recommended):
   ```bash
   go test -v ./...
   ```

4. **Continue with platform-specific installation** (see sections above)

## Docker Installation

### Using Docker

1. **Pull the Docker image**:
   ```bash
   docker pull toxic-development/sysmanix:latest
   ```

2. **Create a configuration directory**:
   ```bash
   mkdir -p /etc/sysmanix
   ```

3. **Run the container**:
   ```bash
   docker run -d \
     --name sysmanix \
     -p 40200:40200 \
     -v /etc/sysmanix:/app/config \
     -v /var/run/docker.sock:/var/run/docker.sock \
     --restart unless-stopped \
     toxic-development/sysmanix:latest
   ```

4. **Edit the generated configuration**:
   ```bash
   nano /etc/sysmanix/config.yaml
   # Update configuration as needed
   ```

5. **Restart the container**:
   ```bash
   docker restart sysmanix
   ```

### Using Docker Compose

1. **Create a docker-compose.yml file**:
   ```yaml
   version: '3'

   services:
     sysmanix:
       image: toxic-development/sysmanix:latest
       container_name: sysmanix
       ports:
         - "40200:40200"
       volumes:
         - ./config:/app/config
         - /var/run/docker.sock:/var/run/docker.sock
       restart: unless-stopped
   ```

2. **Start the service**:
   ```bash
   docker-compose up -d
   ```

3. **Edit the generated configuration**:
   ```bash
   nano ./config/config.yaml
   # Update configuration as needed
   ```

4. **Restart the service**:
   ```bash
   docker-compose restart
   ```

## Post-Installation Steps

### Security Checklist

After installing SysManix, follow these security best practices:

1. **Change all default credentials**:
   - Update the admin and viewer passwords
   - Generate a secure random string for `secretKey`

2. **Set appropriate file permissions**:
   - Windows: Restrict config.yaml to Administrators
   - Linux: `chmod 640 config.yaml` and appropriate ownership

3. **Configure firewall rules**:
   - Windows: Allow inbound connections to the configured port
   - Linux: Configure iptables or ufw to restrict access

4. **Set up HTTPS** (recommended for production):
   - Use a reverse proxy like Nginx or Caddy
   - Configure with proper TLS certificates

### Initial Configuration

1. **Verify the installation**:
   ```bash
   # Check if the service is running
   curl http://localhost:40200/health
   ```

2. **Test authentication**:
   ```bash
   # Using curl
   curl -X POST http://localhost:40200/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"your-password"}'
   ```

3. **Update configuration as needed**:
   - Adjust `tokenDuration` for your security requirements
   - Configure logging levels and locations
   - Add additional users and roles

## Upgrading

### Standard Upgrade Process

1. **Stop the running service**:
   ```bash
   # Windows
   Stop-Service SysManix

   # Linux
   sudo systemctl stop sysmanix
   ```

2. **Backup configuration**:
   ```bash
   # Windows
   Copy-Item "C:\Program Files\SysManix\config.yaml" "C:\Program Files\SysManix\config.yaml.bak"

   # Linux
   sudo cp /opt/sysmanix/config.yaml /opt/sysmanix/config.yaml.bak
   ```

3. **Replace the executable**:
   ```bash
   # Windows
   Move-Item -Path .\SysManix_windows_amd64.exe -Destination "C:\Program Files\SysManix\SysManix.exe" -Force

   # Linux
   sudo cp ./sysmanix /opt/sysmanix/sysmanix
   sudo chmod +x /opt/sysmanix/sysmanix
   ```

4. **Start the service**:
   ```bash
   # Windows
   Start-Service SysManix

   # Linux
   sudo systemctl start sysmanix
   ```

5. **Verify the upgrade**:
   ```bash
   curl http://localhost:40200/health
   ```

### Docker Upgrade

```bash
# Pull the latest image
docker pull toxic-development/sysmanix:latest

# Stop and remove the current container
docker stop sysmanix
docker rm sysmanix

# Start a new container with the latest image
docker run -d \
  --name sysmanix \
  -p 40200:40200 \
  -v /etc/sysmanix:/app/config \
  --restart unless-stopped \
  toxic-development/sysmanix:latest
```

## Troubleshooting Installation

### Common Installation Issues

- **"Configuration file not found" error**:
  - Verify the correct path to the config file
  - Ensure the directory exists and has proper permissions

- **"Port already in use" error**:
  - Change the port in `config.yaml`
  - Check for other applications using the default port

- **Service fails to start**:
  - Check system logs for detailed error messages
  - Verify all dependencies are installed
  - Confirm proper permissions for executable and config

### Viewing Logs

```bash
# Windows (Event Viewer)
Get-EventLog -LogName Application | Where-Object {$_.Source -eq "SysManix"}

# Linux (systemd)
sudo journalctl -u sysmanix.service -f
```

### Getting Help

If you encounter installation issues not covered in this guide:

1. Check the [troubleshooting guide](./TROUBLESHOOTING.md) for more detailed solutions
2. Visit the [GitHub issues page](https://github.com/toxic-development/sysmanix/issues) to check for known issues
3. Ask for help in the community discussion forums
