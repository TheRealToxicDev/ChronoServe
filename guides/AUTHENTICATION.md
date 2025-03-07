# SysManix Authentication Guide

This guide explains the authentication system used by SysManix.

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

### Adding a New User

To add a new user:

1. Add a user entry to the `users` section of the `config.yaml` file:
   ```yaml
   newuser:
     username: "newuser"
     password: "secure-password"  # Will be hashed after first run
     roles:
       - viewer
   ```

2. Restart SysManix
3. SysManix will automatically hash the password and update the configuration

## Authentication Flow

1. The client sends credentials to the login endpoint
2. SysManix verifies the credentials and issues a JWT token
3. The client includes this token in all subsequent requests
4. SysManix validates the token and checks the user's roles for each request

## Obtaining a Token

```bash
curl -X POST http://localhost:40200/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"your-secure-admin-password"}'
```

```powershell
# Using PowerShell
$auth = @{
    username = "admin"
    password = "your-password"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:40200/auth/login" -Method Post -Body $auth -ContentType "application/json"
$token = $response.data.token
```

A successful login response looks like:

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

## Using the Token

Include the token in the `Authorization` header for all authenticated requests:

```bash
# Using curl
curl -X GET http://localhost:40200/services \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

```powershell
# Using PowerShell
$headers = @{
    "Authorization" = "Bearer $token"
}

Invoke-RestMethod -Uri "http://localhost:40200/services" -Headers $headers
```

## Token Management

SysManix provides endpoints to manage your authentication tokens.

### List Active Tokens

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/auth/tokens
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/auth/tokens" -Headers $headers
```

Response:

```json
{
  "status": "success",
  "message": "User tokens retrieved successfully",
  "data": [
    {
      "tokenId": "01H5RZ4WH5P3Z72SF",
      "userId": "admin",
      "roles": ["admin"],
      "issuedAt": "2023-06-15T14:23:45Z",
      "expiresAt": "2023-06-16T14:23:45Z"
    }
  ]
}
```

### Revoke a Specific Token

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"tokenId":"01H5RZ4WH5P3Z72SF"}' http://localhost:40200/auth/tokens/revoke
```

```powershell
# Using PowerShell
$revokeBody = @{
    tokenId = "01H5RZ4WH5P3Z72SF"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:40200/auth/tokens/revoke" -Method Post -Headers $headers -Body $revokeBody -ContentType "application/json"
```

### Revoke All Your Tokens

```bash
# Using curl
curl -X POST http://localhost:40200/auth/tokens/revoke-all \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

```powershell
# Using PowerShell
Invoke-RestMethod -Uri "http://localhost:40200/auth/tokens/revoke-all" -Method Post -Headers $headers
```

### Refresh Your Token

```bash
curl -X POST -H "Authorization: Bearer $TOKEN" http://localhost:40200/auth/tokens/refresh
```

```powershell
# Using PowerShell
$refreshResponse = Invoke-RestMethod -Uri "http://localhost:40200/auth/tokens/refresh" -Method Post -Headers $headers

# Update token for future requests
$token = $refreshResponse.data.token
$headers = @{
    "Authorization" = "Bearer $token"
}
```

## Admin Token Management

Users with the `admin` role have access to additional token management capabilities:

### List All Tokens (Admin)

```bash
# Using curl
curl -X GET http://localhost:40200/auth/admin/tokens \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

```powershell
# Using PowerShell (with admin token)
Invoke-RestMethod -Uri "http://localhost:40200/auth/admin/tokens" -Headers $adminHeaders
```

### List Tokens for a Specific User (Admin)

```bash
# Using curl
curl -X GET "http://localhost:40200/auth/admin/tokens/user?userId=viewer" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE"
```

```powershell
# Using PowerShell (with admin token)
Invoke-RestMethod -Uri "http://localhost:40200/auth/admin/tokens/user?userId=viewer" -Headers $adminHeaders
```

### Revoke All Tokens for a User (Admin)

```bash
# Using curl
curl -X POST http://localhost:40200/auth/admin/tokens/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN_HERE" \
  -d '{"userId":"viewer"}'
```

```powershell
# Using PowerShell (with admin token)
$revokeBody = @{
    userId = "viewer"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:40200/auth/admin/tokens/revoke" -Method Post -Headers $adminHeaders -Body $revokeBody -ContentType "application/json"
```

## Role-Based Access Control

SysManix implements role-based access control with these built-in roles:

| Role    | Permissions |
|---------|------------|
| admin   | Full access to all endpoints including service start/stop and token management |
| viewer  | Read-only access to services, logs, and status |

### Permission Matrix

| Endpoint | Method | admin | viewer |
|----------|--------|-------|--------|
| `/services` | GET | ✅ | ✅ |
| `/services/status/{name}` | GET | ✅ | ✅ |
| `/services/logs/{name}` | GET | ✅ | ✅ |
| `/services/start/{name}` | POST | ✅ | ❌ |
| `/services/stop/{name}` | POST | ✅ | ❌ |
| `/auth/tokens` | GET | ✅ | ✅ |
| `/auth/tokens/revoke` | POST | ✅ | ✅ |
| `/auth/admin/tokens` | GET | ✅ | ❌ |

## Security Considerations

### Password Security

SysManix uses Argon2id for password hashing with these parameters:

- Memory: 64 MB
- Iterations: 1
- Parallelism: 4
- Key length: 32 bytes

These parameters provide strong security against various attack vectors, including brute force and side-channel attacks.

### JWT Security

To enhance JWT security:

1. Use a strong, randomly generated `secretKey` (at least 32 characters)
2. Set an appropriate `tokenDuration` based on your security requirements
3. Rotate the secret key periodically in production environments
4. Use HTTPS in production to protect token transmission

### Best Practices

1. **Secure Storage**: Store tokens securely in your client applications
2. **Token Refresh**: For long-running sessions, refresh tokens periodically
3. **Least Privilege**: Assign the minimum necessary roles to users
4. **Token Revocation**: Implement token logout by revoking tokens when sessions end
5. **Regular Audits**: Periodically review and clean up active tokens

## Troubleshooting

### Common Authentication Issues

1. **Invalid credentials**:
   - Verify username and password
   - Check if the user exists in the configuration

2. **Token expired**:
   - Use the refresh token endpoint
   - Authenticate again to get a new token

3. **Insufficient permissions**:
   - Verify the user has the required role for the operation
   - Check the JWT payload to confirm assigned roles

4. **Token validation errors**:
   - Ensure the token is being sent correctly in the Authorization header
   - Verify the token hasn't been tampered with

### Debugging Authentication

To troubleshoot authentication issues:

1. Check the server logs for authentication failures
2. Decode your JWT token at [jwt.io](https://jwt.io/) to inspect the payload
3. Verify the token expiration time in the payload

## Further Reading

- [JWT.io](https://jwt.io/introduction): Introduction to JWT tokens
- [Argon2 Specifications](https://github.com/P-H-C/phc-winner-argon2/blob/master/argon2-specs.pdf): Details on the Argon2 hashing algorithm
- [API Reference](./API_REFERENCE.md): Complete API documentation
- [Configuration Guide](./CONFIGURATION.md): Detailed authentication configuration
