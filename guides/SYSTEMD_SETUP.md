# SysManix Systemd Integration Guide

This guide explains how to set up SysManix as a systemd service on Linux systems for proper service management.

## Prerequisites

- SysManix installed on a Linux system (see [Installation Guide](./INSTALLATION.md))
- Root or sudo access
- Systemd as the init system (standard on most modern Linux distributions)

## Understanding Systemd Integration

Running SysManix as a systemd service provides several benefits:

- **Automatic startup** on system boot
- **Dependency management** with other services
- **Restart policies** for increased reliability
- **Resource control** via systemd slices and limits
- **Standardized logging** through journald

## Creating a Systemd Service File

1. Create a new systemd service file:

```bash
sudo nano /etc/systemd/system/sysmanix.service
```

2. Add the following content, adjusting paths as necessary:

```ini
[Unit]
Description=SysManix Service Management API
Documentation=https://github.com/toxic-development/SysManix
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=sysmanix
Group=sysmanix
ExecStart=/usr/bin/sysmanix serve
Restart=on-failure
RestartSec=5
TimeoutStartSec=30
TimeoutStopSec=30

# Security hardening
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_BIND_SERVICE
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=read-only
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
PrivateTmp=true

# Working directory and environment
WorkingDirectory=/etc/sysmanix
Environment="CONFIG_PATH=/etc/sysmanix/config.yaml"

[Install]
WantedBy=multi-user.target
```

## Understanding Service Configuration Options

### Basic Service Settings

- **User/Group**: It's recommended to run SysManix as a dedicated user for security
- **ExecStart**: Path to the SysManix binary with necessary arguments
- **Restart**: Determines when systemd should restart the service
- **RestartSec**: How long to wait before restarting the service
- **TimeoutStartSec/TimeoutStopSec**: Maximum time allowed for start/stop operations

### Security Hardening Settings

- **CapabilityBoundingSet**: Limits the capabilities of the process
- **NoNewPrivileges**: Prevents gaining new privileges via execve()
- **ProtectSystem**: Restricts write access to the system
- **ProtectHome**: Controls access to home directories
- **PrivateTmp**: Provides an isolated /tmp directory

## Creating a Dedicated User

For enhanced security, create a dedicated user for running SysManix:

```bash
# Create sysmanix user and group
sudo useradd -r -s /bin/false -m -d /var/lib/sysmanix sysmanix

# Create necessary directories
sudo mkdir -p /etc/sysmanix /var/log/sysmanix

# Set permissions
sudo chown -R sysmanix:sysmanix /etc/sysmanix /var/log/sysmanix
sudo chmod 750 /etc/sysmanix /var/log/sysmanix
```

## Configuring SysManix for Systemd

Ensure your SysManix configuration is placed in `/etc/sysmanix/config.yaml`. Adjust the following settings:

```yaml
server:
  host: "0.0.0.0"  # Listen on all interfaces
  port: 40200

logging:
  directory: "/var/log/sysmanix"

linux:
  logDirectory: "/var/log/sysmanix"
```

## Starting and Enabling the Service

1. Reload systemd to recognize the new service file:

```bash
sudo systemctl daemon-reload
```

2. Start the SysManix service:

```bash
sudo systemctl start sysmanix
```

3. Enable the service to start automatically on boot:

```bash
sudo systemctl enable sysmanix
```

4. Verify the service status:

```bash
sudo systemctl status sysmanix
```

## Viewing Logs

With systemd integration, you can view SysManix logs using journalctl:

```bash
# View all SysManix logs
sudo journalctl -u sysmanix

# Follow logs in real-time
sudo journalctl -u sysmanix -f

# View logs since the last boot
sudo journalctl -u sysmanix -b

# View logs with specific log level
sudo journalctl -u sysmanix -p err
```

## Service Management

### Basic Service Control

```bash
# Stop the service
sudo systemctl stop sysmanix

# Start the service
sudo systemctl start sysmanix

# Restart the service
sudo systemctl restart sysmanix

# Reload the service (if supported)
sudo systemctl reload sysmanix
```

### Checking Service Status

```bash
# View detailed status
sudo systemctl status sysmanix

# Check if service is active
sudo systemctl is-active sysmanix

# Check if service is enabled at boot
sudo systemctl is-enabled sysmanix

# View service dependencies
sudo systemctl list-dependencies sysmanix
```

## Advanced Systemd Configuration

### Resource Limits

You can add resource limits to your service by adding the following to the `[Service]` section:

```ini
# Memory limit (example: 2GB)
MemoryLimit=2G

# CPU limits
CPUQuota=50%

# Process limits
LimitNPROC=64
LimitNOFILE=4096
```

### Restart Policies

For a more robust restart policy, adjust these settings:

```ini
Restart=always
RestartSec=5
StartLimitIntervalSec=500
StartLimitBurst=5
```

### Environment Variables

Add additional environment variables:

```ini
Environment="DEBUG=false"
Environment="LOG_LEVEL=info"
EnvironmentFile=/etc/sysmanix/env
```

## Troubleshooting

### Service Fails to Start

Check the logs for detailed error messages:

```bash
sudo journalctl -u sysmanix -b -n 50
```

Common issues include:
- Incorrect file permissions
- Configuration errors
- Port already in use
- Insufficient permissions for the sysmanix user

### Permission Issues

If SysManix needs access to system services but runs as a non-root user:

```bash
# Add specific capabilities
sudo setcap cap_dac_override,cap_net_bind_service=+ep /usr/bin/sysmanix
```

### Failed to Load Configuration

Ensure the configuration file exists and has correct permissions:

```bash
sudo ls -la /etc/sysmanix/config.yaml
sudo chown sysmanix:sysmanix /etc/sysmanix/config.yaml
sudo chmod 640 /etc/sysmanix/config.yaml
```

## Security Considerations

1. **Firewall Configuration**: Restrict access to the SysManix API port
   ```bash
   sudo ufw allow from 192.168.1.0/24 to any port 40200 proto tcp
   ```

2. **Regular Updates**: Keep SysManix updated to get the latest security patches
   ```bash
   sudo systemctl stop sysmanix
   # update SysManix
   sudo systemctl start sysmanix
   ```

3. **Audit Logging**: Enable audit logging in your systemd configuration
   ```ini
   LogExtraFields=AUDIT_SESSION_ID=123 ACTION=start
   ```

## Further Reading

- [Systemd Documentation](https://www.freedesktop.org/software/systemd/man/systemd.service.html)
- [Linux Security Best Practices](./LINUX_SECURITY.md)
- [SysManix Configuration Guide](./CONFIGURATION.md)
