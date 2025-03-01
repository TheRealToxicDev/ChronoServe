package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	config     Config
	configLock sync.RWMutex
	configPath string
)

// InitConfig initializes the configuration from a file or creates a new one
func InitConfig(filePath string) error {
	configPath = filePath

	// Check if config file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create new config with defaults
		config = defaultConfig

		if err := SaveConfig(filePath); err != nil {
			return fmt.Errorf("error creating initial config: %w", err)
		}

		fmt.Printf("\n=== Configuration Notice ===\n")
		fmt.Printf("A new configuration file has been created at: %s\n", filePath)
		fmt.Println("Please update the following values in the config file:")
		fmt.Println("1. auth.secretKey")
		fmt.Println("2. auth.users.admin.password")
		fmt.Println("3. auth.users.viewer.password")
		fmt.Println("\nAfter updating, restart the application.")
		os.Exit(0)
	}

	// Load existing config
	if err := LoadConfig(filePath); err != nil {
		return err
	}

	// Process any plain passwords and convert to hashes
	if err := processPlainPasswords(); err != nil {
		return err
	}

	return config.Validate()
}

func initializeDefaultPasswords() error {
	if err := config.SetUserPassword("admin", "change-me-admin"); err != nil {
		return fmt.Errorf("failed to set admin password: %w", err)
	}
	if err := config.SetUserPassword("viewer", "change-me-viewer"); err != nil {
		return fmt.Errorf("failed to set viewer password: %w", err)
	}
	return nil
}

// Validate validates the configuration
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

	// Validate user credentials
	for username, creds := range c.Auth.Users {
		if creds.PasswordHash == "" {
			return fmt.Errorf("user %s has no password hash set", username)
		}
		if !strings.HasPrefix(creds.PasswordHash, "$argon2id$") {
			return fmt.Errorf("user %s has invalid password hash format", username)
		}
		if len(creds.Roles) == 0 {
			return fmt.Errorf("user %s has no assigned roles", username)
		}
	}

	if c.Logging.MaxSize < 1 {
		return fmt.Errorf("invalid log max size: %d", c.Logging.MaxSize)
	}

	return nil
}

// SetUserPassword updates a user's password with a secure hash
func (c *Config) SetUserPassword(username, password string) error {
	configLock.Lock()
	defer configLock.Unlock()

	user, exists := c.Auth.Users[username]
	if !exists {
		return fmt.Errorf("user not found: %s", username)
	}

	if err := user.SetPassword(password); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	c.Auth.Users[username] = user
	return SaveConfig(configPath)
}

// processPlainPasswords hashes any plain text passwords in the configuration
func processPlainPasswords() error {
	configLock.Lock()
	defer configLock.Unlock()

	passwordsChanged := false

	// Check each user's password
	for username, user := range config.Auth.Users {
		// If there's a plain password and no hash
		if user.Password != "" && user.PasswordHash == "" {
			// Hash the password
			if err := user.SetPassword(user.Password); err != nil {
				return fmt.Errorf("failed to hash password for user %s: %w", username, err)
			}
			// Clear the plain password
			user.Password = ""
			config.Auth.Users[username] = user
			passwordsChanged = true
		}
	}

	// If any passwords were changed, save the config
	if passwordsChanged {
		if err := SaveConfig(configPath); err != nil {
			return fmt.Errorf("failed to save hashed passwords: %w", err)
		}
		fmt.Println("\n=== Security Update ===")
		fmt.Println("Your config.yaml file has been updated with hashed passwords.")
		fmt.Println("Please use the new hashed passwords for authentication.")
	}

	return nil
}

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
