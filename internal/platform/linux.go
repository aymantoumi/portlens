//go:build linux
// +build linux

package platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aymantoumi/portlens/internal/types"
)

type LinuxScanner struct{}

func NewScanner() Scanner {
	return &LinuxScanner{}
}

func (s *LinuxScanner) ScanListening() ([]*types.PortEntry, error) {
	entries := make([]*types.PortEntry, 0)

	tcpPorts, err := parseProcNet("/proc/net/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to parse /proc/net/tcp: %w", err)
	}

	tcp6Ports, err := parseProcNet("/proc/net/tcp6")
	if err != nil {
		return nil, fmt.Errorf("failed to parse /proc/net/tcp6: %w", err)
	}

	ports := append(tcpPorts, tcp6Ports...)

	inodeToEntry := make(map[string]*types.PortEntry)

	for _, p := range ports {
		entry := &types.PortEntry{
			Port:     p.port,
			Proto:    "tcp",
			BindIP:   p.address,
			Kind:     types.SourceUnknown,
			PID:      0,
			Username: "",
		}

		if p.inode != "" {
			inodeToEntry[p.inode] = entry
		}
		entries = append(entries, entry)
	}

	if err := enrichWithPID(entries, inodeToEntry); err != nil {
		return nil, fmt.Errorf("failed to enrich with PID: %w", err)
	}

	return entries, nil
}

type procPort struct {
	port    int
	address string
	inode   string
	state   string
}

func parseProcNet(path string) ([]procPort, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return nil, nil
	}

	var ports []procPort

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localAddr := fields[1]
		inode := fields[9]
		state := fields[3]

		if state != "0A" {
			continue
		}

		parts := strings.Split(localAddr, ":")
		if len(parts) != 2 {
			continue
		}

		hexPort := parts[len(parts)-1]
		port, err := parseHexPort(hexPort)
		if err != nil {
			continue
		}

		addr := parts[0]
		if addr == "00000000" {
			addr = "0.0.0.0"
		} else if addr == "00000000000000000000000000000000" {
			addr = "::"
		} else if len(addr) == 32 {
			addr = fmt.Sprintf("[%s:%s]", addr[:32], hexPort)
		}

		ports = append(ports, procPort{
			port:    int(port),
			address: addr,
			inode:   inode,
			state:   state,
		})
	}

	return ports, nil
}

func parseHexPort(hexPort string) (int64, error) {
	return strconv.ParseInt(hexPort, 16, 64)
}

func enrichWithPID(entries []*types.PortEntry, inodeMap map[string]*types.PortEntry) error {
	procDir, err := os.Open("/proc")
	if err != nil {
		return err
	}
	defer procDir.Close()

	files, err := procDir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, name := range files {
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		fdDir, err := os.Open(fmt.Sprintf("/proc/%d/fd", pid))
		if err != nil {
			continue
		}

		fds, err := fdDir.Readdirnames(0)
		fdDir.Close()
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", pid, fd))
			if err != nil {
				continue
			}

			if strings.HasPrefix(link, "socket:[") {
				inode := strings.Trim(link, "socket:[]")
				if entry, ok := inodeMap[inode]; ok {
					entry.PID = pid
					entry.Kind = types.SourceProcess

					if name, err := readProcessName(pid); err == nil {
						entry.ProcessName = name
					}

					if cmdline, err := readCmdLine(pid); err == nil && len(cmdline) > 0 {
						entry.ExePath = cmdline[0]
					}

					if exe, err := readExePath(pid); err == nil {
						entry.ExePath = exe
					}

					if username, err := readUsername(pid); err == nil {
						entry.Username = username
					}

					if unit, err := readSystemdUnit(pid); err == nil && unit != "" {
						entry.SystemdUnit = unit
						entry.Kind = types.SourceSystemd
					}
				}
			}
		}
	}

	return nil
}

func readProcessName(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func readCmdLine(pid int) ([]string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	return strings.Split(string(data[:len(data)-1]), "\x00"), nil
}

func readExePath(pid int) (string, error) {
	return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
}

func readUsername(pid int) (string, error) {
	_, err := os.Stat(fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return "", err
	}
	return "", nil
}

func readSystemdUnit(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "systemd") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				name := parts[len(parts)-1]
				if strings.HasSuffix(name, ".service") {
					return strings.TrimSuffix(name, ".service"), nil
				}
				return name, nil
			}
		}
	}
	return "", nil
}
