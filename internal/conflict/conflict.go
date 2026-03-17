package conflict

import (
	"fmt"

	"github.com/yourusername/portmap/internal/registry"
	"github.com/yourusername/portmap/internal/types"
)

type Suggestion struct {
	Port        int
	Description string
	Alternative int
	FixCommand  string
}

func Suggest(pair *types.ConflictPair) Suggestion {
	s := Suggestion{Port: pair.Port}

	a := pair.A
	b := pair.B

	s.Description = fmt.Sprintf(
		"%s and %s both bound to port %d",
		identLabel(a), identLabel(b), pair.Port,
	)

	alt := findAlternative(pair.Port)
	if alt != 0 {
		s.Alternative = alt
		s.FixCommand = fmt.Sprintf("portlens suggest --fix %d", pair.Port)
	}

	return s
}

func identLabel(e *types.PortEntry) string {
	switch e.Kind {
	case types.SourceDocker, types.SourceDockerCompose:
		return fmt.Sprintf("container %s", e.ContainerName)
	case types.SourceProcess, types.SourceSystemd:
		if e.SystemdUnit != "" {
			return fmt.Sprintf("systemd unit %s", e.SystemdUnit)
		}
		return fmt.Sprintf("process %s (pid %d)", e.ProcessName, e.PID)
	default:
		return fmt.Sprintf("unknown process on port %d", e.Port)
	}
}

func findAlternative(port int) int {
	candidates := []int{port + 10, port + 1, port + 100}
	for _, c := range candidates {
		e := registry.Lookup(c)
		if e != nil {
			return c
		}
	}
	return 0
}
