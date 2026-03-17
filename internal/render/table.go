package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/portmap/internal/types"
)

const (
	boxTopLeft  = "┌"
	boxTopRight = "┐"
	boxBotLeft  = "└"
	boxBotRight = "┘"
	boxMidLeft  = "├"
	boxMidRight = "┤"
	boxBar      = "│"
	boxHoriz    = "─"
	boxCross    = "┼"
)

func PrintScanHeader(duration string, count int, conflictCount int) {
	fmt.Printf(" %s  %s\n", styleGreen.Render("▣"), styleWhite.Bold(true).Render("portlens scan"))
	if conflictCount > 0 {
		fmt.Printf(" %s  %d ports  ·  %s conflict  ·  %s\n",
			styleGreen.Render("✓"), count, styleRed.Render(fmt.Sprintf("%d", conflictCount)), duration)
	} else {
		fmt.Printf(" %s  %d ports  ·  %s conflict  ·  %s\n",
			styleGreen.Render("✓"), count, styleMuted.Render("0"), duration)
	}
}

func PrintScanTable(entries []*types.PortEntry) {
	if len(entries) == 0 {
		return
	}

	cols := calcColumns()
	width := cols.Port + cols.Svc + cols.Cat + cols.Src + cols.Badge + cols.Info + 12

	border := fmt.Sprintf(" %s%s%s", boxTopLeft, strings.Repeat(boxHoriz, width), boxTopRight)
	fmt.Println(border)

	header := fmt.Sprintf(" %s %s %s %s %s %s %s",
		boxBar,
		styleDim.Width(cols.Port).Render("PORT"),
		styleDim.Width(cols.Svc).Render("SERVICE"),
		styleDim.Width(cols.Cat).Render("CAT"),
		styleDim.Width(cols.Src).Render("SRC"),
		styleDim.Width(cols.Badge).Render("BADGE"),
		styleDim.Width(cols.Info).Render("INFO"),
	)
	header = strings.TrimRight(header, " ")
	fmt.Println(header)

	midBorder := fmt.Sprintf(" %s%s%s", boxMidLeft, strings.Repeat(boxHoriz, width), boxMidRight)
	fmt.Println(midBorder)

	for i, e := range entries {
		printEntryRow(e, cols)
		if i < len(entries)-1 {
			rowSep := fmt.Sprintf(" %s%s%s", boxMidLeft, strings.Repeat(boxHoriz, width), boxMidRight)
			fmt.Println(rowSep)
		}
	}

	botBorder := fmt.Sprintf(" %s%s%s", boxBotLeft, strings.Repeat(boxHoriz, width), boxBotRight)
	fmt.Println(botBorder)
	fmt.Println()
}

func printEntryRow(e *types.PortEntry, cols TableColumns) {
	portStyle := CategoryColor(e.Category)
	if e.ConflictWith != nil {
		portStyle = styleRed
	}

	svcName := e.RegistryName
	if svcName == "" {
		svcName = e.ProcessName
	}
	if svcName == "" {
		svcName = "?"
	}

	infoStr := getInfoString(e)

	badge := BadgeStyle(e.Badge)
	if e.ConflictWith != nil {
		badge = pill("conflict", lipgloss.Color("52"), colorRed)
	}

	sourceTag := SourceTag(e)

	row := fmt.Sprintf(" %s %s %s %s %s %s %s",
		boxBar,
		portStyle.Width(cols.Port).Render(truncate(fmt.Sprintf("%d", e.Port), cols.Port)),
		stylePrimary.Width(cols.Svc).Render(truncate(svcName, cols.Svc)),
		styleMuted.Width(cols.Cat).Render(truncate(e.Category, cols.Cat)),
		sourceTag,
		badge,
		styleDim.Width(cols.Info).Render(truncate(infoStr, cols.Info)),
	)

	if e.ConflictWith != nil {
		row = styleConflictRow.Render(row)
	}

	fmt.Println(strings.TrimRight(row, " "))
}

func getInfoString(e *types.PortEntry) string {
	switch e.Kind {
	case types.SourceDockerCompose:
		parts := []string{}
		if e.ContainerName != "" {
			parts = append(parts, e.ContainerName)
		}
		if e.ImageName != "" {
			parts = append(parts, e.ImageName)
		}
		return strings.Join(parts, " · ")
	case types.SourceDocker:
		return e.ContainerName
	case types.SourceSystemd:
		if e.SystemdUnit != "" {
			return e.SystemdUnit
		}
		return fmt.Sprintf("pid:%d", e.PID)
	case types.SourceProcess:
		if e.ExePath != "" {
			parts := strings.Split(e.ExePath, "/")
			return parts[len(parts)-1]
		}
		return fmt.Sprintf("pid:%d", e.PID)
	case types.SourceUnknown:
		return "sudo lsof"
	}
	return ""
}

