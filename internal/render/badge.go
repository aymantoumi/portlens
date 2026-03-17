package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/aymantoumi/portlens/internal/types"
)

func BadgeStyle(b types.BadgeKind) string {
	switch b {
	case types.BadgeOfficial:
		return pill("official", lipgloss.Color("17"), colorBlue)
	case types.BadgeConvention:
		return pill("convention", lipgloss.Color("22"), colorGreen)
	case types.BadgeCustom:
		return pill("custom", lipgloss.Color("58"), colorAmber)
	default:
		return pill("unknown", lipgloss.Color("235"), colorMuted)
	}
}

func SourceTag(e *types.PortEntry) string {
	switch e.Kind {
	case types.SourceDockerCompose:
		return pill("compose", lipgloss.Color("17"), colorBlue)
	case types.SourceDocker:
		return pill("docker", lipgloss.Color("54"), colorPurple)
	case types.SourceProcess:
		return pill("process", lipgloss.Color("22"), colorGreen)
	case types.SourceSystemd:
		return pill("systemd", lipgloss.Color("58"), colorAmber)
	default:
		return pill("unknown", lipgloss.Color("52"), colorRed)
	}
}

func CategoryColor(category string) lipgloss.Style {
	switch category {
	case "database", "cache", "search":
		return stylePurple
	case "api", "auth":
		return styleBlue
	case "frontend", "web":
		return styleGreen
	case "storage", "infra", "homelab":
		return styleOrange
	case "monitoring":
		return styleAmber
	default:
		return stylePrimary
	}
}

func pill(text string, bg, fg lipgloss.Color) string {
	return lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Padding(0, 1).
		Render(text)
}
