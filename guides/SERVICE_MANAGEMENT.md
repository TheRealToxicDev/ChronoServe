# SysManix Service Management Guide

This guide explains how to manage system services using SysManix.

## Configuration

Ensure your `config.yaml` file is properly set up:

```yaml
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
```

## Service Management Overview

SysManix provides a unified API for managing services across different operating systems, abstracting away the platform-specific details. The API allows you to:

- List all available services
- Check service status
- Start and stop services
- View service logs

## Authentication

Before performing any service management operations, you need to authenticate and obtain a JWT token:

```bash
# Using curl
curl -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

```powershell
# Using PowerShell
$auth = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" -Method Post -Body $auth -ContentType "application/json"
$token = $response.data.token
$headers = @{ "Authorization" = "Bearer $token" }
```

## Listing Services

To get a list of all available services:

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

The response includes basic information about each service:

```json
{
  "status": "success",
  "message": "Services retrieved successfully",
  "data": [
    {
      "name": "wuauserv",
      "displayName": "Windows Update",
      "isActive": true,
      "updatedAt": "2025-03-06T21:55:34.123Z"
    },
    {
      "name": "nginx",
      "displayName": "nginx",
      "isActive": false,
      "updatedAt": "2025-03-06T21:55:34.123Z"
    }
  ]
}
```

## Checking Service Status

To check the status of a specific service:

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services/status/{service-name}
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/services/status/wuauserv" -Headers $headers
```

The response includes detailed status information:

```json
{
  "status": "success",
  "message": "Service status retrieved successfully",
  "data": {
    "name": "wuauserv",
    "status": "Running",
    "isActive": true,
    "updatedAt": "2025-03-06T21:56:34.567Z"
  }
}
```

## Starting a Service

To start a service:

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:40200/services/start/{service-name}
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/services/start/wuauserv" -Method Post -Headers $headers
```

The response confirms the action:

```json
{
  "status": "success",
  "message": "Service wuauserv started successfully"
}
```

## Stopping a Service

To stop a service:

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:40200/services/stop/{service-name}
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/services/stop/wuauserv" -Method Post -Headers $headers
```

The response confirms the action:

```json
{
  "status": "success",
  "message": "Service wuauserv stopped successfully"
}
```

## Viewing Service Logs

To view logs for a service:

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services/logs/{service-name}?lines=50
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv?lines=50" -Headers $headers
```

The response includes log entries:

```json
{
  "status": "success",
  "message": "Service logs retrieved successfully",
  "data": [
    {
      "time": "2025-03-06 21:50:12",
      "level": "Information",
      "message": "Service started successfully"
    },
    {
      "time": "2025-03-06 21:55:30",
      "level": "Information",
      "message": "Service stopped by user request"
    }
  ]
}
```

## Protected Services

SysManix includes protection for critical system services. Attempts to start or stop these services will result in a `403 Forbidden` response:

```json
{
  "status": "error",
  "message": "operation not allowed on protected system service: lsass",
  "code": 403
}
```

### Windows Protected Services

The following Windows services are protected:
- `wininit`: Windows Start-Up Application
- `csrss`: Client Server Runtime Process
- `lsass`: Local Security Authority Process
- `services`: Services and Controller App
- `spooler`: Print Spooler
- and other critical system services

### Linux Protected Services

The following Linux services are protected:
- `systemd`: Core system daemon
- `systemd-journald`: Journal logging service
- `dbus`: System message bus
- `sshd`: SSH daemon
- and other critical system services

## Role-Based Access

Service management is controlled by role-based access:

- **viewer role**: Can list services, view status, and view logs
- **admin role**: Can perform all operations including starting and stopping services

## Platform-Specific Details

### Windows Implementation

On Windows, SysManix uses PowerShell commands to interact with services:

- List services: Uses `Get-Service` cmdlet
- Start services: Uses `Start-Service` cmdlet
- Stop services: Uses `Stop-Service` cmdlet
- View logs: Uses `Get-WinEvent` with filtering

### Linux Implementation

On Linux, SysManix uses systemd commands:

- List services: Uses `systemctl list-units --type=service`
- Start services: Uses `systemctl start <service>`
- Stop services: Uses `systemctl stop <service>`
- View logs: Uses `journalctl -u <service>`

## Error Handling

SysManix uses consistent error responses:

```json
{
  "status": "error",
  "message": "Failed to start service: The service did not respond in a timely fashion",
  "code": 500
}
```

Common error codes:
- `400`: Bad request (invalid input)
- `401`: Unauthorized (missing or invalid token)
- `403`: Forbidden (insufficient permissions)
- `404`: Not found (service doesn't exist)
- `408`: Request timeout (service operation took too long)
- `500`: Internal server error

## Practical Examples

### Service Monitoring Script

```bash
#!/bin/bash
# Service monitoring script using SysManix API

