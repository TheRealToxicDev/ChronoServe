package linux

// List of Linux services that are critical to system operation and should be protected
var criticalLinuxServices = []string{
	"systemd",           // Core system daemon
	"systemd-journald",  // Journal logging service
	"systemd-logind",    // Login service
	"systemd-udevd",     // udev management daemon
	"sshd",              // SSH daemon
	"dbus",              // D-Bus system message bus
	"NetworkManager",    // Network management
	"polkit",            // Authorization manager
	"user@",             // User manager
	"selinux",           // SELinux policy manager
	"accounts-daemon",   // Accounts service
	"agetty",            // Console manager
	"apparmor",          // AppArmor security service
	"cron",              // Scheduled tasks
	"rsyslog",           // System logging
	"systemd-resolved",  // DNS resolver
	"systemd-timesyncd", // Time synchronization
	"systemd-networkd",  // Network configuration
	"init",              // System V init process
	"syslog",            // System logging service
	"wpa_supplicant",    // Wireless authentication
	"firewalld",         // Firewall service
	"iptables",          // Firewall management (legacy)
	"dnsmasq",           // DNS caching and DHCP
	"cups",              // Printing service
	"fail2ban",          // Intrusion prevention
	"nftables",          // Modern firewall replacement for iptables
	"pam",               // Pluggable Authentication Module (PAM) daemon
	"lvm2",              // Logical Volume Manager
	"mdadm",             // RAID array management
	"udisks2",           // Disk management service
	"ntpd",              // Network Time Protocol daemon
	"chronyd",           // Alternative to ntpd for time sync
}
