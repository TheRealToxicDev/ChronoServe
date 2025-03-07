package windows

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/toxic-development/sysmanix/services"
)

// Ensure WindowsService implements ServiceHandler
var _ services.ServiceHandler = (*WindowsService)(nil)

// WindowsService implements the ServiceHandler interface for Windows
type WindowsService struct {
	services.BaseServiceHandler
	cache             map[string]services.ServiceStatus
	cacheMutex        sync.RWMutex
	cacheTTL          time.Duration
	protectedServices map[string]bool
}

// NewWindowsService creates a new Windows service handler
func NewWindowsService() *WindowsService {
	// Initialize the protected services map
	protectedServices := make(map[string]bool)
	for _, service := range criticalWindowsServices {
		protectedServices[strings.ToLower(service)] = true
	}

	return &WindowsService{
		cache:             make(map[string]services.ServiceStatus),
		cacheTTL:          5 * time.Minute,
		protectedServices: protectedServices,
	}
}

// IsProtectedService checks if a service is in the protected list
func (s *WindowsService) IsProtectedService(name string) bool {
	return s.protectedServices[strings.ToLower(name)]
}

// ValidateServiceOperation checks if operations are allowed on this service
func (s *WindowsService) ValidateServiceOperation(name string) error {
	if !s.ValidateServiceName(name) {
		return fmt.Errorf("invalid service name")
	}

	if s.IsProtectedService(name) {
		return fmt.Errorf("operation not allowed on protected system service: %s", name)
	}

	return nil
}
