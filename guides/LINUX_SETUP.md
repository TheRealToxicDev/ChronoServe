# SysManix Linux Setup Guide

This guide provides detailed instructions for setting up SysManix on Linux systems, focusing on Linux-specific considerations and optimizations.

## System Requirements

- A modern Linux distribution with systemd (Ubuntu 20.04+, Debian 11+, CentOS/RHEL 8+, etc.)
- Root or sudo privileges
- systemd as the init system
- journalctl access for log viewing

## Installation Options

### Option 1: Standalone Binary

1. **Download the latest release**:
   ```bash
   mkdir -p /opt/sysmanix
   cd /opt/sysmanix
   wget https://github.com/toxic-development/sysmanix/releases/latest/download/SysManix_linux_amd64 -O sysmanix
   chmod +x sysmanix
   ```

2. **Generate initial configuration**:
   ```bash
   sudo ./sysmanix
   ```

3. **Edit the configuration**:
   ```bash
   sudo nano config.yaml
   # Update settings according to your needs
   ```

### Option 2: Systemd Service (Recommended for Production)

1. **Create a dedicated system user** (optional but recommended):
   ```bash
   sudo useradd -r -s /bin/false sysmanix
   ```

2. **Set up directory structure**:
   ```bash
   sudo mkdir -p /opt/sysmanix
   sudo mkdir -p /etc/sysmanix
   sudo mkdir -p /var/log/sysmanix

   # If you created a dedicated user
   sudo chown sysmanix:sysmanix /var/log/sysmanix
   ```

3. **Install the binary**:
   ```bash
   sudo wget https://github.com/toxic-development/sysmanix/releases/latest/download/SysManix_linux_amd64 -O /opt/sysmanix/sysmanix
   sudo chmod +x /opt/sysmanix/sysmanix
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
   User=root  # Must run as root to control services
   WorkingDirectory=/opt/sysmanix
   ExecStart=/opt/sysmanix/sysmanix -config /etc/sysmanix/config.yaml
   Restart=on-failure
   RestartSec=5s
   SyslogIdentifier=sysmanix

   # Hardening options (optional)
   NoNewPrivileges=true
   ProtectSystem=full
   ProtectHome=true
   PrivateTmp=true

   [Install]
   WantedBy=multi-user.target
   ```

5. **Generate and move the configuration**:
   ```bash
   sudo /opt/sysmanix/sysmanix -config /opt/sysmanix/config.yaml
   sudo cp /opt/sysmanix/config.yaml /etc/sysmanix/
   sudo nano /etc/sysmanix/config.yaml  # Edit as needed
   ```

6. **Enable and start the service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable sysmanix
   sudo systemctl start sysmanix
   ```

## Linux-Specific Configuration

Edit the `config.yaml` file with Linux-specific settings:

```yaml
linux:
  serviceCommand: "systemctl"  # Use systemctl for systemd
  logDirectory: "/var/log/sysmanix"
  services:
    protected:
      # Default protected services (do not modify unless you know what you're doing)
      - systemd
      - systemd-journald
      - dbus
      - sshd
      # Additional custom protected services
      - your-critical-service
```

## Firewall Configuration

If you're using a firewall (which is recommended), you'll need to open the API port:

### UFW (Ubuntu/Debian)
```bash
# Open the default port
sudo ufw allow 40200/tcp

# Alternatively, limit access to specific IP addresses
sudo ufw allow from 192.168.1.0/24 to any port 40200 proto tcp
```

### Firewalld (CentOS/RHEL/Fedora)
```bash
# Open the default port
sudo firewall-cmd --permanent --add-port=40200/tcp
sudo firewall-cmd --reload

# Alternatively, create a service definition
sudo firewall-cmd --permanent --new-service=sysmanix
sudo firewall-cmd --permanent --service=sysmanix --add-port=40200/tcp
sudo firewall-cmd --permanent --add-service=sysmanix
sudo firewall-cmd --reload
```

## SELinux Configuration

If SELinux is enabled on your system (common on CentOS/RHEL), you'll need to configure it to allow SysManix to function properly:

```bash
# Check SELinux status
getenforce

# If Enforcing, create a policy for SysManix
sudo semanage port -a -t http_port_t -p tcp 40200
```

You may also need to set appropriate contexts for SysManix files:

```bash
# Set context for binary
sudo chcon -t bin_t /opt/sysmanix/sysmanix

# Set context for configuration directory
sudo chcon -R -t etc_t /etc/sysmanix
```

## Service Management Specifics

### Linux Service Names

When managing Linux services via SysManix, use the exact service name as shown in systemd:

```bash
# List all services
systemctl list-units --type=service --all

