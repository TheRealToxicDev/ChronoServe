# SysManix Nginx Setup Guide

This guide explains how to set up Nginx as a reverse proxy for SysManix, providing HTTPS encryption, load balancing, and enhanced security.

## Benefits of Using Nginx with SysManix

- **HTTPS Support**: Secure communication with TLS encryption
- **Load Balancing**: Distribute requests across multiple SysManix instances
- **Rate Limiting**: Protect your API from abuse and DDoS attacks
- **Authentication**: Add an additional layer of protection
- **Caching**: Improve performance for frequently requested endpoints
- **Path Rewriting**: Serve SysManix under a specific URL path

## Basic Nginx Setup

### Installing Nginx

#### On Debian/Ubuntu:
```bash
sudo apt update
sudo apt install nginx
```

#### On CentOS/RHEL:
```bash
sudo yum install epel-release
sudo yum install nginx
```

#### On Windows:
Download and install Nginx from [nginx.org](https://nginx.org/en/download.html)

### Creating a Basic Configuration

1. Create a new Nginx server configuration:

   ```bash
   # Linux
   sudo nano /etc/nginx/sites-available/sysmanix

   # Windows
   # Edit C:\nginx\conf\nginx.conf
   ```

2. Add a basic proxy configuration:

   ```nginx
   server {
       listen 80;
       server_name sysmanix.example.com;

       location / {
           proxy_pass http://localhost:40200;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

3. Enable the site configuration (Linux only):

   ```bash
   sudo ln -s /etc/nginx/sites-available/sysmanix /etc/nginx/sites-enabled/
   ```

4. Test the configuration and restart Nginx:

   ```bash
   # Linux
   sudo nginx -t
   sudo systemctl restart nginx

   # Windows
   C:\nginx\nginx.exe -t
   C:\nginx\nginx.exe -s reload
   ```

## Setting Up HTTPS with Let's Encrypt

### Installing Certbot

#### On Debian/Ubuntu:
```bash
sudo apt install certbot python3-certbot-nginx
```

#### On CentOS/RHEL:
```bash
sudo yum install certbot python3-certbot-nginx
```

### Obtaining SSL Certificates

```bash
sudo certbot --nginx -d sysmanix.example.com
```

Follow the prompts to complete the certificate issuance and configuration.

### Manual SSL Configuration

If you prefer to configure SSL manually or are using a different certificate provider:

1. Modify your Nginx configuration:

   ```nginx
   server {
       listen 443 ssl http2;
       server_name sysmanix.example.com;

       ssl_certificate /etc/nginx/ssl/sysmanix.crt;
       ssl_certificate_key /etc/nginx/ssl/sysmanix.key;
       ssl_protocols TLSv1.2 TLSv1.3;
       ssl_ciphers HIGH:!aNULL:!MD5;
       ssl_prefer_server_ciphers on;

       location / {
           proxy_pass http://localhost:40200;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto https;
       }
   }

   # Redirect HTTP to HTTPS
   server {
       listen 80;
       server_name sysmanix.example.com;
       return 301 https://$host$request_uri;
   }
   ```

2. Restart Nginx to apply changes.

## Advanced Nginx Configurations

### Load Balancing

If you're running multiple SysManix instances for high availability:

```nginx
# Define upstream servers
upstream sysmanix_backend {
    server 192.168.1.10:40200 weight=3;
    server 192.168.1.11:40200 weight=1;
    server 192.168.1.12:40200 backup;
}

server {
    listen 443 ssl http2;
    server_name sysmanix.example.com;

    # SSL configuration omitted for brevity

    location / {
        proxy_pass http://sysmanix_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

### Rate Limiting

Protect your API from abuse:

```nginx
# Define rate limiting zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=5r/s;
limit_req_status 429;

server {
    # Server configuration omitted for brevity

    location / {
        # Apply rate limiting (5 requests per second with burst of 10)
        limit_req zone=api_limit burst=10 nodelay;

        proxy_pass http://localhost:40200;
        # Other proxy settings omitted for brevity
    }

    # Allow higher rate for health check endpoint
    location /health {
        proxy_pass http://localhost:40200/health;
        # Other proxy settings omitted for brevity
    }
}
```

### Path-Based Routing

Serve SysManix under a specific URL path:

```nginx
server {
    # Server configuration omitted for brevity

    # Serve SysManix under /api path
    location /api/ {
        # Strip /api prefix when forwarding
        rewrite ^/api/(.*) /$1 break;
        proxy_pass http://localhost:40200/;
        # Other proxy settings omitted for brevity
    }

    # Serve static frontend files
    location / {
        root /var/www/sysmanix-ui;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
}
```

### Additional Security Headers

Improve security with additional headers:

```nginx
server {
    # Server configuration omitted for brevity

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload";
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options SAMEORIGIN;
    add_header X-XSS-Protection "1; mode=block";
    add_header Content-Security-Policy "default-src 'self'; script-src 'self'; img-src 'self'; style-src 'self';";

    location / {
        proxy_pass http://localhost:40200;
        # Other proxy settings omitted for brevity
    }
}
```

### Basic Authentication

Add an additional authentication layer:

1. Create an authentication file:

   ```bash
   # Install htpasswd tool if needed
   sudo apt install apache2-utils

   # Create password file
   sudo htpasswd -c /etc/nginx/.htpasswd admin
   # Follow the prompts to set a password
   ```

2. Configure Nginx to use it:

   ```nginx
   server {
       # Server configuration omitted for brevity

       location / {
           auth_basic "Restricted Access";
           auth_basic_user_file /etc/nginx/.htpasswd;

           proxy_pass http://localhost:40200;
           # Other proxy settings omitted for brevity
       }

       # Skip authentication for health checks
       location /health {
           auth_basic off;
           proxy_pass http://localhost:40200/health;
           # Other proxy settings omitted for brevity
       }
   }
   ```

### WebSocket Support

If SysManix uses WebSockets for real-time updates:

```nginx
server {
    # Server configuration omitted for brevity

    location / {
        proxy_pass http://localhost:40200;

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Other proxy settings omitted for brevity
    }
}
```

### Caching Responses

For endpoints that don't change frequently:

```nginx
# Define cache location
proxy_cache_path /var/cache/nginx/sysmanix levels=1:2 keys_zone=sysmanix_cache:10m max_size=1g inactive=60m;

server {
    # Server configuration omitted for brevity

    location /services {
        proxy_cache sysmanix_cache;
        proxy_cache_valid 200 1m;  # Cache successful responses for 1 minute
        proxy_cache_methods GET;
        proxy_cache_key $request_uri;

        proxy_pass http://localhost:40200/services;
        # Other proxy settings omitted for brevity
    }

    # Don't cache endpoints that modify data
    location ~ ^/(services/start|services/stop) {
        proxy_no_cache 1;
        proxy_cache_bypass 1;

        proxy_pass http://localhost:40200;
        # Other proxy settings omitted for brevity
    }
}
```

## Nginx Configuration for Docker Environments

If running SysManix in Docker:

```nginx
server {
    listen 443 ssl http2;
    server_name sysmanix.example.com;

    # SSL configuration omitted for brevity

    location / {
        # Docker host networking
        proxy_pass http://docker-host-ip:40200;
        # Or use Docker network name
        # proxy_pass http://sysmanix:40200;

        # Other proxy settings omitted for brevity
    }
}
```

## Troubleshooting Nginx Configuration

### Common Issues

1. **502 Bad Gateway**:
   - Verify SysManix is running: `curl http://localhost:40200/health`
   - Check Nginx error logs: `sudo tail -f /var/log/nginx/error.log`
   - Ensure SELinux/AppArmor allows Nginx to connect to SysManix

2. **SSL Certificate Issues**:
   - Verify certificate paths in Nginx configuration
   - Check certificate expiry: `openssl x509 -in /path/to/cert.crt -noout -enddate`
   - Validate certificate chain: `openssl verify -CAfile /path/to/chain.pem /path/to/cert.crt`

3. **Permission Denied**:
   - Check Nginx worker process permissions
   - Ensure SSL private keys have correct permissions (typically 600)

### Verifying Your Configuration

Test your Nginx configuration before applying changes:

```bash
# Linux
sudo nginx -t

# Windows
C:\nginx\nginx.exe -t
```

## Multi-Environment Setup

For managing multiple environments (dev, staging, production):

```nginx
# Include a specific environment config
include /etc/nginx/sysmanix/sysmanix-environment.conf;
```

Create separate environment files:

```bash
# /etc/nginx/sysmanix/sysmanix-production.conf
server {
    listen 443 ssl http2;
    server_name sysmanix.example.com;
    # Production configuration
}

# /etc/nginx/sysmanix/sysmanix-staging.conf
server {
    listen 443 ssl http2;
    server_name staging.sysmanix.example.com;
    # Staging configuration
}
```

Then create a symbolic link to the desired environment:

```bash
sudo ln -sf /etc/nginx/sysmanix/sysmanix-production.conf /etc/nginx/sysmanix/sysmanix-environment.conf
```

## Conclusion

Using Nginx as a reverse proxy for SysManix provides enhanced security, improved performance, and additional management capabilities. The configurations in this guide can be adapted to suit your specific requirements, from simple deployments to complex high-availability setups.

For more help with Nginx configuration, refer to the [official Nginx documentation](https://nginx.org/en/docs/).
