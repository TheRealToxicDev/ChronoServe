# SysManix Linux Setup Guide

This guide provides detailed instructions for setting up SysManix on Linux systems for production use.

## Prerequisites

- A Linux distribution with systemd (Ubuntu 20.04+, CentOS 8+, Debian 11+, etc.)
- Root or sudo access
- Basic familiarity with Linux command line
- SysManix installed (see [Installation Guide](./INSTALLATION.md))

## Installation Review

If you haven't installed SysManix yet, follow the [Installation Guide](./INSTALLATION.md). In summary:

### Package Manager Installation

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

## Linux System Preparation

### Creating a Dedicated User

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

### Directory Structure

The standard directory structure for SysManix on Linux:

```
/etc/sysmanix/          # Configuration files
  └── config.yaml       # Main configuration file
/usr/bin/sysmanix       # Binary executable
/var/lib/sysmanix/      # Application data
/var/log/sysmanix/      # Log files
  ├── app.log           # Application logs
  ├── auth.log          # Authentication logs
  └── access.log        # API access logs
```

## Configuration Setup

Create and edit the configuration file:

```bash
# Create default configuration
sudo mkdir -p /etc/sysmanix
sudo touch /etc/sysmanix/config.yaml
sudo chown -R sysmanix:sysmanix /etc/sysmanix
sudo chmod 750 /etc/sysmanix
sudo chmod 640 /etc/sysmanix/config.yaml

# Edit configuration
sudo nano /etc/sysmanix/config.yaml
```

Add the following configuration, adjusted for your environment:

```yaml
server:
  host: "0.0.0.0"  # Listen on all interfaces
  port: 40200
  readTimeout: "15s"
  writeTimeout: "15s"
  maxHeaderBytes: 1048576

auth:
  secretKey: "generate-a-secure-random-string"  # Change this!
  tokenDuration: 8h
  issuedBy: "SysManix"
  allowedRoles:
    - admin
    - viewer
  users:
    admin:
      username: "admin"
      password: "your-secure-admin-password"  # Will be hashed after first run
      roles:
        - admin
    viewer:
      username: "viewer"
      password: "your-secure-viewer-password"  # Will be hashed after first run
      roles:
        - viewer

logging:
  level: "info"
  directory: "/var/log/sysmanix"
  maxSize: 10
  maxBackups: 5
  maxAge: 30
  compress: true

linux:
  serviceCommand: "systemctl"
  logDirectory: "/var/log/sysmanix"
  services:
    protected:
      - systemd
      - systemd-journald
      - dbus
      - sshd
      - NetworkManager
      - firewalld
      - ufw
```

### Generate a Secure Secret Key

```bash
# Generate a secure random string for JWT signing
SECRET=$(openssl rand -base64 32)
echo "Generated secret key: $SECRET"
```

Update the `secretKey` field in your configuration with this value.

## Setting Up Systemd Integration

For production use, you should run SysManix as a systemd service. See the [Systemd Setup Guide](./SYSTEMD_SETUP.md) for detailed instructions. In summary:

1. Create a systemd service file at `/etc/systemd/system/sysmanix.service`
2. Configure it with appropriate security settings
3. Start and enable the service

## Firewall Configuration

### UFW (Uncomplicated Firewall)

If using UFW (default on Ubuntu/Debian):

```bash
# Allow connections to SysManix API from internal networks only
sudo ufw allow from 10.0.0.0/8 to any port 40200 proto tcp
sudo ufw allow from 172.16.0.0/12 to any port 40200 proto tcp
sudo ufw allow from 192.168.0.0/16 to any port 40200 proto tcp
```

### firewalld (CentOS/RHEL/Fedora)

If using firewalld:

```bash
# Allow connections to SysManix API from internal networks
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="10.0.0.0/8" port protocol="tcp" port="40200" accept'
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="172.16.0.0/12" port protocol="tcp" port="40200" accept'
sudo firewall-cmd --permanent --add-rich-rule='rule family="ipv4" source address="192.168.0.0/16" port protocol="tcp" port="40200" accept'
sudo firewall-cmd --reload
```

### iptables

If directly using iptables:

```bash
# Allow connections to SysManix API from internal networks
sudo iptables -A INPUT -p tcp -s 10.0.0.0/8 --dport 40200 -j ACCEPT
sudo iptables -A INPUT -p tcp -s 172.16.0.0/12 --dport 40200 -j ACCEPT
sudo iptables -A INPUT -p tcp -s 192.168.0.0/16 --dport 40200 -j ACCEPT

# Optionally save iptables rules for persistence
sudo netfilter-persistent save
```

## SELinux Configuration

If SELinux is enabled on your system (common in RHEL/CentOS), you'll need to configure it to work with SysManix:

```bash
# Allow SysManix to bind to its network port
sudo semanage port -a -t http_port_t -p tcp 40200

# Allow SysManix to access systemd services
sudo setsebool -P httpd_manage_sys_content 1

# Create a custom SELinux policy module for SysManix
sudo ausearch -c 'sysmanix' --raw | audit2allow -M sysmanix_selinux
sudo semodule -i sysmanix_selinux.pp
```

