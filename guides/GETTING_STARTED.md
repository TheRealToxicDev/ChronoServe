# Getting Started with SysManix

This guide will help you get up and running with SysManix quickly.

## Installation

### Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_windows_amd64.exe" -OutFile "SysManix.exe"
.\SysManix.exe
```

### Linux

```bash
wget https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_linux_amd64 -O sysmanix
chmod +x sysmanix
sudo ./sysmanix
```

## Configuration

Edit the `config.yaml` file created in the same directory as the executable:

```yaml
auth:
  secretKey: "your-secure-random-string-here"
  users:
    admin:
      password: "your-secure-admin-password"
    viewer:
      password: "your-secure-viewer-password"
```

## Starting SysManix

### Windows

```powershell
.\SysManix.exe
```

### Linux

```bash
sudo ./sysmanix
```

## Accessing the API

### Authentication

```bash
curl -X POST http://localhost:40200/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"your-secure-admin-password"}'
```

### Basic Service Management

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:40200/services
```
