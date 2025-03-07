# SysManix Permissions and Role-Based Access Control

## Overview

SysManix implements a role-based access control (RBAC) system to ensure users can only perform actions they're authorized for. This guide explains how permissions work and how to configure them.

## Role Types

SysManix has two built-in roles:

| Role | Description | Permissions |
|------|-------------|-------------|
| `admin` | Full administrative access | All operations including service management |
| `viewer` | Read-only access | Can view services, status, and logs but cannot modify |

## Permission Matrix

| Endpoint | Method | admin | viewer |
|----------|--------|-------|--------|
| `/health` | GET | ✅ | ✅ |
| `/auth/login` | POST | ✅ | ✅ |
| `/services` | GET | ✅ | ✅ |
| `/services/status/{name}` | GET | ✅ | ✅ |
| `/services/logs/{name}` | GET | ✅ | ✅ |
| `/services/start/{name}` | POST | ✅ | ❌ |
| `/services/stop/{name}` | POST | ✅ | ❌ |

## Configuration

Roles are configured in the `config.yaml` file:

```yaml
auth:
    allowedRoles:
        - admin
        - viewer
    users:
        admin:
            username: admin
            password_hash: "$argon2id$v=19$..."
            roles:
                - admin
        viewer:
            username: viewer
            password_hash: "$argon2id$v=19$..."
            roles:
                - viewer
        poweruser:
            username: poweruser
            password_hash: "$argon2id$v=19$..."
            roles:
                - admin
                - viewer
```

## JWT Implementation

When a user logs in, their roles are encoded in the JWT token:

```json
{
  "uid": "admin",
  "roles": ["admin"],
  "exp": 1735689600,
  "iat": 1735603200,
  "iss": "SysManix"
}
```

The middleware validates these roles on each protected request:

```go
func ValidateRoles(requiredRoles []string, userRoles []string) bool {
    for _, required := range requiredRoles {
        hasRole := false
        for _, role := range userRoles {
            if role == required {
                hasRole = true
                break
            }
        }
        if !hasRole {
            return false
        }
    }
    return true
}
```

## Protected Services

Regardless of a user's role, certain system-critical services are protected from modification:

- `admin` users cannot start or stop protected services
- Protected services are filtered from the service list
- Attempts to modify protected services return a 403 Forbidden response

## Error Responses

### Insufficient Permissions

If a user lacks the required role for an endpoint:

```json
{
    "success": false,
    "error": "Insufficient permissions to access this resource",
    "code": 403
}
```

### Protected Service

If a user attempts to modify a protected service:

```json
{
    "success": false,
    "error": "operation not allowed on protected system service: systemd",
    "code": 403
}
```

### Invalid Token

If the JWT token is invalid or expired:

```json
{
    "success": false,
    "error": "Invalid or expired token",
    "code": 401
}
```

## Best Practices

1. **Principle of Least Privilege**: Assign the minimum required roles to users
2. **Token Expiration**: Use reasonable JWT token expiration times (default: 24h)
3. **Role Separation**: Create separate users for different responsibilities
4. **Audit Logging**: Monitor failed permission attempts in auth logs
5. **Regular Review**: Periodically review and update user roles

## Adding Custom Roles

While SysManix comes with predefined `admin` and `viewer` roles, you can implement custom roles by:

1. Adding new role names to `allowedRoles` in config
2. Assigning the roles to users
3. Adding role checks in your middleware or handlers

Example custom role configuration:

```yaml
auth:
    allowedRoles:
        - admin
        - viewer
        - service_manager
        - log_viewer
    users:
        custom_user:
            username: custom_user
            password_hash: "$argon2id$v=19$..."
            roles:
                - service_manager
                - log_viewer
```
