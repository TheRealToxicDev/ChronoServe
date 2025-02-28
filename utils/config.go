package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the root configuration structure
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Linux   LinuxConfig   `yaml:"linux"`
	Windows WindowsConfig `yaml:"windows"`
	Logging LogConfig     `yaml:"logging"`
}

type ServerConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	ReadTimeout    string `yaml:"readTimeout"`
	WriteTimeout   string `yaml:"writeTimeout"`
	MaxHeaderBytes int    `yaml:"maxHeaderBytes"`
}

type AuthConfig struct {
	SecretKey     string                 `yaml:"secretKey"`
	TokenDuration time.Duration          `yaml:"tokenDuration"`
	IssuedBy      string                 `yaml:"issuedBy"`
	AllowedRoles  []string               `yaml:"allowedRoles"`
	Users         map[string]Credentials `yaml:"users"`
}

type Credentials struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Roles    []string `yaml:"roles"`
}

type LinuxConfig struct {
	ServiceCommand string             `yaml:"serviceCommand"`
	LogDirectory   string             `yaml:"logDirectory"`
	Services       map[string]Service `yaml:"services"`
}

type WindowsConfig struct {
	ServiceCommand string             `yaml:"serviceCommand"`
	LogDirectory   string             `yaml:"logDirectory"`
	Services       map[string]Service `yaml:"services"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Directory  string `yaml:"directory"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

type Service struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Enabled      bool     `yaml:"enabled"`
	AllowedRoles []string `yaml:"allowedRoles"`
}

// Default configuration values
var defaultConfig = Config{
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
		IssuedBy:      "ChronoServe",
		AllowedRoles:  []string{"admin", "viewer"},
		Users: map[string]Credentials{
			"admin": {
				Username: "admin",
				Password: "change-me",
				Roles:    []string{"admin"},
			},
			"viewer": {
				Username: "viewer",
				Password: "change-me",
				Roles:    []string{"viewer"},
			},
		},
	},
	Linux: LinuxConfig{
		ServiceCommand: "systemctl",
		LogDirectory:   "/var/log/chronoserve",
	},
	Windows: WindowsConfig{
		ServiceCommand: "sc",
		LogDirectory:   "C:\\ProgramData\\ChronoServe\\logs",
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

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Port)
	}

	if c.Auth.SecretKey == "" || c.Auth.SecretKey == defaultConfig.Auth.SecretKey {
		return fmt.Errorf("security risk: default secret key must be changed")
	}

	if len(c.Auth.AllowedRoles) == 0 {
		return fmt.Errorf("at least one role must be defined")
	}

	if c.Logging.MaxSize < 1 {
		return fmt.Errorf("invalid log max size: %d", c.Logging.MaxSize)
	}

	return nil
}

func InitConfig(filePath string) error {
	// Check if config file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create default config based on OS
		config = defaultConfig

		// Save the default config
		if err := SaveConfig(filePath); err != nil {
			return fmt.Errorf("error creating initial config: %w", err)
		}

		fmt.Printf("\n=== Security Notice ===\n")
		fmt.Printf("A new configuration file has been created at: %s\n", filePath)
		fmt.Println("Please update the following security-sensitive values before running again:")
		fmt.Println("1. auth.secretKey")
		fmt.Println("2. auth.users.admin.password")
		fmt.Println("3. auth.users.viewer.password")
		fmt.Println("\nExiting for security reasons. Please update the configuration and restart.")
		os.Exit(1)
	}

	// Load existing config
	if err := LoadConfig(filePath); err != nil {
		return err
	}

	// Check for default values in existing config
	if hasDefaultCredentials() {
		fmt.Printf("\n=== Security Risk Detected ===\n")
		fmt.Println("Default credentials found in configuration.")
		fmt.Println("Please update the security-sensitive values before running.")
		fmt.Println("Exiting for security reasons.")
		os.Exit(1)
	}

	// Validate configuration
	return config.Validate()
}

// Add this new helper function
func hasDefaultCredentials() bool {
	// Check if any security-sensitive values are still set to defaults
	if config.Auth.SecretKey == defaultConfig.Auth.SecretKey {
		return true
	}

	// Check admin credentials
	if admin, exists := config.Auth.Users["admin"]; exists {
		if admin.Password == defaultConfig.Auth.Users["admin"].Password {
			return true
		}
	}

	// Check viewer credentials
	if viewer, exists := config.Auth.Users["viewer"]; exists {
		if viewer.Password == defaultConfig.Auth.Users["viewer"].Password {
			return true
		}
	}

	return false
}

var (
	config     Config
	configLock sync.RWMutex
)

// GetConfig returns a copy of the current configuration
func GetConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filePath string) error {
	configLock.Lock()
	defer configLock.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Merge with defaults for any missing values
	mergeWithDefaults(&cfg)

	// Update global config
	config = cfg
	return nil
}

// SaveConfig saves the current configuration to a YAML file
func SaveConfig(filePath string) error {
	configLock.RLock()
	defer configLock.RUnlock()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
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
}

// UpdateConfig updates the configuration and optionally saves it to disk
func UpdateConfig(newConfig Config, save bool) error {
	configLock.Lock()
	defer configLock.Unlock()

	if err := newConfig.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	config = newConfig

	if save {
		return SaveConfig("config.yaml")
	}
	return nil
}

// GetServiceConfig returns the service configuration for the current OS
func GetServiceConfig() interface{} {
	configLock.RLock()
	defer configLock.RUnlock()

	if runtime.GOOS == "linux" {
		return config.Linux
	}
	return config.Windows
}