# Example service names
# - nginx.service
# - postgresql.service
# - apache2.service
```

Note that SysManix will automatically append `.service` if it's not included in your request.

### Service Logs Retrieval

On Linux, SysManix retrieves logs using `journalctl`. By default, it:
- Retrieves logs for the specified unit
- Limits log lines to the requested number
- Parses the journal format into structured log entries

## Configuring HTTPS with Nginx Reverse Proxy

For production environments, you should secure SysManix with HTTPS using Nginx:

1. **Install Nginx**:
   ```bash
   sudo apt install nginx   # Debian/Ubuntu
   # or
   sudo yum install nginx   # CentOS/RHEL
   ```

2. **Create Nginx configuration**:
   ```bash
   sudo nano /etc/nginx/sites-available/sysmanix
   ```

   Add this configuration:
   ```nginx
   server {
       listen 443 ssl;
       server_name sysmanix.example.com;

       ssl_certificate /etc/nginx/ssl/sysmanix.crt;
       ssl_certificate_key /etc/nginx/ssl/sysmanix.key;
       ssl_protocols TLSv1.2 TLSv1.3;
       ssl_prefer_server_ciphers on;

       location / {
           proxy_pass http://localhost:40200;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }

   # Redirect HTTP to HTTPS
   server {
       listen 80;
       server_name sysmanix.example.com;
       return 301 https://$server_name$request_uri;
   }
   ```

3. **Generate SSL certificate** (or use Let's Encrypt):
   ```bash
   sudo mkdir -p /etc/nginx/ssl
   sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout /etc/nginx/ssl/sysmanix.key \
     -out /etc/nginx/ssl/sysmanix.crt
   ```

4. **Enable and restart Nginx**:
   ```bash
   sudo ln -s /etc/nginx/sites-available/sysmanix /etc/nginx/sites-enabled/
   sudo systemctl restart nginx
   ```

## Advanced Linux Configuration

### Limiting Resource Usage

Use systemd resource controls in the service file:

```ini
[Service]
# Add these lines to your existing service file
CPUQuota=50%
MemoryLimit=256M
TasksMax=100
```

### Journal Storage Configuration

Optimize journald configuration to manage SysManix logs better:

```bash
sudo nano /etc/systemd/journald.conf
```

Common settings to adjust:
```
SystemMaxUse=1G         # Maximum disk space used by journals
MaxRetentionSec=2week   # Maximum retention time
```

Then restart journald:
```bash
sudo systemctl restart systemd-journald
```

### Monitoring Integration

For monitoring SysManix on Linux, you can use standard tools:

#### Prometheus Monitoring

1. Create a node_exporter systemd service:
   ```bash
   sudo nano /etc/systemd/system/node-exporter.service
   ```

2. Add configuration:
   ```ini
   [Unit]
   Description=Node Exporter
   After=network.target

   [Service]
   User=node_exporter
   ExecStart=/usr/local/bin/node_exporter --collector.systemd

   [Install]
   WantedBy=multi-user.target
   ```

3. Configure Prometheus to scrape SysManix metrics (if SysManix exposes them).

## Troubleshooting Linux-Specific Issues

### Common Issues and Solutions

1. **"Permission denied" when accessing services**:
   - Ensure SysManix is running as root
   - Check if SELinux is blocking access with: `sudo ausearch -m avc -ts recent`

2. **Service control fails**:
   - Verify systemd service exists: `systemctl list-unit-files | grep your-service`
   - Check for proper permissions: `systemctl status your-service`

3. **Logs not appearing**:
   - Check if the user can access journald logs: `sudo usermod -a -G systemd-journal your-user`
   - Verify journal persistence is enabled: `grep Storage /etc/systemd/journald.conf`

### Systemd Journal Access

To manually view logs for debugging:

```bash
# View SysManix logs
sudo journalctl -u sysmanix

# Follow logs in real time
sudo journalctl -u sysmanix -f

# View logs for a specific service
sudo journalctl -u nginx -n 100
```

## Performance Tuning

For optimal performance on Linux:

1. **Adjust open file limits**:
   ```bash
   # Add to /etc/security/limits.conf
   sysmanix soft nofile 65536
   sysmanix hard nofile 65536
   ```

2. **Enable systemd service optimization**:
   ```ini
   # Add to the [Service] section in sysmanix.service
   LimitNOFILE=65536
   Nice=-5
   IOSchedulingClass=2
   IOSchedulingPriority=0
   ```

3. **Configure system-wide optimizations**:
   ```bash
   # Improve network performance in /etc/sysctl.conf
   net.ipv4.tcp_fin_timeout = 30
   net.core.somaxconn = 1024
   ```

## Next Steps

After completing the Linux setup:

- [Secure your installation](./SECURITY.md)
- [Configure authentication](./AUTHENTICATION.md)
- [Create systemd service files for backup and monitoring](./SYSTEMD_SETUP.md)
- [Troubleshoot common issues](./TROUBLESHOOTING.md)
