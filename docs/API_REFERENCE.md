# ChronoServe API Reference

## Base URL
Default: `http://localhost:40200`

## Authentication

All protected endpoints require a JWT token in the Authorization header.

```http
Authorization: Bearer your-jwt-token
```

### Login
```http
POST /auth/login

Request Body:
{
    "username": "admin",
    "password": "your-password"
}

Response (200 OK):
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJ...",
        "roles": ["admin"]
    }
}

Response (401 Unauthorized):
{
    "success": false,
    "error": "Invalid credentials"
}
```

## Service Management

### List Services
```http
GET /services

Headers:
    Authorization: Bearer your-jwt-token

Response (200 OK):
{
    "success": true,
    "message": "Services retrieved successfully",
    "data": [
        {
            "Name": "service_name",
            "DisplayName": "Service Display Name",
            "Status": "Running"
        }
    ]
}
```

### View Service Logs
```http
GET /services/logs/{service_name}

Headers:
    Authorization: Bearer your-jwt-token

Query Parameters:
    lines (optional): Number of lines to return (default: 100)

Response (200 OK):
{
    "success": true,
    "message": "Logs retrieved successfully",
    "data": [
        {
            "Time": "2025-02-28T22:03:26.010Z",
            "Level": "Information",
            "Message": "Service log entry"
        }
    ]
}
```

### Start Service
```http
POST /services/start/{service_name}

Headers:
    Authorization: Bearer your-jwt-token
    
Required Role: admin

Response (200 OK):
{
    "success": true,
    "message": "Service started successfully"
}

Response (403 Forbidden):
{
    "success": false,
    "error": "Insufficient permissions"
}
```

### Stop Service
```http
POST /services/stop/{service_name}

Headers:
    Authorization: Bearer your-jwt-token
    
Required Role: admin

Response (200 OK):
{
    "success": true,
    "message": "Service stopped successfully"
}
```

### Get Service Status
```http
GET /services/status/{service_name}

Headers:
    Authorization: Bearer your-jwt-token

Response (200 OK):
{
    "success": true,
    "message": "Service status retrieved successfully",
    "data": {
        "Name": "service_name",
        "DisplayName": "Service Display Name",
        "Status": "Running"
    }
}
```

## Health Check
```http
GET /health

Response (200 OK):
{
    "success": true,
    "data": {
        "status": "healthy",
        "uptime": "1h2m3s",
        "version": "0.1.0",
        "goVersion": "go1.23.1",
        "memory": {
            "alloc": 1234567,
            "totalAlloc": 2345678,
            "sys": 3456789,
            "numGC": 12,
            "heapInUse": 45.67
        },
        "startTime": "2025-02-28T22:00:00Z"
    }
}
```

## Error Responses
```http
{
    "success": false,
    "error": "Error description"
}
```

## Role-Based Access
- `admin`: Full access to all endpoints
- `viewer`: Read-only access (list services, view logs, get status)

## Configuration
The application uses a YAML configuration file (`config.yaml`) with the following structure:

```yaml
server:
    host: localhost
    port: 40200
    readTimeout: 15s
    writeTimeout: 15s
    maxHeaderBytes: 1048576
auth:
    secretKey: your-secret-key
    tokenDuration: 24h0m0s
    issuedBy: ChronoServe
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            username: admin
            password_hash: "$argon2id$v=19$..."
            roles:
                - admin
logging:
    level: debug
    directory: logs
    maxSize: 10
    maxBackups: 5
    maxAge: 30
    compress: true
```

## Examples

### PowerShell Authentication
```powershell
$body = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

$token = $response.data.token

# Use token in subsequent requests
$headers = @{
    Authorization = "Bearer $token"
}

# Example: List services
Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

### Common Windows Services
- Windows Update (`wuauserv`)
- Event Log (`EventLog`)
- Windows Remote Management (`WinRM`)
- Background Intelligent Transfer (`BITS`)

Example viewing Windows Update service logs:
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services/logs/wuauserv?lines=50" -Headers $headers
```