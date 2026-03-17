package resolver

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/yourusername/portmap/internal/docker"
	"github.com/yourusername/portmap/internal/proc"
	"github.com/yourusername/portmap/internal/registry"
	"github.com/yourusername/portmap/internal/systemd"
	"github.com/yourusername/portmap/internal/types"
)

type Options struct {
	SkipDocker bool
	SkipProc   bool
	Proto      string
}

func Run(opts Options) (*types.ScanResult, error) {
	start := time.Now()
	byPort := map[int]*types.PortEntry{}

	if !opts.SkipDocker {
		dc, err := docker.NewClient()
		if err == nil {
			defer dc.Close()
			dockerEntries, err := dc.ScanPorts(context.Background())
			if err != nil {
				return nil, fmt.Errorf("resolver: docker scan: %w", err)
			}
			for _, e := range dockerEntries {
				enrichEntry(e)
				byPort[e.Port] = e
			}
		}
	}

	if !opts.SkipProc {
		procEntries, err := proc.ScanListening()
		if err != nil {
			return nil, fmt.Errorf("resolver: proc scan: %w", err)
		}
		for _, e := range procEntries {
			enrichEntry(e)
			enrichSystemd(e)

			existing, ok := byPort[e.Port]
			if !ok {
				byPort[e.Port] = e
				continue
			}
			if existing.Kind == types.SourceDockerCompose ||
				existing.Kind == types.SourceDocker {
				// Only conflict if proc has a real process with valid PID
				// Don't conflict if PID is 0 (unresolved) or is docker-proxy
				if e.PID > 0 && e.ProcessName != "docker-proxy" {
					existing.ConflictWith = e
				}
			}
		}
	}

	var entries []*types.PortEntry
	for _, e := range byPort {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Port < entries[j].Port
	})

	var conflicts []*types.ConflictPair
	for _, e := range entries {
		if e.ConflictWith != nil {
			conflicts = append(conflicts, &types.ConflictPair{
				Port: e.Port,
				A:    e,
				B:    e.ConflictWith,
			})
		}
	}

	return &types.ScanResult{
		Entries:   entries,
		Conflicts: conflicts,
		Duration:  fmt.Sprintf("%.1fs", time.Since(start).Seconds()),
	}, nil
}

func FreeInRange(min, max, count int) (*types.FreeResult, error) {
	result, err := Run(Options{})
	if err != nil {
		return nil, err
	}

	bound := map[int]string{}
	for _, e := range result.Entries {
		bound[e.Port] = describeEntry(e)
	}

	var available []int
	var blocked []types.BlockedPort

	for p := min; p <= max && len(available) < count; p++ {
		if reason, ok := bound[p]; ok {
			blocked = append(blocked, types.BlockedPort{Port: p, Reason: reason})
			continue
		}
		if registry.IsConventionPort(p) {
			blocked = append(blocked, types.BlockedPort{
				Port:   p,
				Reason: "reserved by portlens convention",
			})
			continue
		}
		available = append(available, p)
	}

	return &types.FreeResult{
		Range:     types.PortRange{Min: min, Max: max},
		Available: available,
		Blocked:   blocked,
	}, nil
}

func enrichEntry(e *types.PortEntry) {
	reg := registry.Lookup(e.Port)
	if reg == nil {
		return
	}
	e.RegistryName = reg.Name
	e.Category = reg.Category
	e.Badge = reg.Badge
}

func enrichSystemd(e *types.PortEntry) {
	if e.PID == 0 {
		return
	}
	unit := systemd.ResolveByPID(e.PID)
	if unit == nil {
		return
	}
	e.SystemdUnit = unit.UnitName
	e.Kind = types.SourceSystemd
}

func describeEntry(e *types.PortEntry) string {
	if e.ContainerName != "" {
		return "in use by container " + e.ContainerName
	}
	if e.ProcessName != "" {
		return fmt.Sprintf("in use by process %s (pid %d)", e.ProcessName, e.PID)
	}
	return "in use"
}
