package windows

// List of Windows services that are critical to system operation and should be protected
var criticalWindowsServices = []string{
	"wininit",           // Windows Start-Up Application
	"csrss",             // Client Server Runtime Process
	"services",          // Services and Controller app
	"lsass",             // Local Security Authority Process
	"winlogon",          // Windows Logon
	"smss",              // Windows Session Manager
	"svchost",           // Service Host
	"spooler",           // Print Spooler
	"explorer",          // Windows Explorer
	"fontdrvhost",       // Font Driver Host
	"dwm",               // Desktop Window Manager
	"taskmgr",           // Task Manager
	"conhost",           // Console Window Host
	"dllhost",           // COM Surrogate
	"audiodg",           // Windows Audio Device Graph Isolation
	"wuauserv",          // Windows Update
	"EventLog",          // Windows Event Log
	"TermService",       // Remote Desktop Services
	"Schedule",          // Task Scheduler
	"Dnscache",          // DNS Client
	"BITS",              // Background Intelligent Transfer Service
	"TrustedInstaller",  // Windows Modules Installer
	"PcaSvc",            // Program Compatibility Assistant Service
	"LanmanServer",      // Server
	"LanmanWorkstation", // Workstation
	"Dhcp",              // DHCP Client
	"WinDefend",         // Windows Defender Antivirus
	"wscsvc",            // Windows Security Center
	"samss",             // Security Accounts Manager
	"RpcSs",             // Remote Procedure Call (RPC)
	"nsi",               // Network Store Interface Service
	"netlogon",          // Net Logon Service
	"PlugPlay",          // Plug and Play device detection
	"VSS",               // Volume Shadow Copy
	"fdPHost",           // Function Discovery Provider Host
	"fdResPub",          // Function Discovery Resource Publication
	"DiagTrack",         // Connected User Experiences and Telemetry
	"sppsvc",            // Software Protection Platform (Windows activation)
	"wmiApSrv",          // WMI Performance Adapter
}
