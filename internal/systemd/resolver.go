package systemd

import (
	"fmt"
	"os"
	"strings"
)

type UnitInfo struct {
	UnitName string
	State    string
}

func ResolveByPID(pid int) *UnitInfo {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		return nil
	}

	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, "system.slice") {
			continue
		}
		parts := strings.Split(line, "/")
		for _, part := range parts {
			if strings.HasSuffix(part, ".service") ||
				strings.HasSuffix(part, ".socket") ||
				strings.HasSuffix(part, ".timer") {
				return &UnitInfo{
					UnitName: part,
					State:    "active",
				}
			}
		}
	}
	return nil
}
