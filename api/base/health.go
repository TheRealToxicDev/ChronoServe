package base

import (
	"net/http"
	"time"

	"github.com/toxic-development/sysmanix/utils"
)

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
	GoVersion string    `json:"goVersion"`
	Memory    MemStats  `json:"memory"`
	StartTime time.Time `json:"startTime"`
}

type MemStats struct {
	Alloc      uint64  `json:"alloc"`
	TotalAlloc uint64  `json:"totalAlloc"`
	Sys        uint64  `json:"sys"`
	NumGC      uint32  `json:"numGC"`
	HeapInUse  float64 `json:"heapInUse"`
}

// @Summary      Get service health
// @Description  Returns service health information and status
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthStatus
// @Router       /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.HealthCheck(w, r)
}