## System Limits Configuration

Optimize system limits for SysManix in `/etc/security/limits.conf`:

```bash
# Add the following lines to /etc/security/limits.conf
sudo tee -a /etc/security/limits.conf << EOF
sysmanix soft nofile 65536
sysmanix hard nofile 65536
sysmanix soft nproc 4096
sysmanix hard nproc 4096
EOF
```

## Monitoring Integration

### Setting up Prometheus Monitoring

To monitor SysManix with Prometheus:

1. Install node_exporter on the server
2. Configure Prometheus to scrape node_exporter metrics
3. Create a script to expose SysManix metrics via the health endpoint

```bash
# Example script to expose SysManix metrics to Prometheus
sudo tee /usr/local/bin/sysmanix-metrics.sh << 'EOF'
#!/bin/bash
set -e

# Configuration
API_URL="http://localhost:40200/health"
OUTPUT_FILE="/var/lib/node_exporter/sysmanix.prom"

# Get SysManix health data as JSON
HEALTH_DATA=$(curl -s ${API_URL})

# Extract metrics and format for Prometheus
echo "# HELP sysmanix_uptime_seconds SysManix uptime in seconds" > ${OUTPUT_FILE}
echo "# TYPE sysmanix_uptime_seconds gauge" >> ${OUTPUT_FILE}
UPTIME_SEC=$(echo ${HEALTH_DATA} | jq -r '.data.uptime_seconds')
echo "sysmanix_uptime_seconds ${UPTIME_SEC}" >> ${OUTPUT_FILE}

echo "# HELP sysmanix_memory_usage_bytes Memory usage in bytes" >> ${OUTPUT_FILE}
echo "# TYPE sysmanix_memory_usage_bytes gauge" >> ${OUTPUT_FILE}
MEMORY_ALLOC=$(echo ${HEALTH_DATA} | jq -r '.data.memory.alloc')
echo "sysmanix_memory_usage_bytes ${MEMORY_ALLOC}" >> ${OUTPUT_FILE}

echo "# HELP sysmanix_gc_count Number of garbage collections" >> ${OUTPUT_FILE}
echo "# TYPE sysmanix_gc_count counter" >> ${OUTPUT_FILE}
GC_COUNT=$(echo ${HEALTH_DATA} | jq -r '.data.memory.numGC')
echo "sysmanix_gc_count ${GC_COUNT}" >> ${OUTPUT_FILE}
EOF

# Make script executable
sudo chmod +x /usr/local/bin/sysmanix-metrics.sh

# Create cron job to run the script every minute
sudo tee /etc/cron.d/sysmanix-metrics << EOF
* * * * * root /usr/local/bin/sysmanix-metrics.sh > /dev/null 2>&1
EOF
```

## Log Rotation

Configure logrotate for SysManix logs:

```bash
# Create logrotate configuration
sudo tee /etc/logrotate.d/sysmanix << EOF
/var/log/sysmanix/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 sysmanix sysmanix
    sharedscripts
    postrotate
        systemctl kill -s HUP sysmanix.service
    endscript
}
EOF
```

## Automatic Updates

Set up automatic security updates for the system:

### Ubuntu/Debian

```bash
# Install unattended-upgrades
sudo apt install unattended-upgrades apt-listchanges

# Configure automatic updates
sudo tee /etc/apt/apt.conf.d/50unattended-upgrades << EOF
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
};
Unattended-Upgrade::Package-Blacklist {
};
Unattended-Upgrade::Automatic-Reboot "false";
EOF

# Enable automatic updates
sudo tee /etc/apt/apt.conf.d/20auto-upgrades << EOF
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
EOF
```

### CentOS/RHEL

```bash
# Install dnf-automatic
sudo yum install dnf-automatic

# Configure automatic updates
sudo sed -i 's/apply_updates = no/apply_updates = yes/' /etc/dnf/automatic.conf

# Enable and start the service
sudo systemctl enable --now dnf-automatic.timer
```

## Performance Tuning

### Optimizing for High Traffic

For servers with many connections, optimize kernel parameters:

```bash
# Add the following to /etc/sysctl.conf
sudo tee -a /etc/sysctl.conf << EOF
# Increase TCP max connections
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 65535

# Increase file descriptor limits
fs.file-max = 2097152

# Optimize TCP keepalive parameters
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 60
net.ipv4.tcp_keepalive_probes = 5
EOF

# Apply changes
sudo sysctl -p
```

## Backup and Recovery

Set up automated backups of SysManix configuration:

