package config

import (
	"fmt"
	"time"

	"github.com/therealtoxicdev/chronoserve/utils"
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
	Username     string   `yaml:"username"`
	Password     string   `yaml:"password,omitempty"` // Temporary field for plain password
	PasswordHash string   `yaml:"password_hash"`      // Stored hash
	Roles        []string `yaml:"roles"`
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
