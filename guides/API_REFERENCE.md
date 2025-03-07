# SysManix API Reference

This document provides a comprehensive reference for the SysManix API endpoints, including request and response formats, authentication requirements, and examples.

## API Overview

The SysManix API is a RESTful interface for managing system services across different operating systems. It provides a consistent set of endpoints regardless of whether the underlying system is Windows or Linux.

Base URL: `http://localhost:40200` (default)

## Authentication

All protected endpoints require authentication using a JWT token provided in the `Authorization` header.

**Authentication Header Format:**
```
Authorization: Bearer your-jwt-token
```

### Obtaining a Token

```
POST /auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "your-password"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "roles": ["admin"]
  }
}
```

## Endpoints

### Health Endpoint

```
GET /health
```

Returns system health information. Does not require authentication.

**Response:**
```json
{
  "status": "healthy",
  "uptime": "3h 24m 15s",
  "version": "0.1.0",
  "goVersion": "go1.23.1",
  "memory": {
    "alloc": 8453120,
    "totalAlloc": 28974656,
    "sys": 20103168,
    "numGC": 14,
    "heapInUse": 7.87
  },
  "startTime": "2023-05-12T14:22:31Z"
}
```

### Service Management

#### List All Services

```
GET /services
```

Lists all available services on the system.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`

**Response:**
```json
{
  "status": "success",
  "message": "Services retrieved successfully",
  "data": [
    {
      "name": "nginx",
      "displayName": "Nginx Web Server",
      "status": "Running",
      "isActive": true,
      "updatedAt": "2023-05-12T17:46:23Z"
    },
    {
      "name": "postgresql",
      "displayName": "PostgreSQL Database Server",
      "status": "Stopped",
      "isActive": false,
      "updatedAt": "2023-05-12T17:46:23Z"
    }
  ]
}
```

#### Get Service Status

```
GET /services/status/{service-name}
```

Returns the current status of a specific service.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`
**Path Parameters:**
- `service-name`: Name of the service

**Response:**
```json
{
  "status": "success",
  "message": "Service status retrieved successfully",
  "data": {
    "name": "nginx",
    "status": "Running",
    "isActive": true,
    "updatedAt": "2023-05-12T17:46:23Z"
  }
}
```

#### Start Service

```
POST /services/start/{service-name}
```

Starts a specific service.

**Authentication Required:** Yes
**Required Roles:** `admin`
**Path Parameters:**
- `service-name`: Name of the service

**Response:**
```json
{
  "status": "success",
  "message": "Service nginx started successfully"
}
```

#### Stop Service

```
POST /services/stop/{service-name}
```

Stops a specific service.

**Authentication Required:** Yes
**Required Roles:** `admin`
**Path Parameters:**
- `service-name`: Name of the service

**Response:**
```json
{
  "status": "success",
  "message": "Service nginx stopped successfully"
}
```

#### View Service Logs

```
GET /services/logs/{service-name}
```

Returns log entries for a specific service.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`
**Path Parameters:**
- `service-name`: Name of the service

**Query Parameters:**
- `lines`: Number of log lines to return (default: 100)

**Response:**
```json
{
  "status": "success",
  "message": "Service logs retrieved successfully",
  "data": [
    {
      "time": "2023-05-12 17:30:45",
      "level": "INFO",
      "message": "Starting service"
    },
    {
      "time": "2023-05-12 17:30:46",
      "level": "INFO",
      "message": "Service started successfully"
    }
  ]
}
```

### Token Management

#### List User Tokens

```
GET /auth/tokens
```

Lists all active tokens for the current user.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`

**Response:**
```json
{
  "status": "success",
  "message": "User tokens retrieved successfully",
  "data": [
    {
      "tokenId": "01H6P4QZ3T9VWXYZ",
      "userId": "admin",
      "roles": ["admin"],
      "issuedAt": "2023-05-12T15:23:45Z",
      "expiresAt": "2023-05-13T15:23:45Z"
    }
  ]
}
```

#### Revoke Token

```
POST /auth/tokens/revoke
```

Revokes a specific token.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer` (can only revoke own tokens unless `admin`)

**Request Body:**
```json
{
  "tokenId": "01H6P4QZ3T9VWXYZ"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Token revoked successfully"
}
```

#### Revoke All User Tokens

```
POST /auth/tokens/revoke-all
```

Revokes all tokens for the current user.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`

