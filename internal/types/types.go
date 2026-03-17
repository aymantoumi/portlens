package types

// SourceKind describes what kind of entity owns a port.
type SourceKind string

const (
	SourceDocker        SourceKind = "docker"
	SourceDockerCompose SourceKind = "compose"
	SourceProcess       SourceKind = "process"
	SourceSystemd       SourceKind = "systemd"
	SourceUnknown       SourceKind = "unknown"
)

// BadgeKind describes how conventional a port assignment is.
type BadgeKind string

const (
	BadgeOfficial   BadgeKind = "official"
	BadgeConvention BadgeKind = "convention"
	BadgeCustom     BadgeKind = "custom"
	BadgeFree       BadgeKind = "free"
)

// PortEntry is the unified identity model for a single port binding.
// Every scanner (docker, proc, systemd) produces PortEntry values.
// The resolver merges and enriches them from the registry.
type PortEntry struct {
	Port   int
	Proto  string // "tcp" or "udp"
	BindIP string // "0.0.0.0", "127.0.0.1", "::"

	Kind SourceKind

	// Docker / Compose fields. Non-empty when Kind is docker or compose.
	ContainerName  string
	ContainerID    string // 12-char short ID
	ImageName      string // e.g. "postgres:16"
	ComposeProject string // label: com.docker.compose.project
	ComposeService string // label: com.docker.compose.service
	ComposeFile    string // label: com.docker.compose.project.config_files

	// Process fields. Non-empty when Kind is process or systemd.
	PID         int
	ProcessName string // argv[0] basename
	ExePath     string // resolved /proc/[pid]/exe symlink
	CmdLine     string // full command line, space-joined
	Username    string // resolved from uid via os/user
	SystemdUnit string // e.g. "nginx.service", empty if not managed

	// Registry enrichment. Set by resolver after registry lookup.
	RegistryName string    // human name, e.g. "PostgreSQL"
	Category     string    // e.g. "database"
	Badge        BadgeKind // official / convention / custom

	// Conflict. Non-nil when another PortEntry claims the same port.
	ConflictWith *PortEntry
}

// ScanResult holds all port entries from a full scan.
type ScanResult struct {
	Entries   []*PortEntry
	Conflicts []*ConflictPair
	Duration  string
}

// ConflictPair holds two entries that conflict on the same port.
type ConflictPair struct {
	Port int
	A    *PortEntry
	B    *PortEntry
}

// FreeResult holds available ports from a free scan.
type FreeResult struct {
	Range     PortRange
	Available []int
	Blocked   []BlockedPort
}

// PortRange is a min/max port range.
type PortRange struct {
	Min int
	Max int
}

// BlockedPort is a port within a requested range that is not free.
type BlockedPort struct {
	Port   int
	Reason string // e.g. "in use by grafana-1", "reserved by convention"
}

// RegistryEntry is a single record from the port registry.
type RegistryEntry struct {
	Port     int
	Name     string
	Category string
	Desc     string
	Badge    BadgeKind
}