```bash
# Create backup script
sudo tee /usr/local/bin/backup-sysmanix.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/var/backups/sysmanix"
DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/sysmanix-config-${DATE}.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p ${BACKUP_DIR}

# Create backup
tar -czf ${BACKUP_FILE} -C /etc sysmanix

# Rotate backups (keep last 7 days)
find ${BACKUP_DIR} -type f -name "sysmanix-config-*.tar.gz" -mtime +7 -delete
EOF

# Make script executable
sudo chmod +x /usr/local/bin/backup-sysmanix.sh

# Create daily cron job
sudo tee /etc/cron.daily/sysmanix-backup << EOF
#!/bin/bash
/usr/local/bin/backup-sysmanix.sh
EOF

# Make cron job executable
sudo chmod +x /etc/cron.daily/sysmanix-backup
```

## Security Hardening

Additional security measures for your Linux system:

```bash
# Disable root SSH login
sudo sed -i 's/^PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Restrict SSH to strong encryption algorithms
sudo tee -a /etc/ssh/sshd_config << EOF
Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com,aes256-ctr,aes192-ctr,aes128-ctr
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com,hmac-sha2-512,hmac-sha2-256
KexAlgorithms curve25519-sha256@libssh.org,diffie-hellman-group-exchange-sha256
EOF
sudo systemctl restart sshd

# Enable process accounting
sudo apt install acct || sudo yum install psacct
sudo systemctl enable --now acct.service || sudo systemctl enable --now psacct.service
```

## Troubleshooting Linux-Specific Issues

### Service Management Issues

If SysManix has trouble managing systemd services:

1. Check that the sysmanix user has appropriate permissions:

```bash
# Create a sudoers file for SysManix
sudo tee /etc/sudoers.d/sysmanix << EOF
sysmanix ALL=(ALL) NOPASSWD: /bin/systemctl start *, /bin/systemctl stop *, /bin/systemctl status *, /bin/systemctl is-active *, /usr/bin/journalctl *
EOF

# Set proper permissions
sudo chmod 440 /etc/sudoers.d/sysmanix
```

2. Update the systemd service file to use sudo for service commands:

```bash
sudo systemctl edit sysmanix.service
```

Add the following:

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/sysmanix serve --use-sudo
```

Then reload and restart:

```bash
sudo systemctl daemon-reload
sudo systemctl restart sysmanix
```

### Journalctl Access Issues

If SysManix can't access journalctl logs:

```bash
# Add sysmanix user to systemd-journal group
sudo usermod -a -G systemd-journal sysmanix

# Give the systemd-journal group permission to read journal files
sudo chmod -R g+r /var/log/journal
```

## System Maintenance Procedures

### Rolling Restarts

For updating SysManix with minimal downtime in clustered environments:

```bash
# Script for rolling restart (/usr/local/bin/sysmanix-rolling-restart.sh)
#!/bin/bash
# This script performs a rolling restart of SysManix instances

# Load balancer hostnames - replace with your actual load balancer hostnames
LOAD_BALANCERS=("lb1.example.com" "lb2.example.com")

# SysManix server hostnames - replace with your actual server hostnames
SERVERS=("sysmanix1.example.com" "sysmanix2.example.com" "sysmanix3.example.com")

for SERVER in "${SERVERS[@]}"; do
  echo "Processing server: $SERVER"

  # Remove from load balancer
  for LB in "${LOAD_BALANCERS[@]}"; do
    echo "Removing $SERVER from $LB"
    ssh $LB "sudo lb-disable-backend $SERVER"
  done

  # Wait for connections to drain
  echo "Waiting for connections to drain (30s)..."
  sleep 30

  # Update and restart SysManix
  echo "Updating and restarting SysManix on $SERVER"
  ssh $SERVER "sudo systemctl restart sysmanix"

  # Wait for service to be fully up
  echo "Waiting for service to start (10s)..."
  sleep 10

  # Verify service is healthy
  HEALTH=$(ssh $SERVER "curl -s http://localhost:40200/health | jq -r '.data.status'")
  if [ "$HEALTH" != "healthy" ]; then
    echo "ERROR: Service on $SERVER is not healthy!"
    exit 1
  fi

  # Add back to load balancer
  for LB in "${LOAD_BALANCERS[@]}"; do
    echo "Adding $SERVER back to $LB"
    ssh $LB "sudo lb-enable-backend $SERVER"
  done

  echo "Server $SERVER successfully updated"
  echo "Waiting before proceeding to next server (10s)..."
  sleep 10
done

echo "All servers updated successfully"
```

## Additional Resources

- [Linux Hardening Guide](https://github.com/trimstray/the-practical-linux-hardening-guide)
- [Systemd Documentation](https://systemd.io/)
- [Linux Performance Tuning Guide](https://www.kernel.org/doc/Documentation/sysctl/vm.txt)

## Next Steps

Once you have set up SysManix on your Linux system:

1. Set up [NGINX as a reverse proxy](./NGINX_SETUP.md) for TLS termination
2. Configure [monitoring and alerts](./MONITORING.md) for proactive management
3. Implement [backup and disaster recovery](./DISASTER_RECOVERY.md) procedures
4. Review the [Security Guide](./SECURITY.md) for additional security hardening steps
