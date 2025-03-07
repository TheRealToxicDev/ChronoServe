package config

import (
	"fmt"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// Config represents the root configuration structure
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Auth    AuthConfig    `yaml:"auth"`
	Linux   LinuxConfig   `yaml:"linux,omitempty"`   // Use omitempty to exclude empty fields
	Windows WindowsConfig `yaml:"windows,omitempty"` // Use omitempty to exclude empty fields
	Logging LogConfig     `yaml:"logging"`
	API     APIConfig     `yaml:"api"`
}

// APIConfig contains API-related settings including Swagger documentation
type APIConfig struct {
	EnableSwagger bool   `yaml:"enableSwagger"` // Whether to enable Swagger UI
	SwaggerPath   string `yaml:"swaggerPath"`   // Path to Swagger UI
	Version       string `yaml:"version"`       // API version
	Title         string `yaml:"title"`         // API title for documentation
	Description   string `yaml:"description"`   // API description for documentation
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
	ID           string   `yaml:"id"` // Unique user ID
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password,omitempty"` // Temporary field for plain password
	PasswordHash string   `yaml:"-"`                  // Exclude from initial config creation
	Roles        []string `yaml:"roles"`
	AvatarURL    string   `yaml:"avatarUrl,omitempty"`   // Optional avatar URL
	BannerURL    string   `yaml:"bannerUrl,omitempty"`   // Optional banner URL
	DisplayName  string   `yaml:"displayName,omitempty"` // Optional display name
	Bio          string   `yaml:"bio,omitempty"`         // Optional bio
}

func (c *Credentials) SetPassword(password string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	c.PasswordHash = hash
	return nil
}

func (c *Credentials) VerifyPassword(password string) (bool, error) {
	return utils.VerifyPassword(password, c.PasswordHash)
}

type LinuxConfig struct {
	ServiceCommand string `yaml:"serviceCommand"`
	LogDirectory   string `yaml:"logDirectory"`
}

type WindowsConfig struct {
	ServiceCommand string `yaml:"serviceCommand"`
	LogDirectory   string `yaml:"logDirectory"`
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
