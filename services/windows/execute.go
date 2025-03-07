package windows

import (
	"bytes"
	"fmt"
	"os/exec"
)

// executePowershell executes a PowerShell script and returns its output
func (s *WindowsService) executePowershell(script string) (*bytes.Buffer, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

	return &out, nil
}
