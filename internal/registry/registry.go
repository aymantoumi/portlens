package registry

import "github.com/yourusername/portmap/internal/types"

// Lookup returns the RegistryEntry for a port number.
// Returns nil if the port is not in the registry.
func Lookup(port int) *types.RegistryEntry {
	for i := range entries {
		if entries[i].Port == port {
			return &entries[i]
		}
	}
	return nil
}

// All returns a copy of the full registry slice.
func All() []types.RegistryEntry {
	out := make([]types.RegistryEntry, len(entries))
	copy(out, entries)
	return out
}

// ByCategory returns all entries matching a category string.
func ByCategory(cat string) []types.RegistryEntry {
	var out []types.RegistryEntry
	for _, e := range entries {
		if e.Category == cat {
			out = append(out, e)
		}
	}
	return out
}

// IsConventionPort returns true if the port is in the registry
// with official or convention badge. Used by portmap free to mark
// ports that should not be suggested even if technically unbound.
func IsConventionPort(port int) bool {
	e := Lookup(port)
	if e == nil {
		return false
	}
	return e.Badge == types.BadgeOfficial || e.Badge == types.BadgeConvention
}
