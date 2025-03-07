# SysManix Systemd Setup Guide

This guide explains how to properly configure SysManix as a systemd service on Linux systems.

## Basic Systemd Service Setup

### Creating the Service File

1. Create a systemd service file:

```bash
sudo nano /etc/systemd/system/sysmanix.service
```

2. Add the following content:

```ini
[Unit]
Description=SysManix Service Manager
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/sysmanix
ExecStart=/opt/sysmanix/sysmanix -config /etc/sysmanix/config.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

3. Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable sysmanix
sudo systemctl start sysmanix
```

## Advanced Service Configuration

### Environment Variables

You can set environment variables in the systemd service file:

```ini
[Service]
Environment="SYSMANIX_LOG_LEVEL=debug"
Environment="SYSMANIX_SERVER_PORT=8080"
```

Alternatively, create an environment file:

```bash
sudo nano /etc/sysmanix/environment
```

Add your environment variables:

```
SYSMANIX_LOG_LEVEL=debug
SYSMANIX_SERVER_PORT=8080
```

Then reference it in your service file:

```ini
[Service]
EnvironmentFile=/etc/sysmanix/environment
```

### Resource Limits

Control system resources used by SysManix:

```ini
[Service]
CPUQuota=50%
MemoryLimit=512M
TasksMax=100
LimitNOFILE=65535
```

## Service Dependencies

If SysManix depends on other services (like a database), you can specify them:

```ini
[Unit]
Description=SysManix Service Manager
After=network.target postgresql.service
Requires=postgresql.service
```

## Monitoring and Logging

### Journald Integration

Systemd automatically captures stdout/stderr from the service. You can view the logs with:

```bash
sudo journalctl -u sysmanix
```

### Custom Log Configuration

If you prefer separate log files, modify your service file:

```ini
[Service]
StandardOutput=append:/var/log/sysmanix/stdout.log
StandardError=append:/var/log/sysmanix/stderr.log
```

Ensure the log directory exists and has appropriate permissions:

```bash
sudo mkdir -p /var/log/sysmanix
sudo chown -R sysmanix:sysmanix /var/log/sysmanix
```

## Service Recovery

Configure how systemd handles service failures:

```ini
[Service]
Restart=on-failure
RestartSec=5s
StartLimitInterval=500s
StartLimitBurst=5
```

For critical services, you can add additional actions:

```ini
[Service]
ExecStopPost=/opt/sysmanix/scripts/notify-failure.sh
```

## Watchdog Integration

If SysManix supports watchdog notifications, configure it:

```ini
[Service]
WatchdogSec=30s
```

Your application would need to periodically call `sd_notify("WATCHDOG=1")` to inform systemd it's still functioning properly.

## Creating a Socket-Activated Service

For more responsive startup, you can use socket activation:

1. Create a socket file:

```bash
sudo nano /etc/systemd/system/sysmanix.socket
```

2. Add socket configuration:

```ini
[Socket]
ListenStream=40200
BindIPv6Only=both

[Install]
WantedBy=sockets.target
```

3. Modify the service file to work with socket activation:

```ini
[Unit]
Description=SysManix Service Manager
After=network.target
Requires=sysmanix.socket

[Service]
ExecStart=/opt/sysmanix/sysmanix -config /etc/sysmanix/config.yaml
```

4. Enable and start the socket:

```bash
sudo systemctl enable sysmanix.socket
sudo systemctl start sysmanix.socket
```

## Creating Auxiliary Services

### Backup Service

Create an automated backup service for SysManix:

1. Create a backup script:

```bash
sudo nano /opt/sysmanix/scripts/backup.sh
```

Add script content:

```bash
#!/bin/bash
BACKUP_DIR="/var/backups/sysmanix"
DATE=$(date +"%Y-%m-%d_%H-%M-%S")
mkdir -p "$BACKUP_DIR"
cp -r /etc/sysmanix "$BACKUP_DIR/config_$DATE"
find "$BACKUP_DIR" -type d -mtime +30 -exec rm -rf {} \; 2>/dev/null || true
```

Make it executable:

```bash
sudo chmod +x /opt/sysmanix/scripts/backup.sh
```

2. Create a timer and service:

```bash
sudo nano /etc/systemd/system/sysmanix-backup.service
```

Add service content:

```ini
[Unit]
Description=SysManix Configuration Backup Service

[Service]
Type=oneshot
ExecStart=/opt/sysmanix/scripts/backup.sh
```