func printSubLine(e *types.PortEntry) {
	line := buildSubLine(e)
	if line == "" {
		return
	}
	fmt.Printf(" %s %s %s\n", boxBar, styleDim.Render("└"), styleMuted.Render(line))
}

func buildSubLine(e *types.PortEntry) string {
	switch e.Kind {
	case types.SourceDockerCompose:
		parts := []string{}
		if e.ContainerName != "" {
			parts = append(parts, stylePrimary.Render(e.ContainerName))
		}
		if e.ImageName != "" {
			parts = append(parts, styleMuted.Render(e.ImageName))
		}
		if e.ContainerID != "" {
			parts = append(parts, styleDim.Render(e.ContainerID[:12]))
		}
		if e.ComposeFile != "" {
			parts = append(parts, styleDim.Render(truncatePath(e.ComposeFile)))
		}
		return strings.Join(parts, " · ")

	case types.SourceDocker:
		return fmt.Sprintf("%s · %s · %s",
			stylePrimary.Render(e.ContainerName),
			styleMuted.Render(e.ImageName),
			styleDim.Render(e.ContainerID[:12]),
		)

	case types.SourceSystemd:
		parts := []string{}
		if e.SystemdUnit != "" {
			parts = append(parts, styleAmber.Render(e.SystemdUnit))
		}
		parts = append(parts, fmt.Sprintf("pid %s", stylePrimary.Render(fmt.Sprintf("%d", e.PID))))
		parts = append(parts, styleMuted.Render("active"))
		return strings.Join(parts, " · ")

	case types.SourceProcess:
		parts := []string{}
		if e.ProcessName != "" {
			parts = append(parts, stylePrimary.Render(e.ProcessName))
		}
		parts = append(parts, fmt.Sprintf("pid %s", stylePrimary.Render(fmt.Sprintf("%d", e.PID))))
		if e.Username != "" {
			parts = append(parts, fmt.Sprintf("user: %s", stylePrimary.Render(e.Username)))
		}
		if e.ExePath != "" {
			parts = append(parts, styleMuted.Render(truncatePath(e.ExePath)))
		}
		return strings.Join(parts, " · ")

	case types.SourceUnknown:
		return styleMuted.Render(fmt.Sprintf("sudo lsof -i :%d", e.Port))
	}
	return ""
}

func truncatePath(path string) string {
	if len(path) <= 40 {
		return path
	}
	parts := strings.Split(path, "/")
	if len(parts) <= 2 {
		return "..." + path[len(path)-37:]
	}
	return "..." + strings.Join(parts[len(parts)-2:], "/")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

func PrintScanSummary(result *types.ScanResult) {
	unknown := 0
	for _, e := range result.Entries {
		if e.Kind == types.SourceUnknown {
			unknown++
		}
	}
	clean := len(result.Entries) - len(result.Conflicts) - unknown

	activePill := strings.TrimSpace(lipgloss.NewStyle().Background(lipgloss.Color("22")).Foreground(colorGreen).Padding(0, 2).Render(fmt.Sprintf("%d active", len(result.Entries))))
	conflictPill := ""
	if len(result.Conflicts) > 0 {
		conflictPill = "  " + strings.TrimSpace(lipgloss.NewStyle().Background(lipgloss.Color("52")).Foreground(colorRed).Padding(0, 2).Render(fmt.Sprintf("%d conflict", len(result.Conflicts))))
	}
	unresolvedPill := ""
	if unknown > 0 {
		unresolvedPill = "  " + strings.TrimSpace(lipgloss.NewStyle().Background(lipgloss.Color("58")).Foreground(colorAmber).Padding(0, 2).Render(fmt.Sprintf("%d unresolved", unknown)))
	}
	cleanPill := "  " + strings.TrimSpace(lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(colorMuted).Padding(0, 2).Render(fmt.Sprintf("%d clean", clean)))

	summaryLine := fmt.Sprintf("%s%s%s%s", activePill, conflictPill, unresolvedPill, cleanPill)
	fmt.Println(strings.TrimRight(summaryLine, " "))
	fmt.Println()
}

func PrintCommandHints() {
	hints := []struct {
		cmd  string
		desc string
	}{
		{"portlens free", "find free ports"},
		{"portlens suggest", "get port layout"},
		{"portlens watch", "monitor changes"},
		{"portlens scan --category <cat>", "filter category"},
	}

	fmt.Println("Commands:")
	for _, h := range hints {
		fmt.Printf("  %s  %s  %s\n", styleGreen.Render("▸"), stylePrimary.Render(h.cmd), styleMuted.Render(h.desc))
	}
}
