# SysManix Authentication System

## Table of Contents

1. [Overview](#overview)
2. [Authentication Flow](#authentication-flow)
3. [Password Security](#password-security)
4. [JWT Token System](#jwt-token-system)
5. [Role-Based Access Control](#role-based-access-control)
6. [Security Implementation Details](#security-implementation-details)
7. [Command-Line Authentication](#command-line-authentication)
8. [Troubleshooting](#troubleshooting)
9. [Best Practices](#best-practices)
10. [API Examples](#api-examples)

## Overview

SysManix implements a modern, secure authentication system combining:

- **JWT (JSON Web Tokens)** for stateless authentication
- **Argon2id** password hashing (winner of the Password Hashing Competition)
- **Role-based access control** for fine-grained permissions
- **Constant-time comparison** for secure credential checking

This document explains how authentication works, how to use it, and best practices for security.

## Authentication Flow

The SysManix authentication flow consists of the following steps:

1. **User Registration** (Config-based)
   - Users are defined in the `config.yaml` file
   - On first run, plain-text passwords are hashed automatically

2. **Login Process**
   - Client sends username and password to `/auth/login`
   - Server validates credentials against stored hashes
   - If valid, server generates a JWT token containing the user's roles
   - Token is returned to the client for subsequent requests

3. **Request Authorization**
   - Client includes the JWT token in the `Authorization` header
   - Server validates the token signature and expiration
   - Server checks the user's roles against required permissions
   - If authorized, the request is processed

4. **Token Expiration**
   - Tokens expire after a configurable period (default: 24 hours)
   - When expired, users must re-authenticate

### Authentication Sequence Diagram

```
┌──────┐                                  ┌──────────┐
│Client│                                  │SysManix  │
└──┬───┘                                  └────┬─────┘
   │       POST /auth/login                    │
   │      {username, password}                 │
   │─────────────────────────────────────────>│
   │                                          │
   │                                          │ Validate credentials
   │                                          │ Generate JWT
   │                                          │
   │      200 OK                              │
   │      {token: "eyJhbGciOiJI..."}          │
   │<─────────────────────────────────────────│
   │                                          │
   │       GET /services                       │
   │       Authorization: Bearer eyJhbGciOiJI...│
   │─────────────────────────────────────────>│
   │                                          │
   │                                          │ Validate JWT
   │                                          │ Check permissions
   │                                          │
   │       200 OK                             │
   │       {services: [...]}                  │
   │<─────────────────────────────────────────│
```

## Password Security

### Argon2id Configuration

SysManix uses the Argon2id algorithm, which is designed to be resistant to both side-channel and brute force attacks. The default configuration balances security and performance:

```go
type PasswordConfig struct {
    time    uint32 // CPU cost parameter (1)
    memory  uint32 // Memory cost in KiB (64 * 1024 = 64MB)
    threads uint8  // Number of threads (4)
    keyLen  uint32 // Hash output length in bytes (32)
}
```

These parameters can be adjusted in more security-sensitive environments. Higher values for time and memory increase security but slow down the hashing process.

### Password Hashing Flow

1. **Initial Configuration**:
   When setting up SysManix, passwords are initially stored in plain text:

```yaml
users:
    admin:
        username: "admin"
        password: "your-secure-password"  # Plain text (temporary)
        roles: ["admin"]
```

2. **First Run Hash Generation**:
   On first run, SysManix automatically:
   - Detects plain text passwords
   - Generates a secure random salt for each password
   - Hashes each password using Argon2id with the salt
   - Updates the config file, replacing plain text with secure hashes:

```yaml
users:
    admin:
        username: "admin"
        password_hash: "$argon2id$v=19$m=65536,t=1,p=4$[salt-base64]$[hash-base64]"
        roles: ["admin"]
```

3. **Future Authentication**:
   For subsequent authentication attempts, SysManix:
   - Extracts the salt and parameters from the hash
   - Hashes the provided password with the same salt and parameters
   - Uses constant-time comparison to verify the hashes match

### Hash Format

The Argon2id hash format follows this pattern:
`$argon2id$v=[version]$m=[memory],t=[time],p=[parallelism]$[salt-base64]$[hash-base64]`

- `argon2id`: The algorithm identifier
- `v=19`: The Argon2 version number
- `m=65536`: Memory cost (64MB)
- `t=1`: Time cost (1 iteration)
- `p=4`: Parallelism (4 threads)
- `[salt-base64]`: Base64-encoded random salt
- `[hash-base64]`: Base64-encoded password hash

This format ensures all necessary information for validation is stored within the hash itself.

## JWT Token System

### Token Structure

JWT tokens used by SysManix have three parts: header, payload, and signature:

#### Header
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

#### Payload
```json
{
  "uid": "admin",
  "roles": ["admin"],
  "exp": 1735689600,
  "iat": 1735603200,
  "iss": "SysManix"
}
```

Key fields in the payload:
- `uid`: User identifier
- `roles`: Array of role names assigned to the user
- `exp`: Expiration timestamp
- `iat`: Issued-at timestamp
- `iss`: Issuer (always "SysManix")

#### Signature
The signature is created using the HMAC-SHA256 algorithm:
```
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secretKey
)
```

The complete token is the concatenation of these three parts with periods:
```
base64UrlEncode(header) + "." + base64UrlEncode(payload) + "." + signature
```

### Token Generation Process

1. **Login Request**:
   The client submits credentials:
```json
{
    "username": "admin",
    "password": "your-secure-password"
}
```

2. **Server Processing**:
   - The server retrieves the user's data from config
   - Verifies the password using Argon2id
   - If valid, creates a new JWT with the user's roles
   - Signs the JWT with the secret key

3. **Login Response**:
   The server returns the JWT token:
```json
{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1...",
        "roles": ["admin"]
    }
}
```

### Token Validation Process

When a protected endpoint receives a request:

1. The `Authorization` header is extracted
2. The JWT token is separated from the "Bearer " prefix
3. The token signature is verified using the secret key
4. The token expiration is checked
5. User roles from the token are checked against required roles
6. If all checks pass, the request is processed

## Role-Based Access Control

SysManix implements role-based access control (RBAC) to manage permission levels for different users.

### Available Roles

| Role    | Permissions |
|---------|-------------|
| admin   | Full access to all endpoints, including service control operations |
| viewer  | Read-only access (list services, view status, view logs) |

### Role Validation

1. When a request arrives with a JWT token:
   - The token is validated for authenticity and expiration
   - User roles are extracted from the token claims

2. The required roles for the endpoint are checked:
   - Most endpoints require either `admin` or `viewer` role
   - Service control operations (start/stop) require the `admin` role

3. Access is granted or denied based on whether the user has the required roles:
   - If the user has appropriate roles, the request is processed
   - Otherwise, a 403 Forbidden error is returned

### Role Configuration

Roles are defined in the `config.yaml` file:

```yaml
auth:
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            username: "admin"
            password_hash: "$argon2id$v=19$..."
            roles: ["admin"]
        viewer:
            username: "viewer"
            password_hash: "$argon2id$v=19$..."
            roles: ["viewer"]
        poweruser:
            username: "poweruser"
            password_hash: "$argon2id$v=19$..."
            roles: ["admin", "viewer"]
```

Users can have multiple roles assigned to them. The middleware checks if a user has any of the required roles for an endpoint.

### Protected Endpoint Matrix

| Endpoint | Method | Required Roles | Description |
|----------|--------|----------------|-------------|
| /health | GET | None (Public) | Health check endpoint |
| /auth/login | POST | None (Public) | Authentication endpoint |
| /services | GET | admin, viewer | List all services |
| /services/status/{name} | GET | admin, viewer | Get service status |
| /services/logs/{name} | GET | admin, viewer | View service logs |
| /services/start/{name} | POST | admin | Start a service |
| /services/stop/{name} | POST | admin | Stop a service |

## Security Implementation Details

### Password Verification

SysManix uses constant-time comparison for password verification to prevent timing attacks:

```go
func VerifyPassword(password, hash string) (bool, error) {
    // Extract parameters from hash
    params := strings.Split(hash, "$")
    if len(params) != 5 {
        return false, fmt.Errorf("invalid hash format")
    }

    // Parse parameters, salt, and stored hash
    // ...

    // Generate comparison hash with same parameters and salt
    newHash := argon2.IDKey(
        []byte(password),
        salt,
        time,
        memory,
        threads,
        keyLen,
    )

    // Constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(existingHash, newHash) == 1, nil
}
```

### Token Validation

JWT tokens are validated with multiple security checks:

```go
func ValidateToken(token string) (*Claims, error) {
    claims := &Claims{}
    parsedToken, err := jwt.ParseWithClaims(
        token,
        claims,
        func(t *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return []byte(config.GetConfig().Auth.SecretKey), nil
        },
    )

    if err != nil {
        return nil, err
    }

    if !parsedToken.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}
```

### Token Extraction

The `Authorization` header is safely extracted:

```go
func extractToken(r *http.Request) (string, error) {
    authHeader := r.Header.Get("Authorization")
    if (authHeader == "") {
        return "", fmt.Errorf("no authorization header")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
        return "", fmt.Errorf("invalid authorization header format")
    }

    // Use constant-time comparison for token type
    if subtle.ConstantTimeCompare([]byte(strings.ToLower(parts[0])), []byte("bearer")) != 1 {
        return "", fmt.Errorf("invalid authorization type")
    }

    return parts[1], nil
}
```

## Command-Line Authentication

### PowerShell Authentication

```powershell
# Login and retrieve token
$body = @{
    username = "admin"
    password = "your-secure-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body

# Store the token
$token = $response.data.token

# Create headers for future requests
$headers = @{
    Authorization = "Bearer $token"
}

# Example: List all services
$services = Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

### Linux/macOS Authentication (bash/curl)

```bash
# Login and retrieve token
TOKEN=$(curl -s -X POST http://localhost:40200/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-secure-password"}' \
  | jq -r '.data.token')

# Example: List all services
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services
```

### Python Example

```python
import requests
import json

# Login and retrieve token
response = requests.post(
    "http://localhost:40200/auth/login",
    headers={"Content-Type": "application/json"},
    data=json.dumps({"username": "admin", "password": "your-secure-password"})
)

token = response.json()["data"]["token"]

# Example: List all services
services = requests.get(
    "http://localhost:40200/services",
    headers={"Authorization": f"Bearer {token}"}
)

print(json.dumps(services.json(), indent=2))
```

## Troubleshooting

### Common Authentication Issues

#### Invalid Credentials Error

**Symptom:**
```json
{
    "success": false,
    "error": "Invalid credentials"
}
```

**Causes and Solutions:**
1. **Username mismatch**
   - Verify that the username matches exactly what's in config.yaml
   - Check for case sensitivity (usernames are case-sensitive)

2. **Password mismatch**
   - If you've forgotten the password, add a new plain-text password to config.yaml and remove the password_hash
   - Restart the app to generate a new hash

3. **Config file issues**
   - Ensure the username in login request matches the key in the users map:
   ```yaml
   users:
       admin:            # This must match the login username
           username: "admin"
           password_hash: "..."
   ```

#### Token Validation Failed

**Symptom:**
```json
{
    "success": false,
    "error": "Invalid or expired token"
}
```

**Causes and Solutions:**
1. **Token expired**
   - Re-authenticate to get a new token
   - Check if the server's system clock is accurate
   - Consider increasing token duration in config.yaml

2. **Invalid token format**
   - Ensure the token is passed correctly in the Authorization header
   - Format should be: `Authorization: Bearer your-token-here`
   - Check for extra spaces or encoding issues

3. **Secret key mismatch**
   - If the server's secret key was changed, all existing tokens become invalid
   - Re-authenticate to get a new token

#### Role Access Denied

**Symptom:**
```json
{
    "success": false,
    "error": "Forbidden",
    "code": 403
}
```

**Causes and Solutions:**
1. **Insufficient permissions**
   - Check that your user has the required role for the endpoint
   - Service control endpoints require admin role

2. **Role configuration**
   - Verify roles are correctly assigned in config.yaml
   - Check that the role exists in the allowedRoles list

## Best Practices

### Password Security

1. **Strong Passwords**
   - Minimum length: 12 characters
   - Mix of uppercase, lowercase, numbers, and symbols
   - Avoid dictionary words and common patterns
   - Unique for each system

2. **Secret Key Protection**
   - Change the default secret key to a strong random value
   - Minimum 32 characters, random alphanumeric + symbols
   - Keep the secret key secure and consider using environment variables

3. **Password Storage**
   - Never store plain-text passwords in production
   - Let SysManix handle password hashing automatically

### Token Management

1. **Token Expiration**
   - Use a reasonable expiration time (24h is default)
   - Shorter for higher security environments (e.g., 1h)
   - Re-authenticate when tokens expire

2. **Token Transmission**
   - Always use HTTPS in production
   - Keep tokens secure in your client applications
   - Clear tokens when logging out or on session end

3. **Error Handling**
   - Don't expose sensitive data in error messages
   - Implement rate limiting for login attempts

## API Examples

### Login Request and Response

**Request:**
```http
POST /auth/login HTTP/1.1
Host: localhost:40200
Content-Type: application/json

{
    "username": "admin",
    "password": "your-secure-password"
}
```

**Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
    "success": true,
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "roles": ["admin"]
    }
}
```

### Protected API Request

```http
GET /services HTTP/1.1
Host: localhost:40200
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Unauthorized Response

```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
    "success": false,
    "error": "Invalid or expired token",
    "code": 401
}
```

### Forbidden Response

```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
    "success": false,
    "error": "Insufficient permissions to access this resource",
    "code": 403
}
```

---

For more information on how authentication integrates with service management, see [SERVICE_MANAGEMENT.md](SERVICE_MANAGEMENT.md) and [PERMISSIONS.md](PERMISSIONS.md).
