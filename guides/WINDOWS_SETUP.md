# SysManix Windows Setup Guide

This guide provides detailed instructions for setting up SysManix on Windows systems.

## Installation

1. **Download the latest release**:

   ```powershell
   Invoke-WebRequest -Uri "https://github.com/toxic-development/SysManix/releases/latest/download/SysManix_windows_amd64.exe" -OutFile "SysManix.exe"
   ```

2. **Run the executable to generate initial config**:

   ```powershell
   .\SysManix.exe
   ```

3. **Edit the configuration file**:

   ```yaml
   auth:
     secretKey: "your-secure-random-string-here"
     users:
       admin:
         password: "your-secure-admin-password"
       viewer:
         password: "your-secure-viewer-password"
   ```

4. **Start SysManix**:

   ```powershell
   .\SysManix.exe
   ```

## Running as a Windows Service

1. **Using NSSM**:

   ```powershell
   nssm.exe install SysManix "C:\path\to\SysManix.exe"
   nssm.exe set SysManix DisplayName "SysManix Service Manager"
   nssm.exe set SysManix Start SERVICE_AUTO_START
   Start-Service SysManix
   ```

## Firewall Configuration

```powershell
New-NetFirewallRule -DisplayName "SysManix API" -Direction Inbound -Protocol TCP -LocalPort 40200 -Action Allow
```

## Viewing Logs

```powershell
Get-EventLog -LogName Application | Where-Object {$_.Source -eq "SysManix"}
```
