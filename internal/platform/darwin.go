//go:build darwin
// +build darwin

package platform

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aymantoumi/portlens/internal/types"
)

type DarwinScanner struct{}

func NewScanner() Scanner {
	return &DarwinScanner{}
}

func (s *DarwinScanner) ScanListening() ([]*types.PortEntry, error) {
	entries := make([]*types.PortEntry, 0)

	cmd := exec.Command("lsof", "-i", "-P", "-n", "+c", "0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsof: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		processName := fields[0]
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		nameField := fields[8]
		if !strings.Contains(nameField, "(LISTEN)") && !strings.Contains(nameField, "LISTEN") {
			continue
		}

		var port int
		var ip string
		for i := len(fields) - 1; i >= 0; i-- {
			if strings.Contains(fields[i], "->") {
				continue
			}
			if strings.Contains(fields[i], ":") {
				parts := strings.Split(fields[i], ":")
				portStr := parts[len(parts)-1]
				port, err = strconv.Atoi(strings.Split(portStr, "-")[0])
				if err != nil {
					continue
				}
				if len(parts) > 1 {
					ip = strings.Join(parts[:len(parts)-1], ":")
				}
				break
			}
		}

		if port == 0 {
			continue
		}

		entry := &types.PortEntry{
			Port:        port,
			Proto:       "tcp",
			BindIP:      ip,
			Kind:        types.SourceProcess,
			PID:         pid,
			ProcessName: processName,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
