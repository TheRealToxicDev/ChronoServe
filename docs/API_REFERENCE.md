# ChronoServe API Reference

## Authentication

All protected endpoints require a JWT token in the Authorization header.

```http
Authorization: Bearer your-jwt-token
```

### Login

Authenticate and receive a JWT token.

```http
POST /auth/login

Request Body:
{
    "username": "string",
    "password": "string"
}

Response (200 OK):
{
    "status": "success",
    "data": {
        "token": "eyJhbGciOiJ...",
        "roles": ["admin"]
    }
}

Response (401 Unauthorized):
{
    "status": "error",
    "message": "Invalid credentials",
    "code": 401
}
```

## Service Management

### List Services

Returns a list of all available services.

```http
GET /services

Response (200 OK):
{
    "status": "success",
    "data": {
        "services": [
            {
                "name": "string",
                "description": "string",
                "status": "running|stopped|unknown",
                "enabled": boolean
            }
        ]
    }
}
```

### Get Service Status

Get the current status of a specific service.

```http
GET /services/status/{name}

Response (200 OK):
{
    "status": "success",
    "data": {
        "name": "string",
        "status": "running|stopped|unknown",
        "enabled": boolean,
        "description": "string"
    }
}
```

### Start Service

Start a specific service (admin only).

```http
POST /services/start/{name}

Response (200 OK):
{
    "status": "success",
    "message": "Service started successfully"
}

Response (403 Forbidden):
{
    "status": "error",
    "message": "Insufficient permissions",
    "code": 403
}
```

### Stop Service

Stop a specific service (admin only).

```http
POST /services/stop/{name}

Response (200 OK):
{
    "status": "success",
    "message": "Service stopped successfully"
}
```

### View Service Logs

Retrieve logs for a specific service.

```http
GET /services/logs/{name}

Query Parameters:
- lines (optional): Number of lines to return (default: 100)
- since (optional): Return logs since timestamp (ISO 8601)

Response (200 OK):
{
    "status": "success",
    "data": {
        "logs": [
            {
                "timestamp": "2025-02-28T15:04:05Z",
                "level": "string",
                "message": "string"
            }
        ]
    }
}
```

## Health Check

Check the API server's health status.

```http
GET /health

Response (200 OK):
{
    "status": "success",
    "data": {
        "status": "healthy",
        "version": "1.0.0",
        "uptime": "24h0m0s"
    }
}
```

## Error Responses

Common error response format:

```http
{
    "status": "error",
    "message": "Error description",
    "code": number
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200  | Success |
| 400  | Bad Request |
| 401  | Unauthorized |
| 403  | Forbidden |
| 404  | Not Found |
| 500  | Internal Server Error |

## Rate Limiting

- 100 requests per minute for authenticated users
- 10 requests per minute for unauthenticated endpoints
- Rate limit headers included in responses:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1735689600
```