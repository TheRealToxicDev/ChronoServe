package utils

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	serviceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	pathRegex        = regexp.MustCompile(`^[a-zA-Z0-9\-_./\\]+$`)
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

// ValidateServiceName checks if a service name is safe to use
// This helps prevent command injection
func ValidateServiceName(name string) bool {
	if name == "" {
		return false
	}

	// Only allow alphanumeric characters, dashes, underscores, and periods
	// This is a common pattern for service names
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-_.]+$`, name)
	return matched
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
