# SysManix Windows Setup Guide

This guide provides detailed instructions for setting up SysManix as a Windows Service for reliable operation on Windows systems.

## Prerequisites

- Windows 10/11 or Windows Server 2019/2022
- Administrator privileges
- PowerShell 5.1 or newer
- SysManix installed (see [Installation Guide](./INSTALLATION.md))

## Understanding Windows Services

Running SysManix as a Windows Service provides several benefits:

- **Automatic startup** on system boot
- **Background operation** without user login
- **Crash recovery** with automatic restart options
- **Service dependencies** management
- **Integrated logging** with Windows Event Log

## Setting Up SysManix as a Windows Service

### Using the Built-in Installer

If you used the SysManix Windows installer, it should have already set up the Windows Service. You can verify this by:

```powershell
Get-Service -Name SysManix
```

If the service exists, you can skip to the [Configuring the Windows Service](#configuring-the-windows-service) section.

### Manual Service Setup

If you need to manually set up the Windows Service:

1. Open PowerShell as Administrator

2. Navigate to your SysManix installation directory:

```powershell
cd "C:\Program Files\SysManix"  # Adjust path if needed
```

3. Create the Windows Service using the New-Service cmdlet:

```powershell
New-Service -Name "SysManix" `
            -DisplayName "SysManix Service Management" `
            -Description "Provides API access to system service management" `
            -BinaryPathName '"C:\Program Files\SysManix\SysManix.exe" --service' `
            -StartupType Automatic `
            -ErrorAction Stop
```

4. Alternatively, use the Windows Service Control (sc) command:

```powershell
sc.exe create SysManix binPath= "\"C:\Program Files\SysManix\SysManix.exe\" --service" start= auto DisplayName= "SysManix Service Management"
```

## Configuring the Windows Service

### Service Properties Configuration

1. Open the Windows Services Manager:
   - Press `Win+R`, type `services.msc`, and press Enter

2. Find the SysManix service in the list and double-click it

3. Configure the service properties:
   - **Startup Type**: Automatic
   - **Recovery**: On first, second, and subsequent failures: Restart the Service
   - **Dependencies**: Consider adding dependencies on services that SysManix needs (e.g., HTTP service)

### Setting Service Credentials

By default, the SysManix service runs under the Local System account. For enhanced security, you may want to:

1. Create a dedicated service account with limited permissions
2. Configure the service to run under this account:
   - In the service properties, go to the "Log On" tab
   - Select "This account" and provide the account credentials

### Configuring SysManix for Windows Service

Ensure the SysManix configuration is properly set up for running as a Windows service:

1. Open the configuration file (typically at `C:\Program Files\SysManix\config.yaml`)

2. Adjust the following settings:

```yaml
server:
  host: "0.0.0.0"  # Listen on all interfaces
  port: 40200

logging:
  directory: "C:\\ProgramData\\SysManix\\logs"

windows:
  logDirectory: "C:\\ProgramData\\SysManix\\logs"
```

3. Create the log directory if it doesn't exist:

```powershell
New-Item -Path "C:\ProgramData\SysManix\logs" -ItemType Directory -Force
```

## Service Control Operations

### Starting and Stopping the Service

```powershell
# Start the SysManix service
Start-Service -Name SysManix

# Stop the SysManix service
Stop-Service -Name SysManix

# Restart the SysManix service
Restart-Service -Name SysManix
```

### Checking Service Status

```powershell
# Check service status
Get-Service -Name SysManix

# Get detailed service information
Get-WmiObject -Class Win32_Service -Filter "Name='SysManix'"
```

### Setting Startup Type

```powershell
# Set to start automatically
Set-Service -Name SysManix -StartupType Automatic

# Set to start manually
Set-Service -Name SysManix -StartupType Manual

# Set to start automatically with delayed start
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\SysManix" -Name "DelayedAutostart" -Value 1 -Type DWORD
```

## Windows Firewall Configuration

Allow incoming connections to the SysManix API port:

```powershell
# Add a firewall rule for SysManix
New-NetFirewallRule -DisplayName "SysManix API" `
                   -Direction Inbound `
                   -Protocol TCP `
                   -LocalPort 40200 `
                   -Action Allow `
                   -Profile Domain,Private
```

For more selective access, restrict by IP address:

```powershell
# Allow access only from specific IP ranges
New-NetFirewallRule -DisplayName "SysManix API - Internal Access" `
                   -Direction Inbound `
                   -Protocol TCP `
                   -LocalPort 40200 `
                   -RemoteAddress 10.0.0.0/8,192.168.1.0/24 `
                   -Action Allow `
                   -Profile Domain,Private
```

## Event Log Integration

SysManix logs events to the Windows Event Log. You can view these logs using:

1. Event Viewer (GUI):
   - Press `Win+R`, type `eventvwr.msc`, and press Enter
   - Navigate to "Windows Logs" > "Application"
   - Filter for events with source "SysManix"

2. PowerShell:
```powershell
Get-WinEvent -LogName Application -MaxEvents 100 | Where-Object { $_.ProviderName -eq "SysManix" }
```

## Advanced Service Configuration

### Service Dependencies

If SysManix depends on other services, you can configure dependencies:

```powershell
# Add dependencies on HTTP service
sc.exe config SysManix depend= HTTP
```

### Service Recovery Options

Configure advanced service recovery options:

```powershell
# Set service recovery options
sc.exe failure SysManix reset= 86400 actions= restart/30000/restart/60000/restart/120000
```

### Auto-Recovery Script

Create a monitoring script that ensures SysManix is always running:

```powershell
# Create a script at C:\Scripts\Check-SysManix.ps1
$scriptContent = @'
$service = Get-Service -Name SysManix
if ($service.Status -ne 'Running') {
    Start-Service -Name SysManix
    Write-EventLog -LogName Application -Source "SysManix Monitor" -EventId 1001 -EntryType Information -Message "SysManix service was not running and has been restarted."
}
'@

Set-Content -Path "C:\Scripts\Check-SysManix.ps1" -Value $scriptContent

# Create a scheduled task that runs every 5 minutes
$action = New-ScheduledTaskAction -Execute 'powershell.exe' -Argument '-ExecutionPolicy Bypass -File "C:\Scripts\Check-SysManix.ps1"'
$trigger = New-ScheduledTaskTrigger -Once -At (Get-Date) -RepetitionInterval (New-TimeSpan -Minutes 5)
$principal = New-ScheduledTaskPrincipal -UserID "NT AUTHORITY\SYSTEM" -LogonType ServiceAccount -RunLevel Highest
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable
Register-ScheduledTask -TaskName "SysManix Monitor" -Action $action -Trigger $trigger -Principal $principal -Settings $settings
```

## Troubleshooting

### Service Won't Start

If the SysManix service fails to start:

1. Check the Windows Application Event Log for error details
2. Verify the path to the executable is correct in the service definition
3. Ensure the configuration file is valid and accessible
4. Check that the specified port is available and not in use
5. Verify that the service account has necessary permissions

### Access Denied Errors

If SysManix encounters "Access Denied" errors when managing other services:

1. Ensure the service is running as an account with appropriate permissions (typically Local System)
2. For a custom account, add it to the Windows "Administrators" group or give specific permissions:
   ```powershell
   # Grant service control permissions to a custom account
   sc.exe sdset SysManix D:(A;;CCLCSWRPWPDTLOCRRC;;;SY)(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;BA)(A;;CCLCSWLOCRRC;;;IU)(A;;CCLCSWLOCRRC;;;SU)
   ```

### Port Conflicts

If another application is already using port 40200:

1. Change the port in the SysManix configuration file:
   ```yaml
   server:
     port: 40201  # Change to an available port
   ```
2. Restart the SysManix service
3. Update any firewall rules to reflect the new port

## Security Considerations

1. **Least Privilege Principle**: Use a dedicated service account with minimal required permissions
2. **Network Isolation**: Restrict access to the SysManix API port using Windows Firewall
3. **Configuration Security**: Secure the configuration file with proper NTFS permissions
4. **Regular Updates**: Keep SysManix updated with the latest security patches
5. **Audit Logging**: Enable detailed audit logging for service operations

## Further Reading

- [Windows Service Documentation](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.management/new-service)
- [SysManix Configuration Guide](./CONFIGURATION.md)
- [SysManix Security Guide](./SECURITY.md)
- [IIS Integration Guide](./IIS_SETUP.md)
