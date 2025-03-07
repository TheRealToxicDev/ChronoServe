# SysManix Installation Guide

This guide covers the installation process for SysManix on both Windows and Linux platforms.

## Prerequisites

Before installing SysManix, ensure your system meets the following requirements:

### Windows Requirements
- Windows 10/11 or Windows Server 2019/2022
- Administrator privileges
- PowerShell 5.1 or newer
- .NET Framework 4.5 or newer (for certain Windows services)

### Linux Requirements
- A modern Linux distribution (Ubuntu 20.04+, CentOS 8+, Debian 11+, etc.)
- Systemd as the init system
- Root or sudo privileges
- journalctl access for log viewing

### Common Requirements
- Outbound network access (for update checking)
- Minimum 1GB RAM (2GB recommended)
- 100MB free disk space

## Installation Methods

### Windows Installation

#### Using the Installer

1. Download the latest SysManix Windows installer from the [releases page](https://github.com/toxic-development/SysManix/releases)
2. Run the installer as administrator
3. Follow the installation wizard instructions
4. Choose an installation directory (default: `C:\Program Files\SysManix`)
5. Select whether to create a desktop shortcut and/or start menu entry
6. Complete the installation

#### Manual Installation

1. Download the latest Windows ZIP package from the [releases page](https://github.com/toxic-development/SysManix/releases)
2. Extract the ZIP to your desired location (e.g., `C:\Program Files\SysManix`)
3. Open PowerShell as administrator and navigate to the installation directory
4. Run the setup script: `.\setup.ps1`
5. Follow the instructions to complete the installation

### Linux Installation

#### Using the Package Manager (Recommended)

**Ubuntu/Debian:**
```bash
# Add the SysManix repository
curl -fsSL https://download.sysmanix.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/sysmanix-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/sysmanix-archive-keyring.gpg] https://download.sysmanix.io/apt stable main" | sudo tee /etc/apt/sources.list.d/sysmanix.list

# Update package lists and install SysManix
sudo apt update
sudo apt install sysmanix
```

**RHEL/CentOS/Fedora:**
```bash
# Add the SysManix repository
sudo tee /etc/yum.repos.d/sysmanix.repo << EOF
[sysmanix]
name=SysManix Repository
baseurl=https://download.sysmanix.io/rpm
enabled=1
gpgcheck=1
gpgkey=https://download.sysmanix.io/gpg
EOF

# Install SysManix
sudo yum install sysmanix
```

#### Using the Binary Package

1. Download the latest Linux tarball from the [releases page](https://github.com/toxic-development/SysManix/releases)
2. Extract the tarball to a temporary location
   ```bash
   tar -xzf sysmanix-0.1.0-linux-amd64.tar.gz -C /tmp
   ```
3. Run the installation script
   ```bash
   cd /tmp/sysmanix-0.1.0
   sudo ./install.sh
   ```
4. Follow the instructions to complete the installation

## Building From Source

### Prerequisites
- Go 1.23.1 or newer
- Git
- Make

### Steps

1. Clone the repository
   ```bash
   git clone https://github.com/toxic-development/SysManix.git
   cd SysManix
   ```

2. Build the application
   ```bash
   make build
   ```

3. Install the application
   ```bash
   make install
   ```

## Post-Installation

After installing SysManix, you'll need to:

1. Configure the application (see the [Configuration Guide](./CONFIGURATION.md))
2. Set up secure user credentials
3. Configure the firewall to allow access to the API port (default: 40200)
4. Set up the service to start automatically (see [Systemd Setup](./SYSTEMD_SETUP.md) for Linux)

## Verifying the Installation

To verify that SysManix is installed correctly:

### On Windows
```powershell
# Check the service status
Get-Service -Name SysManix

# Access the health endpoint
Invoke-RestMethod -Uri "http://localhost:40200/health"
```

### On Linux
```bash
# Check the service status
systemctl status sysmanix

# Access the health endpoint
curl http://localhost:40200/health
```

You should see a JSON response indicating the service is healthy.

## Common Installation Issues

### Port Conflicts
If another application is using port 40200, you can change the port in the configuration file (`config.yaml`).

### Permission Issues
Ensure you have the necessary permissions to install and run the application. On Linux, you may need to use `sudo`. On Windows, run the installer or PowerShell as administrator.

### Service Startup Failures
Check the logs in the configured log directory for error messages. The most common issues are configuration errors or permission problems.

## Next Steps

- [Quick Start Guide](./QUICKSTART.md): Get up and running with SysManix
- [Configuration Guide](./CONFIGURATION.md): Learn about configuration options
- [Systemd Setup](./SYSTEMD_SETUP.md): Configure SysManix as a systemd service on Linux
- [Windows Service Setup](./WINDOWS_SETUP.md): Configure SysManix as a Windows service
- [Security Guide](./SECURITY.md): Best practices for securing your SysManix installation
