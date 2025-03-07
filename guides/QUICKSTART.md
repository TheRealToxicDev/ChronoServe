# SysManix Quick Start Guide

This guide will help you quickly set up and start using SysManix for service management on both Windows and Linux systems.

## Installation in 5 Minutes

### Windows Quick Install

1. Download the latest Windows installer from the [releases page](https://github.com/toxic-development/SysManix/releases)
2. Run the installer as administrator
3. Follow the installation wizard, accepting defaults
4. SysManix will start automatically as a Windows service

### Linux Quick Install

```bash
# Ubuntu/Debian
curl -fsSL https://download.sysmanix.io/install.sh | sudo bash

# RHEL/CentOS/Fedora
curl -fsSL https://download.sysmanix.io/install.sh | sudo bash
```

## First-Time Configuration

### Accessing the Configuration

#### Windows
```powershell
# Open the configuration file
notepad "C:\Program Files\SysManix\config.yaml"
```

#### Linux
```bash
# Open the configuration file
sudo nano /etc/sysmanix/config.yaml
```

### Essential Security Settings

Change these security-sensitive settings immediately:

```yaml
auth:
  secretKey: "generate-a-secure-random-string-here"
  users:
    admin:
      username: "admin"
      password: "your-secure-admin-password"  # Will be hashed after first run
    viewer:
      username: "viewer"
      password: "your-secure-viewer-password"  # Will be hashed after first run
```

Save the changes and restart SysManix:

#### Windows
```powershell
Restart-Service -Name SysManix
```

#### Linux
```bash
sudo systemctl restart sysmanix
```

## Getting Authentication Token

Before using the API, you need to authenticate and get a JWT token:

### Using PowerShell (Windows)

```powershell
$loginBody = @{
    username = "admin"
    password = "your-secure-admin-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
                              -Method Post `
                              -ContentType "application/json" `
                              -Body $loginBody

# Store the token for later use
$token = $response.data.token
$headers = @{
    "Authorization" = "Bearer $token"
}
```

### Using cURL (Linux/macOS)

```bash
# Login and extract token
TOKEN=$(curl -s -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-secure-admin-password"}' \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"\(.*\)"/\1/')

# Store headers for later use
HEADERS="Authorization: Bearer $TOKEN"
```

## Basic API Operations

### List All Services

#### Windows (PowerShell)
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

#### Linux (cURL)
```bash
curl -H "$HEADERS" http://localhost:40200/services
```

### Check Service Status

#### Windows (PowerShell)
```powershell
# For Windows Update service
Invoke-RestMethod -Uri "http://localhost:40200/services/status/wuauserv" -Headers $headers
```

#### Linux (cURL)
```bash
# For SSH service
curl -H "$HEADERS" http://localhost:40200/services/status/sshd
```

### View Service Logs

#### Windows (PowerShell)
```powershell
# Get the last 50 log entries for Windows Update service
Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv?lines=50" -Headers $headers
```

#### Linux (cURL)
```bash
# Get the last 50 log entries for SSH service
curl -H "$HEADERS" "http://localhost:40200/services/logs/sshd?lines=50"
```

### Start a Service

#### Windows (PowerShell)
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services/start/wuauserv" `
                 -Method Post `
                 -Headers $headers
```

#### Linux (cURL)
```bash
curl -X POST -H "$HEADERS" http://localhost:40200/services/start/nginx
```

### Stop a Service

#### Windows (PowerShell)
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services/stop/wuauserv" `
                 -Method Post `
                 -Headers $headers
```

#### Linux (cURL)
```bash
curl -X POST -H "$HEADERS" http://localhost:40200/services/stop/nginx
```

## Common Windows Services

Here are some common Windows services you can manage with SysManix:

| Service Name | Display Name | Description |
|--------------|--------------|-------------|
| wuauserv | Windows Update | Manages Windows Update |
| EventLog | Windows Event Log | Logs system events |
| WinRM | Windows Remote Management | Enables remote management |
| BITS | Background Intelligent Transfer | Manages file transfers |
| LanmanServer | Server | Enables file and printer sharing |
| Spooler | Print Spooler | Manages printer jobs |

## Common Linux Services

Here are some common Linux services you can manage with SysManix:

| Service Name | Description |
|--------------|-------------|
| nginx | Web server |
| apache2 | Web server |
| sshd | SSH server |
| postgresql | PostgreSQL database |
| mysql | MySQL database |
| docker | Docker container service |

## Health Check

Verify that SysManix is running properly:

#### Windows (PowerShell)
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/health"
```

#### Linux (cURL)
```bash
curl http://localhost:40200/health
```

You should see a response with the status "healthy" and various system metrics.

## Troubleshooting Common Issues

### Can't Connect to API

1. Verify SysManix is running:
   ```powershell
   # Windows
   Get-Service -Name SysManix
   ```
   ```bash
   # Linux
   systemctl status sysmanix
   ```

2. Check the firewall settings:
   ```powershell
   # Windows
   Get-NetFirewallRule | Where-Object { $_.DisplayName -like "*SysManix*" }
   ```
   ```bash
   # Linux
   sudo ufw status
   ```

### Authentication Failures

1. Verify the credentials in your request match the configured user credentials
2. Check if the token has expired (default lifetime is 24 hours)
3. Ensure the configuration file has been properly saved and the service restarted

### Permission Denied

1. Check that your user has the necessary role for the operation:
   - `viewer`: Can list services, view status and logs
   - `admin`: Can perform all operations including start/stop services

2. Verify you're not trying to modify a protected service

## Next Steps

Now that you've set up SysManix and executed basic operations, consider:

1. Setting up [HTTPS with Nginx](./NGINX_SETUP.md)
2. Configuring [Windows Service](./WINDOWS_SETUP.md) or [Systemd](./SYSTEMD_SETUP.md) for production use
3. Implementing [security best practices](./SECURITY.md)
4. Exploring advanced [service management capabilities](./SERVICE_MANAGEMENT.md)
5. Learning about [authentication options](./AUTHENTICATION.md)

## Sample Scripts

### Windows Service Dashboard (PowerShell)

Save the following script as `SysManixDashboard.ps1`:

```powershell
# SysManix Service Dashboard
param (
    [string]$username = "admin",
    [string]$password = "your-secure-admin-password",
    [string]$baseUrl = "http://localhost:40200"
)

# Get authentication token
$loginBody = @{
    username = $username
    password = $password
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -ContentType "application/json" -Body $loginBody
    $token = $response.data.token
    $headers = @{
        "Authorization" = "Bearer $token"
    }

    Write-Host "Authentication successful" -ForegroundColor Green
} catch {
    Write-Host "Authentication failed: $_" -ForegroundColor Red
    exit 1
}

# Get list of services
try {
    $services = Invoke-RestMethod -Uri "$baseUrl/services" -Headers $headers

    # Display services in a formatted table
    Write-Host "Services:" -ForegroundColor Cyan
    $services.data | Format-Table -Property Name, DisplayName, Status -AutoSize

    # Offer action menu
    $selectedService = Read-Host "Enter service name to manage (or 'exit' to quit)"

    while ($selectedService -ne "exit") {
        Write-Host "Actions for $selectedService:" -ForegroundColor Yellow
        Write-Host "1. View Status"
        Write-Host "2. View Logs"
        Write-Host "3. Start Service"
        Write-Host "4. Stop Service"
        Write-Host "5. Back to services list"

        $action = Read-Host "Select action (1-5)"

        switch ($action) {
            1 {
                $status = Invoke-RestMethod -Uri "$baseUrl/services/status/$selectedService" -Headers $headers
                $status.data | Format-List
            }
            2 {
                $lines = Read-Host "Number of log lines to view (default: 20)"
                if ([string]::IsNullOrEmpty($lines)) { $lines = 20 }
                $logs = Invoke-RestMethod -Uri "$baseUrl/services/logs/$selectedService?lines=$lines" -Headers $headers
                $logs.data | Format-Table -Property Time, Level, Message -AutoSize -Wrap
            }
            3 {
                $confirm = Read-Host "Are you sure you want to start $selectedService? (y/n)"
                if ($confirm -eq "y") {
                    $result = Invoke-RestMethod -Uri "$baseUrl/services/start/$selectedService" -Method Post -Headers $headers
                    Write-Host $result.message -ForegroundColor Green
                }
            }
            4 {
                $confirm = Read-Host "Are you sure you want to stop $selectedService? (y/n)"
                if ($confirm -eq "y") {
                    $result = Invoke-RestMethod -Uri "$baseUrl/services/stop/$selectedService" -Method Post -Headers $headers
                    Write-Host $result.message -ForegroundColor Green
                }
            }
            5 {
                break
            }
            default {
                Write-Host "Invalid option" -ForegroundColor Red
            }
        }

        if ($action -ne "5") {
            Read-Host "Press Enter to continue"
        }

        if ($action -eq "5") {
            $services = Invoke-RestMethod -Uri "$baseUrl/services" -Headers $headers
            $services.data | Format-Table -Property Name, DisplayName, Status -AutoSize
            $selectedService = Read-Host "Enter service name to manage (or 'exit' to quit)"
        } else {
            Write-Host "Actions for $selectedService:" -ForegroundColor Yellow
            Write-Host "1. View Status"
            Write-Host "2. View Logs"
            Write-Host "3. Start Service"
            Write-Host "4. Stop Service"
            Write-Host "5. Back to services list"
            $action = Read-Host "Select action (1-5)"
        }
    }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
}
```

### Linux Monitor Script (Bash)

Save the following script as `sysmanix_monitor.sh`:

```bash
#!/bin/bash
# SysManix Service Monitor

# Configuration
USERNAME="admin"
PASSWORD="your-secure-admin-password"
BASE_URL="http://localhost:40200"

# Get authentication token
TOKEN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
  | grep -o '"token":"[^"]*"' | sed 's/"token":"\(.*\)"/\1/')

