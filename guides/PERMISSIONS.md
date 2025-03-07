# SysManix Permissions Guide

This guide explains SysManix's role-based access control system and how to configure permissions for users.

## Permission Model Overview

SysManix uses a role-based access control (RBAC) system to determine what operations each user can perform. The key components are:

- **Users**: Individual accounts that authenticate with the system
- **Roles**: Named collections of permissions (e.g., "admin", "viewer")
- **Permissions**: Authorization to perform specific operations

## Built-in Roles

SysManix comes with two pre-defined roles:

### Admin Role

The `admin` role has full access to all features:

- View services and statuses
- Start and stop services
- View service logs
- Manage tokens (view, revoke, refresh)
- View and manage other users' tokens
- Access protected API endpoints
- Manage system configuration

### Viewer Role

The `viewer` role has read-only access:

- View services and statuses
- View service logs
- Manage their own tokens (view, revoke, refresh)
- No access to start/stop services
- No access to other users' tokens
- No access to protected API endpoints

## Role-Based API Access

This table shows which endpoints each role can access:

| Endpoint | Method | admin | viewer |
|----------|--------|-------|--------|
| `/services` | GET | ✅ | ✅ |
| `/services/status/{name}` | GET | ✅ | ✅ |
| `/services/logs/{name}` | GET | ✅ | ✅ |
| `/services/start/{name}` | POST | ✅ | ❌ |
| `/services/stop/{name}` | POST | ✅ | ❌ |
| `/auth/login` | POST | ✅ | ✅ |
| `/auth/tokens` | GET | ✅ | ✅ |
| `/auth/tokens/revoke` | POST | ✅ | ✅ |
| `/auth/tokens/revoke-all` | POST | ✅ | ✅ |
| `/auth/tokens/refresh` | POST | ✅ | ✅ |
| `/auth/admin/tokens` | GET | ✅ | ❌ |
| `/auth/admin/tokens/user` | GET | ✅ | ❌ |
| `/auth/admin/tokens/revoke` | POST | ✅ | ❌ |
| `/health` | GET | ✅ | ✅ |

## User Configuration

Users and their roles are defined in the configuration file (`config.yaml`):

```yaml
auth:
  # Other auth settings...
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
    custom_user:
      username: "custom_user"
      password: "plain_password"  # Will be hashed on next restart
      roles:
        - viewer
```

## Creating Custom Roles

SysManix allows you to define custom roles for more granular permission control:

1. Define the custom roles in the configuration:
   ```yaml
   auth:
     allowedRoles:
       - admin
       - viewer
       - operator
       - db_admin
       - web_admin
   ```

2. Assign custom roles to users:
   ```yaml
   auth:
     users:
       db_user:
         username: "db_user"
         password: "secure_password"
         roles:
           - db_admin
       web_user:
         username: "web_user"
         password: "secure_password"
         roles:
           - web_admin
   ```

3. Define custom permissions for services:
   ```yaml
   windows:
     services:
       customPermissions:
         sqlserver:
           allowedRoles:
             - admin
             - db_admin
         iis:
           allowedRoles:
             - admin
             - web_admin

   linux:
     services:
       customPermissions:
         postgresql:
           allowedRoles:
             - admin
             - db_admin
         nginx:
           allowedRoles:
             - admin
             - web_admin
   ```

## Example Role Combinations

### Operator Role

The `operator` role might have permission to start and stop services, but not manage users or tokens:

```yaml
auth:
  allowedRoles:
    - admin
    - viewer
    - operator

  users:
    operator:
      username: "operator"
      password: "secure_password"
      roles:
        - operator
```

With the corresponding middleware configuration in code:

```go
// Example middleware registration for an operator role
registerRouteWithMiddleware(mux, "services", services.ListServices, true, []string{"admin", "viewer", "operator"})
registerRouteWithMiddleware(mux, "services/start/", services.StartService, true, []string{"admin", "operator"})
registerRouteWithMiddleware(mux, "services/stop/", services.StopService, true, []string{"admin", "operator"})
registerRouteWithMiddleware(mux, "services/logs/", services.ViewServiceLogs, true, []string{"admin", "viewer", "operator"})
registerRouteWithMiddleware(mux, "services/status/", services.GetServiceStatus, true, []string{"admin", "viewer", "operator"})
```

