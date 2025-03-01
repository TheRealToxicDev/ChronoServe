# ChronoServe Authentication System

## Overview

ChronoServe uses a JWT-based authentication system with Argon2id password hashing for secure user management.

## Password Security

### Argon2id Configuration

```go
type PasswordConfig struct {
    time    uint32 // CPU/memory cost parameter (1)
    memory  uint32 // Memory size in KiB (64 * 1024)
    threads uint8  // Number of threads (4)
    keyLen  uint32 // Hash length in bytes (32)
}
```

### Password Hashing Flow

1. Initial password storage in `config.yaml`:
```yaml
users:
    admin:
        username: "admin"
        password: "your-password"  # Plain text (temporary)
        roles: ["admin"]
```

2. On first run, the application:
   - Detects plain text passwords
   - Generates a secure salt
   - Hashes using Argon2id
   - Updates config:

```yaml
users:
    admin:
        username: "admin"
        password_hash: "$argon2id$v=19$m=65536,t=1,p=4$[salt]$[hash]"
        roles: ["admin"]
```

### Hash Format
The Argon2id hash format follows: `$argon2id$v=[version]$m=[memory],t=[time],p=[parallelism]$[salt-base64]$[hash-base64]`

## JWT Authentication

### Token Structure

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "uid": "user-id",
    "roles": ["admin"],
    "exp": 1735689600,
    "iat": 1735603200,
    "iss": "ChronoServe"
  }
}
```

### Token Generation Process

1. User submits credentials:
```json
{
    "username": "admin",
    "password": "your-password"
}
```

2. Server:
   - Retrieves stored user data
   - Verifies password using constant-time comparison
   - Generates JWT with user roles
   - Returns token response:

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

## Role-Based Access Control (RBAC)

### Available Roles

| Role    | Permissions |
|---------|------------|
| admin   | Full access to all endpoints |
| viewer  | Read-only access (list, status, logs) |

### Role Validation

1. Request arrives with JWT token
2. Token is validated for:
   - Signature validity
   - Expiration time
   - Required roles
3. Access granted or denied based on roles

### Role Configuration

```yaml
auth:
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            roles: ["admin"]
        viewer:
            roles: ["viewer"]
```

## Security Implementation

### Password Storage
- Uses Argon2id (winner of Password Hashing Competition)
- Secure random salt generation
- Automatic removal of plain text passwords
- Constant-time comparison for verification

### Token Security
- HMAC-SHA256 signing
- Configurable expiration time
- Secure secret key storage
- Token validation on every request

### Code Examples

#### Password Verification
```go
func VerifyPassword(password, hash string) (bool, error) {
    // Extract parameters from hash
    params := strings.Split(hash, "$")
    if len(params) != 5 {
        return false, fmt.Errorf("invalid hash format")
    }

    // Constant-time comparison of hashes
    newHash := argon2.IDKey(
        []byte(password),
        salt,
        time,
        memory,
        threads,
        keyLen,
    )
    return subtle.ConstantTimeCompare(existingHash, newHash) == 1, nil
}
```

#### Token Validation
```go
func ValidateToken(token string) (*Claims, error) {
    // Parse and validate JWT token
    claims := &Claims{}
    parsedToken, err := jwt.ParseWithClaims(
        token,
        claims,
        func(t *jwt.Token) (interface{}, error) {
            return []byte(config.GetConfig().Auth.SecretKey), nil
        },
    )

    if err != nil || !parsedToken.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}
```

## Security Best Practices

1. Password Requirements
   - Minimum length: 12 characters
   - Mix of uppercase, lowercase, numbers, symbols
   - No common dictionary words

2. Token Management
   - Short expiration times (24h default)
   - Secure transmission over HTTPS
   - Token invalidation on logout

3. Configuration Security
   - Change default secret key
   - Use strong passwords
   - Regular key rotation

4. Error Handling
   - Generic error messages
   - No information disclosure
   - Proper logging of authentication attempts

## Troubleshooting

### Common Issues

1. Invalid Credentials
   - Verify username matches config key
   - Check password was entered correctly
   - Ensure password was hashed properly

2. Token Validation Failures
   - Check token expiration
   - Verify secret key in config
   - Confirm proper token format

3. Role Access Denied
   - Verify user has required roles
   - Check role configuration
   - Ensure token contains roles

## API Examples

### Login Request
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
```

### Protected Request
```powershell
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```