# Configuration
API_URL="http://localhost:40200"
TOKEN="your-token-here"
SERVICES=("nginx" "postgresql" "mysql")

# Check each service
for service in "${SERVICES[@]}"; do
  echo "Checking $service..."
  status=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/services/status/$service")
  isActive=$(echo $status | grep -o '"isActive":true' | wc -l)

  if [ $isActive -eq 0 ]; then
    echo "Service $service is not running! Starting..."
    curl -s -X POST -H "Authorization: Bearer $TOKEN" "$API_URL/services/start/$service"
    echo "Service start requested."
  else
    echo "Service $service is running normally."
  fi
done
```

### PowerShell Service Dashboard

```powershell
# PowerShell Service Dashboard using SysManix API
param(
    [string]$ApiUrl = "http://localhost:40200",
    [string]$Username = "admin",
    [string]$Password = "your-password"
)

# Authenticate
$authBody = @{
    username = $Username
    password = $Password
} | ConvertTo-Json

$authResult = Invoke-RestMethod -Uri "$ApiUrl/auth/login" -Method Post -Body $authBody -ContentType "application/json"
$token = $authResult.data.token
$headers = @{ "Authorization" = "Bearer $token" }

# Get services
$services = Invoke-RestMethod -Uri "$ApiUrl/services" -Headers $headers
$activeCount = ($services.data | Where-Object { $_.isActive -eq $true }).Count
$inactiveCount = ($services.data | Where-Object { $_.isActive -eq $false }).Count

Write-Host "====== SysManix Service Dashboard ======"
Write-Host "Total services: $($services.data.Count)"
Write-Host "Active services: $activeCount"
Write-Host "Inactive services: $inactiveCount"
Write-Host ""
Write-Host "Service Status:"

$services.data | Format-Table -Property name, displayName, isActive -AutoSize
```

## Best Practices

1. **Token Management**:
   - Store tokens securely
   - Use token refresh for long-running applications
   - Set appropriate token expiration times

2. **Error Handling**:
   - Implement retry logic for transient errors
   - Handle 408 Request Timeout errors by checking status afterward
   - Log failed operations for troubleshooting

3. **Performance Considerations**:
   - Cache service status where appropriate
   - Limit log retrieval to necessary lines
   - Implement pagination for large service lists

## Troubleshooting

### Common Issues

1. **Service operation takes too long**
   - SysManix has a timeout for service operations (10 seconds by default)
   - Check the service status after the timeout to verify completion

2. **Service not found**
   - Verify the service name is correct
   - Different operating systems may use different service names

3. **Insufficient permissions**
   - Ensure your user has the required role
   - Admin role is needed for start/stop operations

## Next Steps

- Learn about [Authentication and Token Management](./AUTHENTICATION.md)
- Configure [HTTPS with Nginx](./NGINX_SETUP.md) for production use
- Explore [Systemd Integration](./SYSTEMD_SETUP.md) for Linux deployments