### Multi-Role User

Users can have multiple roles for combined permissions:

```yaml
auth:
  users:
    power_user:
      username: "power_user"
      password: "secure_password"
      roles:
        - db_admin
        - web_admin
```

## Assigning Permissions Programmatically

Roles and permissions are enforced by middleware in the code:

```go
// Example permission middleware
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := GetClaimsFromContext(r.Context())
            if claims == nil {
                utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
                return
            }

            // Check if user has any of the required roles
            hasRole := false
            for _, requiredRole := range roles {
                for _, userRole := range claims.Roles {
                    if userRole == requiredRole {
                        hasRole = true
                        break
                    }
                }
                if hasRole {
                    break
                }
            }

            if !hasRole {
                utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## Best Practices

### Principle of Least Privilege

Always assign the minimum necessary permissions:

- Give users only the roles they need
- Create specific roles for specific tasks
- Avoid using the admin role for routine tasks

### Service-Specific Permissions

For sensitive services, create dedicated roles:

```yaml
windows:
  services:
    customPermissions:
      BackupService:
        allowedRoles:
          - admin
          - backup_operator
      FinanceApp:
        allowedRoles:
          - admin
          - finance_team
```

### Role Review Schedule

Regularly review user roles:

1. Audit users and their assigned roles quarterly
2. Remove unused or unnecessary roles
3. Verify that service-specific permissions are still appropriate
4. Document role assignments and changes

## Advanced Permission Configurations

### Implementing Time-Based Restrictions

SysManix can be extended to support time-based access controls:

```yaml
auth:
  timeRestrictions:
    maintenance_window:
      allowedRoles:
        - admin
        - operator
      schedule:
        - weekdays: [Saturday, Sunday]
          hours: [0-23]
        - weekdays: [Monday-Friday]
          hours: [18-23, 0-8]
```

### IP-Based Restrictions

Limit access to specific IP ranges:

```yaml
auth:
  ipRestrictions:
    officeNetwork:
      ranges:
        - 192.168.1.0/24
        - 10.0.0.0/8
      allowedRoles:
        - admin
        - viewer
        - operator

    adminNetwork:
      ranges:
        - 192.168.10.0/24
      allowedRoles:
        - admin
```

### Service Group Permissions

Group services for easier permission management:

```yaml
auth:
  serviceGroups:
    webServices:
      - nginx
      - apache
      - haproxy

    databaseServices:
      - postgresql
      - mysql
      - mongodb

  roles:
    web_admin:
      serviceGroups:
        - webServices

    db_admin:
      serviceGroups:
        - databaseServices
```

## Troubleshooting

### Permission Denied Errors

If you encounter "Forbidden" or "Permission Denied" errors:

1. Verify the user has the correct role assigned
2. Check if the role has permission for the specific endpoint
3. For service-specific operations, verify any custom permissions
4. Check the logs for authentication and authorization details

### Debugging Role Assignments

To debug role-related issues:

1. Use the `/health` endpoint to verify the API is functioning
2. Check the JWT token payload (at [jwt.io](https://jwt.io/)) to verify roles
3. Review the configuration file for correct role assignments
4. Monitor authentication logs for login attempts and role information

### Common Issues

1. **Missing roles in token**: Ensure roles are properly assigned and included in JWT
2. **Configuration not reloaded**: Restart SysManix after changing roles
3. **Case sensitivity**: Role names are case-sensitive
4. **Protected service restrictions**: Some service operations are restricted regardless of role

## Upgrading Permissions

When upgrading SysManix, be aware of permission changes:

1. Back up the existing configuration
2. Review any new roles or permission changes in the release notes
3. Update your custom roles as needed
4. Test permissions after upgrade

## Next Steps

- Learn about [Service Management](./SERVICE_MANAGEMENT.md) to understand service operations
- Explore the [API Reference](./API_REFERENCE.md) for detailed endpoint information
- Review [Security Guide](./SECURITY.md) for security best practices
