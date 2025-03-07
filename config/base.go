package config

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

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

		// Remove OS-specific config based on the current system
		if runtime.GOOS == "linux" {
			config.Windows = WindowsConfig{}
		} else if runtime.GOOS == "windows" {
			config.Linux = LinuxConfig{}
		}

		// Assign a numeric ID to each user and exclude password_hash
		for username, user := range config.Auth.Users {
			userID := generateNumericID()
			user.ID = userID
			user.PasswordHash = "" // Ensure password_hash is excluded
			config.Auth.Users[username] = user
		}

		if err := saveConfigInternal(filePath); err != nil {
			return fmt.Errorf("error creating initial config: %w", err)
		}

		fmt.Printf("\n=== Configuration Notice ===\n")
		fmt.Printf("A new configuration file has been created at: %s\n", filePath)
		fmt.Println("Please update the following values in the config file:")
		fmt.Println("1. auth.secretKey")
		fmt.Println("2. auth.users.superadmin.password")
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

	// Assign IDs to users if not already set
	if err := assignUserIDs(); err != nil {
		return err
	}

	return config.Validate()
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
		// Allow either password or password_hash during validation
		if creds.Password == "" && creds.PasswordHash == "" {
			return fmt.Errorf("user %s has no password set", username)
		}
		if creds.PasswordHash != "" && !strings.HasPrefix(creds.PasswordHash, "$argon2id$") {
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

// processPlainPasswords hashes any plain text passwords in the configuration
func processPlainPasswords() error {
	configLock.Lock()
	defer configLock.Unlock()

	passwordsChanged := false

	// Check each user's password
	for username, user := range config.Auth.Users {
		// If there's a plain password
		if user.Password != "" {
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
		if err := saveConfigInternal(configPath); err != nil {
			return fmt.Errorf("failed to save hashed passwords: %w", err)
		}
		fmt.Println("\n=== Security Update ===")
		fmt.Println("Plain text passwords detected in the config.yaml file, which have been hashed.")
		fmt.Println("Please use the new hashed password for authentication with the API.")
	}

	return nil
}

// assignUserIDs assigns a numeric ID to users if not already set
func assignUserIDs() error {
	configLock.Lock()
	defer configLock.Unlock()

	idsAssigned := false

	for username, user := range config.Auth.Users {
		if user.ID == "" {
			userID := generateNumericID()
			user.ID = userID
			config.Auth.Users[username] = user
			idsAssigned = true
		}
	}

	// If any IDs were assigned, save the config
	if idsAssigned {
		if err := saveConfigInternal(configPath); err != nil {
			return fmt.Errorf("failed to save user IDs: %w", err)
		}
		fmt.Println("\n=== Configuration Update ===")
		fmt.Println("User IDs have been assigned to users without IDs in the config.yaml file.")
	}

	return nil
}

// generateNumericID generates a numeric ID similar to Discord
func generateNumericID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int63())
}

// Internal function to save config without acquiring locks
func saveConfigInternal(filePath string) error {
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
	configLock.Lock()
	defer configLock.Unlock()
	return saveConfigInternal(filePath)
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
		return saveConfigInternal(configPath)
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
