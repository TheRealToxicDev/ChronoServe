package utils

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
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

var startTime = time.Now()

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	status := HealthStatus{
		Status:    "healthy",
		Uptime:    time.Since(startTime).String(),
		GoVersion: runtime.Version(),
		Memory: MemStats{
			Alloc:      mem.Alloc,
			TotalAlloc: mem.TotalAlloc,
			Sys:        mem.Sys,
			NumGC:      mem.NumGC,
			HeapInUse:  float64(mem.HeapInuse) / 1024 / 1024, // MB
		},
		StartTime: startTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
