package utils

import (
	"regexp"
	"strings"
)

var (
	serviceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	pathRegex        = regexp.MustCompile(`^[a-zA-Z0-9\-_./\\]+$`)
)

func ValidateServiceName(name string) bool {
	return serviceNameRegex.MatchString(name) && len(name) <= 255
}

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

func ExtractServiceName(path string) string {
	// Extract the service name from the URL path
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
