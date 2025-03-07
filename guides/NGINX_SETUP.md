# SysManix Nginx Integration Guide

This guide explains how to set up Nginx as a reverse proxy for SysManix, enabling HTTPS, load balancing, and additional security features.

## Prerequisites

- SysManix installed and running
- Nginx installed on your server
- Basic understanding of Nginx configuration
- (Optional) Domain name for your SysManix instance
- (Optional) SSL certificate for HTTPS

## Why Use Nginx with SysManix?

Placing Nginx in front of SysManix provides several benefits:

- **HTTPS Support**: Encrypt API traffic with TLS/SSL
- **Load Balancing**: Distribute requests across multiple SysManix instances
- **Additional Security**: Rate limiting, IP filtering, and WAF capabilities
- **Caching**: Improve performance by caching responses
- **Path Rewriting**: Serve SysManix under a specific URL path
- **Authentication**: Add an additional layer of authentication

## Basic Nginx Configuration

### Simple HTTP Proxy

Create a new configuration file for SysManix:

```bash
sudo nano /etc/nginx/sites-available/sysmanix
```

Add the following basic configuration:

```nginx
server {
    listen 80;
    server_name sysmanix.example.com;  # Replace with your domain or IP

    access_log /var/log/nginx/sysmanix-access.log;
    error_log /var/log/nginx/sysmanix-error.log;

    location / {
        proxy_pass http://localhost:40200;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site and reload Nginx:

```bash
sudo ln -s /etc/nginx/sites-available/sysmanix /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### HTTPS Configuration with Let's Encrypt

1. Install Certbot for Let's Encrypt certificates:

```bash
sudo apt update
sudo apt install certbot python3-certbot-nginx
```

2. Obtain and configure SSL certificate:

```bash
sudo certbot --nginx -d sysmanix.example.com
```

3. Certbot will automatically modify your Nginx configuration for HTTPS. The final configuration should look similar to:

```nginx
server {
    listen 80;
    server_name sysmanix.example.com;
    return 301 https://$server_name$request_uri;  # Redirect HTTP to HTTPS
}

server {
    listen 443 ssl;
    server_name sysmanix.example.com;

    ssl_certificate /etc/letsencrypt/live/sysmanix.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/sysmanix.example.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # SSL configuration
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_stapling on;
    ssl_stapling_verify on;

    # Security headers
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";

    access_log /var/log/nginx/sysmanix-access.log;
    error_log /var/log/nginx/sysmanix-error.log;

    location / {
        proxy_pass http://localhost:40200;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Advanced Configurations

### Load Balancing Multiple SysManix Instances

If you're running multiple SysManix instances for high availability, configure load balancing:

```nginx
upstream sysmanix_backend {
    server 192.168.1.10:40200;
    server 192.168.1.11:40200;
    server 192.168.1.12:40200;
}

server {
    listen 443 ssl;
    server_name sysmanix.example.com;

    # SSL configuration here...

    location / {
        proxy_pass http://sysmanix_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Set load balancing method (optional)
        proxy_next_upstream error timeout http_500;

        # Health checks
        health_check interval=10 fails=3 passes=2;
    }
}
```

### Rate Limiting

Add rate limiting to protect your SysManix API from abuse:

```nginx
# Define a limit zone based on client IP
limit_req_zone $binary_remote_addr zone=sysmanix_limit:10m rate=10r/s;

