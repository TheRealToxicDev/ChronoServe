package linux

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/toxic-development/sysmanix/services"
)

// Ensure SystemdService implements ServiceHandler
var _ services.ServiceHandler = (*SystemdService)(nil)

// SystemdService implements the ServiceHandler interface for Linux
type SystemdService struct {
	services.BaseServiceHandler
	cache             map[string]services.ServiceStatus
	cacheMutex        sync.RWMutex
	cacheTTL          time.Duration
	protectedServices map[string]bool
}

// NewSystemdService creates a new systemd service handler
func NewSystemdService() *SystemdService {
	// Initialize the protected services map
	protectedServices := make(map[string]bool)
	for _, service := range criticalLinuxServices {
		protectedServices[strings.ToLower(service)] = true
	}

	return &SystemdService{
		cache:             make(map[string]services.ServiceStatus),
		cacheTTL:          5 * time.Minute,
		protectedServices: protectedServices,
	}
}

// IsProtectedService checks if a service is in the protected list
func (s *SystemdService) IsProtectedService(name string) bool {
	// First check for exact matches
	if s.protectedServices[strings.ToLower(name)] {
		return true
	}

	// Then check for prefix matches (for services like user@1000.service)
	for prefix := range s.protectedServices {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			return true
		}
	}

	return false
}

// ValidateServiceOperation checks if operations are allowed on this service
func (s *SystemdService) ValidateServiceOperation(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	if s.IsProtectedService(name) {
		return fmt.Errorf("operation not allowed on protected system service: %s", name)
	}

	return nil
}
