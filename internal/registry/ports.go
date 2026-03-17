package registry

import "github.com/yourusername/portmap/internal/types"

var entries = []types.RegistryEntry{
	// Web and reverse proxy
	{Port: 80, Name: "HTTP", Category: "web", Desc: "Public web traffic", Badge: types.BadgeOfficial},
	{Port: 443, Name: "HTTPS", Category: "web", Desc: "Public TLS traffic", Badge: types.BadgeOfficial},
	{Port: 8080, Name: "HTTP alt", Category: "web", Desc: "Dev or proxy fallback", Badge: types.BadgeConvention},
	{Port: 8443, Name: "HTTPS alt", Category: "web", Desc: "Dev TLS fallback", Badge: types.BadgeConvention},
	{Port: 8888, Name: "Traefik dashboard", Category: "web", Desc: "Reverse proxy dashboard", Badge: types.BadgeConvention},
	{Port: 9000, Name: "Portainer", Category: "web", Desc: "Container management UI", Badge: types.BadgeConvention},
	{Port: 9090, Name: "Prometheus", Category: "monitoring", Desc: "Metrics scraping endpoint", Badge: types.BadgeConvention},
	{Port: 9443, Name: "Portainer TLS", Category: "web", Desc: "Portainer HTTPS interface", Badge: types.BadgeConvention},

	// Frontend and WebUI
	{Port: 3000, Name: "App UI", Category: "frontend", Desc: "React / Next.js / Vite dev server", Badge: types.BadgeConvention},
	{Port: 3001, Name: "App UI 2", Category: "frontend", Desc: "Second frontend or Storybook", Badge: types.BadgeCustom},
	{Port: 3002, Name: "Admin panel", Category: "frontend", Desc: "Internal admin WebUI", Badge: types.BadgeCustom},
	{Port: 3010, Name: "Grafana", Category: "monitoring", Desc: "Metrics dashboard WebUI", Badge: types.BadgeConvention},
	{Port: 3100, Name: "Loki", Category: "monitoring", Desc: "Log aggregation HTTP API", Badge: types.BadgeConvention},

	// Backend APIs
	{Port: 4000, Name: "GraphQL API", Category: "api", Desc: "Primary GraphQL endpoint", Badge: types.BadgeCustom},
	{Port: 4100, Name: "REST API", Category: "api", Desc: "Main backend REST service", Badge: types.BadgeCustom},
	{Port: 4200, Name: "REST API 2", Category: "api", Desc: "Secondary microservice", Badge: types.BadgeCustom},
	{Port: 4300, Name: "gRPC", Category: "api", Desc: "Internal RPC service", Badge: types.BadgeCustom},
	{Port: 5000, Name: "Flask/FastAPI", Category: "api", Desc: "Python backend default", Badge: types.BadgeConvention},
	{Port: 5001, Name: "FastAPI 2", Category: "api", Desc: "Second Python service", Badge: types.BadgeCustom},
	{Port: 5050, Name: "pgAdmin", Category: "database", Desc: "PostgreSQL web client", Badge: types.BadgeConvention},

	// Relational databases
	{Port: 5432, Name: "PostgreSQL", Category: "database", Desc: "Primary Postgres instance", Badge: types.BadgeOfficial},
	{Port: 5433, Name: "PostgreSQL 2", Category: "database", Desc: "Replica or second instance", Badge: types.BadgeCustom},
	{Port: 3306, Name: "MySQL/MariaDB", Category: "database", Desc: "MySQL official port", Badge: types.BadgeOfficial},
	{Port: 3307, Name: "MySQL 2", Category: "database", Desc: "Second MySQL instance", Badge: types.BadgeCustom},
	{Port: 1433, Name: "MSSQL", Category: "database", Desc: "SQL Server default port", Badge: types.BadgeOfficial},
	{Port: 1521, Name: "Oracle DB", Category: "database", Desc: "Oracle listener port", Badge: types.BadgeOfficial},

	// NoSQL and search
	{Port: 6379, Name: "Redis", Category: "cache", Desc: "Cache and session store", Badge: types.BadgeOfficial},
	{Port: 6380, Name: "Redis 2", Category: "cache", Desc: "Second Redis or TLS mode", Badge: types.BadgeCustom},
	{Port: 27017, Name: "MongoDB", Category: "database", Desc: "Document store default", Badge: types.BadgeOfficial},
	{Port: 9200, Name: "Elasticsearch", Category: "search", Desc: "Search HTTP API", Badge: types.BadgeOfficial},
	{Port: 9300, Name: "Elasticsearch cluster", Category: "search", Desc: "Cluster transport port", Badge: types.BadgeOfficial},
	{Port: 8529, Name: "ArangoDB", Category: "database", Desc: "Multi-model DB WebUI", Badge: types.BadgeOfficial},
	{Port: 7474, Name: "Neo4j", Category: "database", Desc: "Graph DB browser", Badge: types.BadgeOfficial},
	{Port: 8086, Name: "InfluxDB", Category: "database", Desc: "Time-series DB HTTP API", Badge: types.BadgeOfficial},

	// Message queues and streaming
	{Port: 5672, Name: "RabbitMQ AMQP", Category: "queue", Desc: "Message broker AMQP protocol", Badge: types.BadgeOfficial},
	{Port: 15672, Name: "RabbitMQ UI", Category: "queue", Desc: "RabbitMQ management dashboard", Badge: types.BadgeOfficial},
	{Port: 9092, Name: "Kafka", Category: "queue", Desc: "Event streaming broker", Badge: types.BadgeOfficial},
	{Port: 2181, Name: "Zookeeper", Category: "queue", Desc: "Kafka cluster coordination", Badge: types.BadgeOfficial},
	{Port: 4222, Name: "NATS", Category: "queue", Desc: "Lightweight messaging server", Badge: types.BadgeOfficial},
	{Port: 8222, Name: "NATS monitor", Category: "queue", Desc: "NATS HTTP monitoring endpoint", Badge: types.BadgeOfficial},

	// Storage and object stores
	{Port: 9001, Name: "MinIO Console", Category: "storage", Desc: "MinIO WebUI dashboard", Badge: types.BadgeConvention},
	{Port: 2049, Name: "NFS", Category: "storage", Desc: "Network file system", Badge: types.BadgeOfficial},
	{Port: 445, Name: "SMB/Samba", Category: "storage", Desc: "Windows file sharing protocol", Badge: types.BadgeOfficial},
	{Port: 22, Name: "SFTP/SSH", Category: "storage", Desc: "Secure file and shell access", Badge: types.BadgeOfficial},

	// Auth and identity
	{Port: 7000, Name: "Auth service", Category: "auth", Desc: "Custom JWT or session API", Badge: types.BadgeCustom},
	{Port: 7001, Name: "OAuth server", Category: "auth", Desc: "OAuth2 or OIDC provider", Badge: types.BadgeCustom},
	{Port: 389, Name: "LDAP", Category: "auth", Desc: "Directory services", Badge: types.BadgeOfficial},
	{Port: 636, Name: "LDAPS", Category: "auth", Desc: "LDAP over TLS", Badge: types.BadgeOfficial},

	// Infra and monitoring
	{Port: 2376, Name: "Docker TLS", Category: "infra", Desc: "Docker daemon TLS port", Badge: types.BadgeOfficial},
	{Port: 2377, Name: "Docker Swarm", Category: "infra", Desc: "Swarm cluster management", Badge: types.BadgeOfficial},
	{Port: 8125, Name: "StatsD", Category: "monitoring", Desc: "Metrics UDP intake", Badge: types.BadgeConvention},
	{Port: 9411, Name: "Zipkin", Category: "monitoring", Desc: "Distributed tracing UI", Badge: types.BadgeConvention},
	{Port: 16686, Name: "Jaeger UI", Category: "monitoring", Desc: "Distributed tracing dashboard", Badge: types.BadgeOfficial},
	{Port: 4317, Name: "OTEL gRPC", Category: "monitoring", Desc: "OpenTelemetry gRPC collector", Badge: types.BadgeOfficial},
	{Port: 4318, Name: "OTEL HTTP", Category: "monitoring", Desc: "OpenTelemetry HTTP collector", Badge: types.BadgeOfficial},

	// Homelab and self-hosted
	{Port: 8096, Name: "Jellyfin", Category: "media", Desc: "Media server HTTP interface", Badge: types.BadgeConvention},
	{Port: 8920, Name: "Jellyfin TLS", Category: "media", Desc: "Media server HTTPS interface", Badge: types.BadgeConvention},
	{Port: 8123, Name: "Home Assistant", Category: "homelab", Desc: "Home automation platform", Badge: types.BadgeConvention},
	{Port: 1880, Name: "Node-RED", Category: "homelab", Desc: "Flow-based automation tool", Badge: types.BadgeConvention},
	{Port: 7878, Name: "Radarr", Category: "media", Desc: "Movie manager WebUI", Badge: types.BadgeConvention},
	{Port: 8989, Name: "Sonarr", Category: "media", Desc: "TV show manager WebUI", Badge: types.BadgeConvention},
	{Port: 9696, Name: "Prowlarr", Category: "media", Desc: "Indexer manager WebUI", Badge: types.BadgeConvention},
	{Port: 8384, Name: "Syncthing", Category: "storage", Desc: "File sync WebUI", Badge: types.BadgeConvention},
	{Port: 51820, Name: "WireGuard", Category: "vpn", Desc: "WireGuard VPN UDP port", Badge: types.BadgeOfficial},
}
