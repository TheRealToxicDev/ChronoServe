package linux

import (
	"bytes"
	"fmt"
	"os/exec"
)

// executeCommand executes a command and returns its output
func (s *SystemdService) executeCommand(cmd *exec.Cmd) (*bytes.Buffer, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}

	return &out, nil
}