if [ -z "$TOKEN" ]; then
  echo "Authentication failed"
  exit 1
fi

# Function to check service status
check_service() {
  local service_name=$1
  local expected_status=$2

  STATUS=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/services/status/$service_name" | grep -o '"Status":"[^"]*"' | sed 's/"Status":"\(.*\)"/\1/')

  if [ "$STATUS" != "$expected_status" ]; then
    echo "WARNING: $service_name is $STATUS (expected: $expected_status)"
    return 1
  else
    echo "OK: $service_name is $STATUS"
    return 0
  }
}

# Main monitoring loop
echo "SysManix Service Monitor"
echo "------------------------"

# List of services to monitor: service_name,expected_status
SERVICES=(
  "nginx,Running"
  "postgresql,Running"
  "sshd,Running"
)

for service in "${SERVICES[@]}"; do
  IFS=',' read -r name status <<< "$service"
  check_service "$name" "$status"
done

exit 0
```

## Integrating with Other Tools

### Grafana Dashboard

You can create a Grafana dashboard that displays SysManix health metrics:

1. Set up a script that pulls data from the health endpoint
2. Store the data in a time-series database like InfluxDB
3. Create a Grafana dashboard that visualizes the data

### Jenkins Pipeline

Integrate SysManix with your CI/CD pipeline:

```groovy
pipeline {
    agent any
    stages {
        stage('Deploy Service') {
            steps {
                // Deploy your application
                sh 'deploy_app.sh'

                // Start the service using SysManix
                sh '''
                    TOKEN=$(curl -s -X POST "http://sysmanix-server:40200/auth/login" \
                      -H "Content-Type: application/json" \
                      -d '{"username":"ci-user","password":"ci-password"}' \
                      | jq -r '.data.token')

                    curl -X POST -H "Authorization: Bearer $TOKEN" \
                      "http://sysmanix-server:40200/services/start/my-app-service"
                '''
            }
        }
        stage('Verify Service') {
            steps {
                // Check if service is running
                sh '''
                    TOKEN=$(curl -s -X POST "http://sysmanix-server:40200/auth/login" \
                      -H "Content-Type: application/json" \
                      -d '{"username":"ci-user","password":"ci-password"}' \
                      | jq -r '.data.token')

                    STATUS=$(curl -s -H "Authorization: Bearer $TOKEN" \
                      "http://sysmanix-server:40200/services/status/my-app-service" \
                      | jq -r '.data.Status')

                    if [ "$STATUS" != "Running" ]; then
                      echo "Service failed to start"
                      exit 1
                    fi
                '''
            }
        }
    }
}
```

## Conclusion

You've now completed the SysManix quick start guide and should be able to:

- Install and configure SysManix
- Authenticate and get a JWT token
- Perform basic service management operations
- Understand common issues and how to resolve them

For more detailed information, refer to the other guides in the documentation.
