# SysManix Service Management Guide

## Overview

SysManix provides a secure API for managing system services across Windows and Linux platforms. This guide explains how to use the service management features effectively.

## Platform Support

### Windows Services
- Uses PowerShell commands for service operations
- Supports Windows services accessible via `Get-Service`
- Retrieves service logs from the Windows Event Log

### Linux Services
- Uses systemd for service operations
- Supports services managed by systemctl
- Retrieves logs using journalctl

## API Endpoints

### List Available Services

```http
GET /services

Headers:
Authorization: Bearer your-jwt-token

Response:
{
    "success": true,
    "message": "Services retrieved successfully",
    "data": [
        {
            "Name": "service1",
            "DisplayName": "Service 1 Display Name",
            "Status": "Running"
        },
        ...
    ]
}
```

### Get Service Status

```http
GET /services/status/{service_name}

Headers:
Authorization: Bearer your-jwt-token

Response:
{
    "success": true,
    "message": "Service status retrieved successfully",
    "data": {
        "Name": "service_name",
        "Status": "Running",
        "IsActive": true,
        "UpdatedAt": "2023-10-13T12:34:56Z"
    }
}
```

### Start a Service

```http
POST /services/start/{service_name}

Headers:
Authorization: Bearer your-jwt-token

Response:
{
    "success": true,
    "message": "Service service_name started successfully"
}
```

### Stop a Service

```http
POST /services/stop/{service_name}

Headers:
Authorization: Bearer your-jwt-token

Response:
{
    "success": true,
    "message": "Service service_name stopped successfully"
}
```

### View Service Logs

```http
GET /services/logs/{service_name}?lines=50

Headers:
Authorization: Bearer your-jwt-token

Response:
{
    "success": true,
    "message": "Service logs retrieved successfully",
    "data": [
        {
            "Time": "2023-10-13T12:34:56Z",
            "Level": "Information",
            "Message": "The service entered the running state."
        },
        ...
    ]
}
```

## Protected Services

For system stability and security, certain critical system services are protected from modification:

### Windows Protected Services
- Core system processes: `wininit`, `csrss`, `lsass`, etc.
- Critical system services: `spooler`, `EventLog`, etc.

### Linux Protected Services
- Core daemons: `systemd`, `systemd-journald`, `dbus`, etc.
- System management: `sshd`, `NetworkManager`, etc.

Attempting to start or stop a protected service will result in a 403 Forbidden response:

```json
{
    "success": false,
    "error": "operation not allowed on protected system service: wininit",
    "code": 403
}
```

## Service Operation Behavior

### Windows Services
- Services are started using `Start-Service` PowerShell command
- Services are stopped using `Stop-Service` PowerShell command
- Service status is polled with a 10-second timeout
- Service logs are retrieved from the System Event Log

### Linux Services
- Services are started using `systemctl start` command
- Services are stopped using `systemctl stop` command
- Service status is polled with a 10-second timeout
- Service logs are retrieved using `journalctl` command

## Error Handling

Common service operation errors:

| Error | HTTP Status | Description |
|-------|-------------|-------------|
| Protected service | 403 | Operation attempted on a protected system service |
| Service not found | 404 | The specified service doesn't exist |
| Permission denied | 403 | User lacks permissions for the operation |
| Start/stop timeout | 500 | Service operation didn't complete in time |
| System error | 500 | Underlying system command failed |

## Examples

### PowerShell Examples

```powershell
# Get auth token
$body = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

$token = $response.data.token
$headers = @{
    Authorization = "Bearer $token"
}

# List services
$services = Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers

# Check service status
$status = Invoke-RestMethod -Uri "http://localhost:40200/services/status/wuauserv" -Headers $headers

# Start a service
$result = Invoke-RestMethod -Uri "http://localhost:40200/services/start/wuauserv" `
    -Method Post `
    -Headers $headers

# Stop a service
$result = Invoke-RestMethod -Uri "http://localhost:40200/services/stop/wuauserv" `
    -Method Post `
    -Headers $headers

# View service logs (50 lines)
$logs = Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv?lines=50" -Headers $headers
```

### cURL Examples (Linux/Mac)

```bash
# Get auth token
token=$(curl -s -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}' \
  | jq -r '.data.token')

# List services
curl -H "Authorization: Bearer $token" http://localhost:40200/services

# Check service status
curl -H "Authorization: Bearer $token" http://localhost:40200/services/status/nginx

# Start a service
curl -X POST -H "Authorization: Bearer $token" http://localhost:40200/services/start/nginx

# Stop a service
curl -X POST -H "Authorization: Bearer $token" http://localhost:40200/services/stop/nginx

# View service logs
curl -H "Authorization: Bearer $token" "http://localhost:40200/services/logs/nginx?lines=50"
```

## Best Practices

1. **Service Discovery**: Always list available services first to ensure service names are correct
2. **Error Handling**: Check response for error messages before assuming success
3. **Permissions**: Ensure admin role for start/stop operations
4. **Protected Services**: Avoid attempting operations on system-critical services
5. **Timeout Awareness**: Some services may take longer to start/stop than the 10-second timeout
6. **Log Size**: For large logs, use the `lines` parameter to limit response size
