package proc

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aymantoumi/portlens/internal/types"
)

func ScanListening() ([]*types.PortEntry, error) {
	var entries []*types.PortEntry

	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		parsed, err := parseProcNet(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("proc: parse %s: %w", path, err)
		}
		entries = append(entries, parsed...)
	}

	return entries, nil
}

func parseProcNet(path string) ([]*types.PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	inodeToPort := map[string]int{}
	scanner := bufio.NewScanner(f)
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}
		state := fields[3]
		if state != "0A" {
			continue
		}
		localAddr := fields[1]
		inode := fields[9]

		port, err := parseHexPort(localAddr)
		if err != nil {
			continue
		}
		inodeToPort[inode] = port
	}

	var entries []*types.PortEntry
	for inode, port := range inodeToPort {
		e := buildEntry(port, inode)
		entries = append(entries, e)
	}
	return entries, nil
}

func parseHexPort(addr string) (int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid addr: %s", addr)
	}
	port64, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, fmt.Errorf("parse port hex %s: %w", parts[1], err)
	}
	return int(port64), nil
}

func buildEntry(port int, inode string) *types.PortEntry {
	e := &types.PortEntry{
		Port:   port,
		Proto:  "tcp",
		BindIP: "0.0.0.0",
		Kind:   types.SourceUnknown,
	}

	pid, err := inodeToPID(inode)
	if err != nil || pid == 0 {
		return e
	}

	e.Kind = types.SourceProcess
	e.PID = pid
	e.ProcessName = readProcessName(pid)
	e.ExePath = readExePath(pid)
	e.CmdLine = readCmdLine(pid)
	e.Username = readUsername(pid)

	return e
}

func inodeToPID(inode string) (int, error) {
	target := fmt.Sprintf("socket:[%s]", inode)

	procs, err := filepath.Glob("/proc/[0-9]*/fd/*")
	if err != nil {
		return 0, err
	}
	for _, fd := range procs {
		link, err := os.Readlink(fd)
		if err != nil {
			continue
		}
		if link == target {
			parts := strings.Split(fd, "/")
			if len(parts) < 3 {
				continue
			}
			pid, err := strconv.Atoi(parts[2])
			if err != nil {
				continue
			}
			return pid, nil
		}
	}
	return 0, nil
}

func readProcessName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readExePath(pid int) string {
	link, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return ""
	}
	return link
}

func readCmdLine(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return ""
	}
	parts := strings.Split(string(data), "\x00")
	var clean []string
	for _, p := range parts {
		if p != "" {
			clean = append(clean, p)
		}
	}
	return strings.Join(clean, " ")
}

func readUsername(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "Uid:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return ""
			}
			u, err := user.LookupId(fields[1])
			if err != nil {
				return fields[1]
			}
			return u.Username
		}
	}
	return ""
}