**Response:**
```json
{
  "status": "success",
  "message": "All tokens revoked successfully",
  "data": {
    "count": 3
  }
}
```

#### Refresh Token

```
POST /auth/tokens/refresh
```

Generates a new token and invalidates the current one.

**Authentication Required:** Yes
**Required Roles:** `admin` or `viewer`

**Response:**
```json
{
  "status": "success",
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Admin Token Management

#### List All Tokens (Admin)

```
GET /auth/admin/tokens
```

Lists all valid tokens in the system.

**Authentication Required:** Yes
**Required Roles:** `admin`

**Response:**
```json
{
  "status": "success",
  "message": "All tokens retrieved successfully",
  "data": [
    {
      "tokenId": "01H6P4QZ3T9VWXYZ",
      "userId": "admin",
      "roles": ["admin"],
      "issuedAt": "2023-05-12T15:23:45Z",
      "expiresAt": "2023-05-13T15:23:45Z"
    },
    {
      "tokenId": "01H6P4R2M8S7PQRS",
      "userId": "viewer",
      "roles": ["viewer"],
      "issuedAt": "2023-05-12T16:34:56Z",
      "expiresAt": "2023-05-13T16:34:56Z"
    }
  ]
}
```

#### List User Tokens (Admin)

```
GET /auth/admin/tokens/user?userId=username
```

Lists all valid tokens for a specific user.

**Authentication Required:** Yes
**Required Roles:** `admin`
**Query Parameters:**
- `userId`: Username to list tokens for

**Response:**
```json
{
  "status": "success",
  "message": "User tokens retrieved successfully",
  "data": [
    {
      "tokenId": "01H6P4R2M8S7PQRS",
      "userId": "viewer",
      "roles": ["viewer"],
      "issuedAt": "2023-05-12T16:34:56Z",
      "expiresAt": "2023-05-13T16:34:56Z"
    }
  ]
}
```

#### Revoke User Tokens (Admin)

```
POST /auth/admin/tokens/revoke
```

Revokes all tokens for a specific user.

**Authentication Required:** Yes
**Required Roles:** `admin`

**Request Body:**
```json
{
  "userId": "viewer"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "All user tokens revoked successfully",
  "data": {
    "count": 2,
    "userId": "viewer"
  }
}
```

## Error Responses

All API errors follow a consistent format:

```json
{
  "status": "error",
  "message": "Description of the error",
  "code": 404
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input parameters |
| 401 | Unauthorized - Missing or invalid authentication token |
| 403 | Forbidden - Insufficient permissions for the action |
| 404 | Not Found - The requested resource does not exist |
| 408 | Request Timeout - Operation took too long to complete |
| 500 | Internal Server Error - Unexpected server error |

## Using the API with Different Tools

### cURL Examples

**Authentication:**
```bash
curl -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}'
```

**List Services:**
```bash
curl -X GET http://localhost:40200/services \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Start Service:**
```bash
curl -X POST http://localhost:40200/services/start/nginx \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### PowerShell Examples

**Authentication:**
```powershell
$auth = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" -Method Post -Body $auth -ContentType "application/json"
$token = $response.data.token
$headers = @{ "Authorization" = "Bearer $token" }
```

**List Services:**
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

**Start Service:**
```powershell
Invoke-RestMethod -Uri "http://localhost:40200/services/start/wuauserv" -Method Post -Headers $headers
```

### JavaScript/Fetch Examples

**Authentication:**
```javascript
const response = await fetch('http://localhost:40200/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'your-password' })
});

const data = await response.json();
const token = data.data.token;
```

**List Services:**
```javascript
const response = await fetch('http://localhost:40200/services', {
  headers: { 'Authorization': `Bearer ${token}` }
});

const services = await response.json();
```

## API Versioning

This documentation covers API version 1.0.0. Future versions may introduce changes to the endpoint structure or response format. The API version can be verified through the `/health` endpoint.

## Rate Limiting

The API may implement rate limiting to protect against abuse. When rate limits are exceeded, the API will return a 429 Too Many Requests status code.

## Further Information

For more details on specific aspects of the API:
- [Configuration Guide](./CONFIGURATION.md) - Details on configuring API behavior
- [Authentication Guide](./AUTHENTICATION.md) - In-depth information on authentication
- [Service Management Guide](./SERVICE_MANAGEMENT.md) - Detailed usage examples
