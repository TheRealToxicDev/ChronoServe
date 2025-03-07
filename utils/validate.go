package utils

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// Validation related constants
const (
	MaxServiceNameLength = 100
	MaxUsernameLength    = 50
	MaxPasswordLength    = 100
)

var (
	serviceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	pathRegex        = regexp.MustCompile(`^[a-zA-Z0-9\-_./\\]+$`)

	// ServiceNamePattern defines the allowed characters in service names
	ServiceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

	// UsernamePattern defines the allowed characters in usernames
	UsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.@]+$`)
)

func ValidatePath(path string) bool {
	return pathRegex.MatchString(path) && len(path) <= 4096
}

func SanitizePath(path string) string {
	// Remove any potential directory traversal attempts
	path = strings.Replace(path, "..", "", -1)
	// Convert windows path separators to unix style
	path = strings.ReplaceAll(path, "\\", "/")
	// Remove multiple consecutive slashes
	path = regexp.MustCompile(`/+`).ReplaceAllString(path, "/")
	return path
}

// ExtractServiceName extracts the service name from a URL path
// Expected format: /services/action/name
func ExtractServiceName(urlPath string) string {
	// Get the last part of the path
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(parts) < 1 {
		return ""
	}
	return parts[len(parts)-1]
}

// ParseInt parses a string into an integer with error handling
func ParseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string cannot be parsed as integer")
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as integer: %v", s, err)
	}

	return i, nil
}

// ValidateServiceName validates that a service name meets requirements:
// - Not empty
// - Only contains allowed characters (alphanumeric, underscore, dash, dot)
// - Not too long
func ValidateServiceName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	if len(name) > MaxServiceNameLength {
		return fmt.Errorf("service name cannot exceed %d characters", MaxServiceNameLength)
	}

	if !ServiceNamePattern.MatchString(name) {
		return fmt.Errorf("service name contains invalid characters")
	}

	return nil
}

// ValidateUsername validates that a username meets requirements:
// - Not empty
// - Only contains allowed characters
// - Not too long
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) > MaxUsernameLength {
		return fmt.Errorf("username cannot exceed %d characters", MaxUsernameLength)
	}

	if !UsernamePattern.MatchString(username) {
		return fmt.Errorf("username contains invalid characters")
	}

	return nil
}

// ValidatePassword validates that a password meets basic requirements:
// - Not empty
// - Not too long
// Note: This doesn't enforce password complexity - that should be handled by the UI/API
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password cannot exceed %d characters", MaxPasswordLength)
	}

	return nil
}

// ValidateLogQueryParams validates log query parameters
func ValidateLogQueryParams(lines int) (int, error) {
	// Default to 100 lines if not specified or invalid
	if lines <= 0 {
		return 100, nil
	}

	// Cap at 1000 lines to prevent resource exhaustion
	if lines > 1000 {
		return 1000, fmt.Errorf("maximum allowed log lines is 1000")
	}

	return lines, nil
}

// ValidateRole checks if a role is allowed in the system
func ValidateRole(role string, allowedRoles []string) bool {
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}

// ValidateRoles checks if all roles are allowed in the system
func ValidateRoles(roles []string, allowedRoles []string) error {
	for _, role := range roles {
		if !ValidateRole(role, allowedRoles) {
			return fmt.Errorf("invalid role: %s", role)
		}
	}
	return nil
}

// IsValidPath checks if a path is valid and safe
func IsValidPath(path string) bool {
	// Disallow empty paths
	if path == "" {
		return false
	}

	// Disallow paths with ".." to prevent directory traversal
	if strings.Contains(path, "..") {
		return false
	}

	// Additional path safety checks can be added here

	return true
}

// SanitizeFilePath sanitizes and validates a file path
func SanitizeFilePath(basePath, filePath string) (string, error) {
	// Prevent directory traversal attacks
	cleanPath := path.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	fullPath := path.Join(basePath, cleanPath)

	// Make sure the resulting path is still under the base path
	if !strings.HasPrefix(fullPath, basePath) {
		return "", fmt.Errorf("path escapes base directory")
	}

	return fullPath, nil
}
