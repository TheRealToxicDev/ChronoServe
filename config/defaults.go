package config

import (
	"runtime"
	"strings"
	"time"
)

// Default configuration values
var defaultConfig = Config{
	API: APIConfig{
		EnableSwagger: true,
		SwaggerPath:   "/swagger/",
		Version:       "1.0.0",
		Title:         "SysManix API",
		Description:   "Cross-platform service management API",
	},
	Server: ServerConfig{
		Host:           "localhost",
		Port:           8080,
		ReadTimeout:    "15s",
		WriteTimeout:   "15s",
		MaxHeaderBytes: 1 << 20, // 1MB
	},
	Auth: AuthConfig{
		SecretKey:     "change-me",
		TokenDuration: 24 * time.Hour,
		IssuedBy:      "SysManix",
		AllowedRoles:  []string{"admin", "viewer"},
		Users: map[string]Credentials{
			"admin": {
				Username: "admin",
				Password: "change-me", // Will be hashed on first run
				Roles:    []string{"admin"},
			},
			"viewer": {
				Username: "viewer",
				Password: "change-me", // Will be hashed on first run
				Roles:    []string{"viewer"},
			},
		},
	},
	Linux: LinuxConfig{
		ServiceCommand: "systemctl",
		LogDirectory:   "/var/log/SysManix",
	},
	Windows: WindowsConfig{
		ServiceCommand: "sc",
		LogDirectory:   "C:\\ProgramData\\SysManix\\logs",
	},
	Logging: LogConfig{
		Level:      "info",
		Directory:  "logs",
		MaxSize:    10, // 10MB
		MaxBackups: 5,
		MaxAge:     30, // 30 days
		Compress:   true,
	},
}

// mergeWithDefaults fills in any missing values with defaults
func mergeWithDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = defaultConfig.Server.Port
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = defaultConfig.Server.Host
	}
	if cfg.Server.ReadTimeout == "" {
		cfg.Server.ReadTimeout = defaultConfig.Server.ReadTimeout
	}
	if cfg.Server.WriteTimeout == "" {
		cfg.Server.WriteTimeout = defaultConfig.Server.WriteTimeout
	}
	if cfg.Server.MaxHeaderBytes == 0 {
		cfg.Server.MaxHeaderBytes = defaultConfig.Server.MaxHeaderBytes
	}

	// Auth defaults
	if cfg.Auth.TokenDuration == 0 {
		cfg.Auth.TokenDuration = defaultConfig.Auth.TokenDuration
	}
	if cfg.Auth.IssuedBy == "" {
		cfg.Auth.IssuedBy = defaultConfig.Auth.IssuedBy
	}
	if len(cfg.Auth.AllowedRoles) == 0 {
		cfg.Auth.AllowedRoles = defaultConfig.Auth.AllowedRoles
	}

	// OS-specific defaults
	if runtime.GOOS == "linux" {
		if cfg.Linux.ServiceCommand == "" {
			cfg.Linux.ServiceCommand = defaultConfig.Linux.ServiceCommand
		}
		if cfg.Linux.LogDirectory == "" {
			cfg.Linux.LogDirectory = defaultConfig.Linux.LogDirectory
		}
	} else if runtime.GOOS == "windows" {
		if cfg.Windows.ServiceCommand == "" {
			cfg.Windows.ServiceCommand = defaultConfig.Windows.ServiceCommand
		}
		if cfg.Windows.LogDirectory == "" {
			cfg.Windows.LogDirectory = defaultConfig.Windows.LogDirectory
		}
	}

	// Logging defaults
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = defaultConfig.Logging.Level
	}
	if cfg.Logging.Directory == "" {
		cfg.Logging.Directory = defaultConfig.Logging.Directory
	}
	if cfg.Logging.MaxSize == 0 {
		cfg.Logging.MaxSize = defaultConfig.Logging.MaxSize
	}
	if cfg.Logging.MaxBackups == 0 {
		cfg.Logging.MaxBackups = defaultConfig.Logging.MaxBackups
	}
	if cfg.Logging.MaxAge == 0 {
		cfg.Logging.MaxAge = defaultConfig.Logging.MaxAge
	}

	// API defaults
	if cfg.API.Version == "" {
		cfg.API.Version = defaultConfig.API.Version
	}
	if cfg.API.Title == "" {
		cfg.API.Title = defaultConfig.API.Title
	}
	if cfg.API.Description == "" {
		cfg.API.Description = defaultConfig.API.Description
	}
	if cfg.API.SwaggerPath == "" {
		cfg.API.SwaggerPath = defaultConfig.API.SwaggerPath
	}
	// We only set EnableSwagger if it's not explicitly set to false
	// This ensures it's enabled by default
}

// hasDefaultCredentials checks if the configuration has the default credentials
func hasDefaultCredentials() bool {
	if config.Auth.SecretKey == defaultConfig.Auth.SecretKey {
		return true
	}

	for _, user := range config.Auth.Users {
		if user.PasswordHash == "" || strings.HasPrefix(user.PasswordHash, "$argon2id$v=19$m=65536,t=1,p=4$") {
			return true
		}
	}

	return false
}