server {
    # ...existing configuration...

    location / {
        # Apply rate limiting
        limit_req zone=sysmanix_limit burst=20 nodelay;

        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...
    }

    # Different rate limit for login endpoint
    location /auth/login {
        limit_req zone=sysmanix_limit burst=5 nodelay;
        proxy_pass http://localhost:40200/auth/login;
        # ...remaining proxy configuration...
    }
}
```

### IP Whitelisting

Restrict access to the SysManix API by IP address:

```nginx
server {
    # ...existing configuration...

    # Allow access only from specific IPs
    allow 10.0.0.0/8;     # Internal network
    allow 192.168.1.0/24; # Office network
    deny all;             # Deny all other IPs

    location / {
        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...
    }
}
```

### Path Rewriting

Serve SysManix under a specific URL path:

```nginx
server {
    # ...existing configuration...

    # Serve SysManix under /api/sysmanix
    location /api/sysmanix/ {
        rewrite ^/api/sysmanix/(.*) /$1 break;
        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...
    }
}
```

### WebSocket Support

If SysManix uses WebSockets for real-time updates:

```nginx
server {
    # ...existing configuration...

    location / {
        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Authentication with Nginx

### Basic Authentication

Add an additional layer of authentication:

1. Create a password file:

```bash
sudo apt install apache2-utils
sudo htpasswd -c /etc/nginx/sysmanix_htpasswd admin
```

2. Configure Nginx to use basic auth:

```nginx
server {
    # ...existing configuration...

    location / {
        auth_basic "Restricted Access";
        auth_basic_user_file /etc/nginx/sysmanix_htpasswd;

        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...
    }

    # Allow login endpoint without basic auth
    location /auth/login {
        proxy_pass http://localhost:40200/auth/login;
        # ...remaining proxy configuration...
    }
}
```

## Optimizing Performance

### Response Caching

Enable caching for appropriate endpoints:

```nginx
# Define a cache zone
proxy_cache_path /var/cache/nginx/sysmanix levels=1:2 keys_zone=sysmanix_cache:10m max_size=1g inactive=60m;

server {
    # ...existing configuration...

    # Cache health endpoint responses
    location /health {
        proxy_pass http://localhost:40200/health;
        proxy_cache sysmanix_cache;
        proxy_cache_valid 200 1m;  # Cache successful responses for 1 minute
        proxy_cache_lock on;
        add_header X-Cache-Status $upstream_cache_status;
        # ...remaining proxy configuration...
    }

    # Don't cache other endpoints like /services
    location / {
        proxy_pass http://localhost:40200;
        proxy_no_cache 1;
        # ...remaining proxy configuration...
    }
}
```

### Compression

Enable compression to reduce bandwidth:

```nginx
server {
    # ...existing configuration...

    # Enable gzip compression
    gzip on;
    gzip_comp_level 5;
    gzip_min_length 256;
    gzip_proxied any;
    gzip_vary on;
    gzip_types
        application/javascript
        application/json
        application/xml
        text/css
        text/plain;

    location / {
        proxy_pass http://localhost:40200;
        # ...remaining proxy configuration...
    }
}
```

## Monitoring and Logging

### Enhanced Logging Format

Create a more detailed log format:

```nginx
log_format sysmanix_detailed '$remote_addr - $remote_user [$time_local] '
                            '"$request" $status $body_bytes_sent '
                            '"$http_referer" "$http_user_agent" '
                            '$request_time $upstream_response_time $pipe';

server {
    # ...existing configuration...
    access_log /var/log/nginx/sysmanix-access.log sysmanix_detailed;
}
```

### Status Page for Monitoring

Enable Nginx status page for monitoring:

```nginx
server {
    # ...existing configuration...

    # Nginx status - restricted to local access
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
}
```

## Troubleshooting

### Checking Nginx Configuration

Always validate your configuration before applying changes:

```bash
sudo nginx -t
```

### Common Issues

1. **502 Bad Gateway**: SysManix is not running or not accessible
   - Check if SysManix is running: `systemctl status sysmanix`
   - Verify the SysManix port configuration

2. **403 Forbidden**: Permission issues with Nginx
   - Check Nginx worker process permissions
   - Verify IP restriction rules

3. **SSL Certificate Issues**:
   - Renew certificates: `sudo certbot renew`
   - Check certificate validity: `sudo certbot certificates`

4. **Headers Not Forwarded**:
   - Ensure proxy_set_header directives are correctly configured

## Security Best Practices

1. **Disable Server Tokens**: Hide Nginx version information
   ```nginx
   server_tokens off;
   ```

2. **Configure Security Headers**: Add appropriate security headers
   ```nginx
   add_header X-Content-Type-Options nosniff;
   add_header X-Frame-Options DENY;
   add_header X-XSS-Protection "1; mode=block";
   ```

3. **Implement Request Filtering**: Block suspicious requests
   ```nginx
   # Block SQL injection attempts
   if ($args ~* "([;']|--)|insert|select|union|update|delete|drop|truncate") {
       return 403;
   }
   ```

4. **Regular Updates**: Keep Nginx and SysManix updated to patch security vulnerabilities

## Further Reading

- [Official Nginx Documentation](https://nginx.org/en/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [SysManix Security Guide](./SECURITY.md)
- [Load Balancing Best Practices](./LOAD_BALANCING.md)
