package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	// Version is the current version of SysManix
	Version = "0.1.0"

	// githubAPI is the endpoint for checking latest releases
	githubAPI = "https://api.github.com/repos/toxic-development/sysmanix/releases/latest"

	// Version information
	AppName    = "SysManix"
	AppVersion = "0.1.0"
	BuildDate  = "2023-06-01"
)

var (
	// StartTime records when the application was started
	StartTime = time.Now()
)

type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}

// VersionInfo contains version information about the application
type VersionInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	BuildDate string    `json:"buildDate"`
	GoVersion string    `json:"goVersion"`
	OS        string    `json:"os"`
	Arch      string    `json:"arch"`
	StartTime time.Time `json:"startTime"`
	Uptime    string    `json:"uptime"`
}

// GetVersionInfo returns the current version information
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Name:      AppName,
		Version:   AppVersion,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		StartTime: StartTime,
		Uptime:    GetUptime(),
	}
}

// GetUptime returns the application uptime as a human-readable string
func GetUptime() string {
	uptime := time.Since(StartTime)
	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds",
		int(uptime.Hours())/24,
		int(uptime.Hours())%24,
		int(uptime.Minutes())%60,
		int(uptime.Seconds())%60)
}

// CheckVersion compares current version with latest GitHub release
func CheckVersion() (isLatest bool, latestVersion string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", githubAPI, nil)
	if err != nil {
		return false, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "SysManix-Version-Checker")

	resp, err := client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("failed to fetch latest version: HTTP %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Clean version strings (remove 'v' prefix if present)
	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(Version, "v")

	isLatest = compareVersions(current, latest) >= 0
	return isLatest, latest, nil
}

// compareVersions compares two semantic version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Pad shorter version with zeros
	for len(parts1) < 3 {
		parts1 = append(parts1, "0")
	}
	for len(parts2) < 3 {
		parts2 = append(parts2, "0")
	}

	// Compare version parts
	for i := 0; i < 3; i++ {
		if parts1[i] != parts2[i] {
			num1 := parseVersionPart(parts1[i])
			num2 := parseVersionPart(parts2[i])
			if num1 < num2 {
				return -1
			}
			if num1 > num2 {
				return 1
			}
		}
	}
	return 0
}

// parseVersionPart safely converts version string parts to integers
func parseVersionPart(part string) int {
	var num int
	_, err := fmt.Sscanf(part, "%d", &num)
	if err != nil {
		return 0
	}
	return num
}

// GetCurrentVersion returns the current version string
func GetCurrentVersion() string {
	return Version
}

// PrintVersionInfo prints version information to stdout
func PrintVersionInfo() {
	isLatest, latestVersion, err := CheckVersion()
	if err != nil {
		fmt.Printf("SysManix v%s\nFailed to check for updates: %v\n", Version, err)
		return
	}

	if isLatest {
		fmt.Printf("SysManix v%s (latest)\n", Version)
	} else {
		fmt.Printf("SysManix v%s (update available: v%s)\n", Version, latestVersion)
		fmt.Println("Visit https://github.com/toxic-development/sysmanix/releases for the latest version")
	}
}

func CheckVersionInBackground(logger *Logger) {
	go func() {
		// Initial check after short delay
		time.Sleep(5 * time.Second)
		checkAndLogVersion(logger)

		// Periodic checks every 24 hours
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			checkAndLogVersion(logger)
		}
	}()
}

func checkAndLogVersion(logger *Logger) {
	isLatest, latestVersion, err := CheckVersion()
	if err != nil {
		logger.Warn("Failed to check for updates: %v", err)
		return
	}

	if !isLatest {
		logger.Info("Update available: v%s (current: v%s)", latestVersion, Version)
		logger.Info("Visit https://github.com/toxic-development/sysmanix/releases for the latest version")
	}
}
