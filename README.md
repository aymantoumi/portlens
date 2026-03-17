# portlens

[![Cross-Platform](https://img.shields.io/badge/cross--platform-Linux%20%7C%20Windows%20%7C%20macOS-blue.svg)](https://github.com/aymantoumi/portlens)
[![Beta](https://img.shields.io/badge/status-beta-orange.svg)](https://github.com/aymantoumi/portlens)

A CLI tool that shows you what is running on every port on your machine, what Docker container or system process owns it, and whether anything is conflicting.

## Why portlens?

Running multiple Docker projects alongside system services means port collisions happen constantly. `docker compose up` fails, a dev server silently uses the wrong port, or Grafana and your React app both try to bind 3000. Finding the owner of a port requires combining `docker ps`, `netstat`, `lsof`, and `cat /proc/net/tcp` by hand.

portlens does all of that in one command.

## Features

- Scan all port bindings from Docker containers, system processes, and systemd services
- Detect port conflicts between containers and host processes
- Find available ports in any range
- Get suggested port layouts for common stacks (database, backend, frontend, etc.)
- Monitor ports in real-time
- Export port configurations for docker-compose, .env, or JSON
- Built-in registry of 75+ well-known ports
- Color-coded output by service category
- **Cross-platform**: Linux, Windows, and macOS support

## Install

### From source

```bash
go install github.com/aymantoumi/portlens@latest
```

### Binary download (Linux amd64 and arm64)

Download the latest release from the [releases page](https://github.com/aymantoumi/portlens/releases).

```bash
curl -L https://github.com/aymantoumi/portlens/releases/latest/download/portlens-linux-amd64 \
    -o portlens
chmod +x portlens
sudo mv portlens /usr/local/bin/
```

### Build from source

```bash
git clone https://github.com/aymantoumi/portlens
cd portlens
go build -o portlens .
```

### Requirements

- Linux, Windows, or macOS
- Docker daemon socket at /var/run/docker.sock (optional, skip with --skip-docker)
- No root required for Docker and most process scanning. Some unresolved processes may need elevated permissions.

---

## Platform Support

portlens is cross-platform and supports:

| Platform | Status | Method |
|----------|--------|--------|
| Linux | ✅ Stable | Reads `/proc/net/tcp` and `/proc/[pid]/` |
| Windows | �_beta_ | Uses gopsutil library |
| macOS | �_beta_ | Uses `lsof` command |

### Building for Different Platforms

```bash
# Linux (default)
go build -o portlens .

# Windows
GOOS=windows GOARCH=amd64 go build -o portlens.exe .

# macOS
GOOS=darwin GOARCH=amd64 go build -o portlens .
```

### Beta Notice

Windows and macOS support are currently in **beta**. Please report any issues at: https://github.com/aymantoumi/portlens/issues

---

## Commands

### scan

Show every active port binding on the machine.

```bash
portlens scan
```

For each port, portlens shows:
- The port number and service name from the registry
- Whether it is a Docker container, Docker Compose service, system process, or systemd unit
- The container name and image, or the process name and PID
- Any conflicts where two services claim the same port

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON for scripting |
| `--skip-docker` | Skip Docker socket scan |
| `--skip-proc` | Skip /proc scan |
| `--category` | Filter by category (database, api, web, etc.) |

**Examples:**

```bash
# Scan all ports
portlens scan

# Output as JSON
portlens scan --json

# Filter by category
portlens scan --category database
portlens scan --category api

# Skip specific sources
portlens scan --skip-docker
portlens scan --skip-proc
```

**Exit codes:**
- `0` - No conflicts found
- `1` - Conflicts detected (useful for CI scripts)

---

### free

Find available (free) ports in a range.

```bash
portlens free
portlens free 4000-5000
portlens free 4000-5000 --count 3
```

portlens checks which ports are unbound and not reserved by convention in the registry. Ports like 3010 (Grafana) are excluded from suggestions even if nothing is running there, because deploying to them later would cause a conflict.

**Flags:**

| Flag | Description |
|------|-------------|
| `--count` | Number of free ports to return (default: 5) |
| `--format` | Output format: table, json, list |

**Examples:**

```bash
# Default range (3000-9999)
portlens free

# Custom range
portlens free 4000-5000

# Specific number of ports
portlens free 4000-5000 --count 3

# Output formats
portlens free --format table   # grid layout (default)
portlens free --format list    # plain list
portlens free --format json    # JSON output
```

---

### suggest

Show a recommended port layout for a named service stack.

```bash
portlens suggest
portlens suggest --stack database
portlens suggest --stack fullstack
```

The output cross-references the suggestion against your current scan. Ports that are already in use are annotated. Conflicts are flagged with a fix command.

**Available stacks:**

| Stack | Description |
|-------|-------------|
| `fullstack` | API + Frontend + Database |
| `backend` | Backend APIs |
| `frontend` | Frontend dev servers |
| `database` | DB + Cache + Search |
| `monitoring` | Observability tools |
| `homelab` | Home server apps |

**Flags:**

| Flag | Description |
|------|-------------|
| `--stack` | Service stack name |
| `--fix` | Suggest fix for a specific port |
| `--export` | Export format: compose, env, json |

**Examples:**

```bash
# Show all stacks
portlens suggest

# Specific stack
portlens suggest --stack database
portlens suggest --stack backend
portlens suggest --stack homelab

# Fix a conflict
portlens suggest --fix 3000

# Export formats
portlens suggest --stack fullstack --export compose   # docker-compose ports snippet
portlens suggest --stack fullstack --export env        # .env file
portlens suggest --stack fullstack --export json       # JSON array
```

---

### watch

Monitor port bindings in real-time.

```bash
portlens watch
portlens watch --interval 5
```

Each line in the log shows a timestamp, event type, port, service name, and source tag. New bindings are green (+). Released bindings are yellow (-). Conflicts are red (!).

When a conflict appears, portlens prints the full conflict block with a suggested fix immediately below the event line.

**Flags:**

| Flag | Description |
|------|-------------|
| `--interval` | Refresh interval in seconds (default: 2) |
| `--once` | Run once and exit |

**Examples:**

```bash
# Live monitoring
portlens watch

# Custom interval
portlens watch --interval 5

# Run once and exit
portlens watch --once
```

---

## Port Registry

portlens ships with a built-in registry of ~75 ports covering:

- **Databases**: PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch, etc.
- **APIs**: REST, GraphQL, gRPC, GraphQL
- **Frontend**: Vite, Webpack, React, Vue dev servers
- **Monitoring**: Grafana, Prometheus, Node Exporter
- **Message Queues**: RabbitMQ, Kafka, NATS
- **Storage**: MinIO, S3, NFS
- **Auth**: Keycloak, OAuth2, LDAP
- **Homelab**: Home Assistant, Pi-hole, AdGuard, Jellyfin

Each entry has a badge:

- **official** — IANA-registered port for this service
- **convention** — community default, not IANA-registered
- **custom** — portlens recommendation for commonly unassigned slots

The registry is used by all commands. `scan` shows the badge per port. `free` excludes convention ports from suggestions. `suggest` builds layouts from registry categories.

---

## Color System

portlens assigns color by service category, not by port number.

| Category | Color |
|----------|-------|
| database, cache, search | purple |
| api, auth | blue |
| frontend, web | green |
| storage, infra, homelab | orange |
| monitoring | yellow |
| conflicts | red |

Source type is shown as a small tag per row: `compose`, `docker`, `process`, `systemd`, or `unknown`.

---

## How It Works

portlens reads from three sources and merges the results:

### Docker Socket

Reads from `/var/run/docker.sock` to list all containers and their port bindings. For Docker Compose containers, it reads the `com.docker.compose.*` labels to identify the project, service name, and compose file path.

### Proc Filesystem

Scans `/proc/net/tcp` and `/proc/net/tcp6` to find all listening TCP sockets. Walks `/proc/[pid]/fd/` to match inodes to PIDs. Reads process name, exe path, full command line, and UID from the corresponding `/proc/[pid]/` files.

### Systemd

Reads `/proc/[pid]/cgroup` to detect whether a process is managed by a systemd unit. No D-Bus dependency.

### Conflict Detection

The resolver merges all three sources. Docker entries take priority for the same port. If a system process (other than docker-proxy) and a container both claim a port, that is flagged as a conflict.

---

## Building from Source

```bash
git clone https://github.com/aymantoumi/portlens
cd portlens
make build
./portlens --help
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the binary |
| `make build-all` | Build for Linux amd64 and arm64 |
| `make test` | Run tests |
| `make lint` | Run linters |
| `make clean` | Clean build artifacts |
| `make install` | Install to $GOPATH/bin |

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## Limitations

- Linux only. macOS and Windows do not expose `/proc/net/tcp`.
- Some process details (exe path, command line) require the process owner or root. portlens falls back to "unresolved" for those entries and suggests `sudo lsof -i :PORT`.
- UDP ports are not scanned. portlens covers TCP only.
- Docker must be running at `/var/run/docker.sock`. Remote Docker hosts are not supported.

---

## License

MIT
