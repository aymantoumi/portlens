//go:build windows
// +build windows

package platform

import (
	"fmt"

	"github.com/aymantoumi/portlens/internal/types"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type WindowsScanner struct{}

func NewScanner() Scanner {
	return &WindowsScanner{}
}

func (s *WindowsScanner) ScanListening() ([]*types.PortEntry, error) {
	entries := make([]*types.PortEntry, 0)

	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	pidConnections := make(map[int32][]net.ConnectionStat)
	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			pidConnections[conn.Pid] = append(pidConnections[conn.Pid], conn)
		}
	}

	for pid, conns := range pidConnections {
		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			continue
		}

		name, _ := proc.Name()
		cmdline, _ := proc.CmdlineSlice()
		username, _ := proc.Username()

		exePath, _ := proc.Exe()

		for _, conn := range conns {
			entry := &types.PortEntry{
				Port:        int(conn.Laddr.Port),
				Proto:       "tcp",
				BindIP:      conn.Laddr.IP,
				Kind:        types.SourceProcess,
				PID:         int(pid),
				ProcessName: name,
				ExePath:     exePath,
				Username:    username,
			}

			if len(cmdline) > 0 {
				entry.ExePath = cmdline[0]
			}

			entry.Kind = types.SourceProcess

			entries = append(entries, entry)
		}
	}

	return entries, nil
}