Create the timer:

```bash
sudo nano /etc/systemd/system/sysmanix-backup.timer
```

Add timer content:

```ini
[Unit]
Description=Run SysManix backup daily

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
```

3. Enable and start the timer:

```bash
sudo systemctl enable sysmanix-backup.timer
sudo systemctl start sysmanix-backup.timer
```

## Service Health Monitoring

### Creating a Health Check Service

1. Create a health check script:

```bash
sudo nano /opt/sysmanix/scripts/health-check.sh
```

Add script content:

```bash
#!/bin/bash
HEALTH_ENDPOINT="http://localhost:40200/health"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$HEALTH_ENDPOINT")

if [ "$STATUS" != "200" ]; then
    echo "SysManix health check failed with status $STATUS"
    exit 1
fi

exit 0
```

Make it executable:

```bash
sudo chmod +x /opt/sysmanix/scripts/health-check.sh
```

2. Create a service and timer:

```bash
sudo nano /etc/systemd/system/sysmanix-health.service
```

Add service content:

```ini
[Unit]
Description=SysManix Health Check Service

[Service]
Type=oneshot
ExecStart=/opt/sysmanix/scripts/health-check.sh
```

Create the timer:

```bash
sudo nano /etc/systemd/system/sysmanix-health.timer
```

Add timer content:

```ini
[Unit]
Description=Run SysManix health check every 5 minutes

[Timer]
OnCalendar=*:0/5
Persistent=true

[Install]
WantedBy=timers.target
```

3. Enable and start the timer:

```bash
sudo systemctl enable sysmanix-health.timer
sudo systemctl start sysmanix-health.timer
```

## Troubleshooting Systemd Services

### Common Issues and Solutions

1. **Service fails to start**:

   Check the service status:
   ```bash
   sudo systemctl status sysmanix
   ```

   View detailed error logs:
   ```bash
   sudo journalctl -u sysmanix -n 100
   ```

2. **"Permission denied" errors**:

   Check the service user and file permissions:
   ```bash
   sudo ls -la /opt/sysmanix/sysmanix
   sudo ls -la /etc/sysmanix/config.yaml
   ```

   Ensure the service user has access to all required resources.

3. **Service starts but exits immediately**:

   Look for potential configuration errors:
   ```bash
   sudo journalctl -u sysmanix -b
   ```

   Try running the application manually to see errors:
   ```bash
   sudo -u root /opt/sysmanix/sysmanix -config /etc/sysmanix/config.yaml
   ```

## Performance Optimization

### Adjusting systemd Settings

For optimal performance and faster service startup:

1. Modify `/etc/systemd/system.conf`:

```ini
DefaultTimeoutStartSec=30s
DefaultTimeoutStopSec=30s
DefaultRestartSec=5s
```

2. Adjust service process priority:

```ini
[Service]
Nice=-10  # Higher priority (range: -20 to 19, lower is higher priority)
CPUSchedulingPolicy=fifo
CPUSchedulingPriority=50  # Range: 1-99
IOSchedulingClass=1  # 1 for real-time, 2 for best-effort
IOSchedulingPriority=0  # Range: 0-7, lower is higher priority
```

## Security Hardening

Enhance the security of your SysManix service:

```ini
[Service]
# Restrict capabilities
CapabilityBoundingSet=CAP_NET_BIND_SERVICE

# File system restrictions
ProtectSystem=strict
ProtectHome=true
PrivateDevices=true
PrivateTmp=true
ReadOnlyPaths=/etc/sysmanix
ReadWritePaths=/var/log/sysmanix /var/lib/sysmanix
NoExecPaths=/tmp

# Process restrictions
RestrictRealtime=true
RestrictSUIDSGID=true
RestrictNamespaces=true

# Memory protection
MemoryDenyWriteExecute=true

# Network security
PrivateNetwork=false  # Set to true if you don't need network access
```

Note: Some restrictions may need to be adjusted depending on SysManix's requirements.

## Conclusion

With proper systemd configuration, SysManix can run reliably as a system service with:

- Automatic startup and recovery
- Resource limits to prevent system overload
- Security hardening to protect your system
- Monitoring to ensure continuous operation
- Scheduled backups to prevent data loss

For more information, refer to the [systemd documentation](https://www.freedesktop.org/software/systemd/man/systemd.service.html).
