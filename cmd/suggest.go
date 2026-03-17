package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/portmap/internal/registry"
	"github.com/yourusername/portmap/internal/render"
	"github.com/yourusername/portmap/internal/resolver"
	"github.com/yourusername/portmap/internal/types"
)

var suggestStack string
var suggestFix int
var suggestExport string

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Suggest a conventional port layout",
	Long: `Display a recommended port layout for a named service stack. Annotations
show which ports are already in use on your system and flag any conflicts.`,
	RunE: runSuggest,
}

func init() {
	suggestCmd.Flags().StringVar(&suggestStack, "stack", "fullstack", "Stack name: fullstack, backend, frontend, database, monitoring, homelab")
	suggestCmd.Flags().IntVar(&suggestFix, "fix", 0, "Suggest an alternative for a specific port")
	suggestCmd.Flags().StringVar(&suggestExport, "export", "", "Export format: compose, env, json")
}

var stacks = map[string][]string{
	"fullstack":  {"frontend", "api", "auth", "database", "cache", "storage", "monitoring"},
	"backend":    {"api", "auth", "database", "cache", "queue"},
	"frontend":   {"frontend"},
	"database":   {"database", "cache", "search"},
	"monitoring": {"monitoring"},
	"homelab":    {"homelab", "media", "vpn", "storage"},
}

func runSuggest(cmd *cobra.Command, args []string) error {
	if suggestFix > 0 {
		return runFix(suggestFix)
	}

	cats, ok := stacks[suggestStack]
	if !ok {
		return fmt.Errorf("suggest: unknown stack %q. Valid stacks: fullstack, backend, frontend, database, monitoring, homelab", suggestStack)
	}

	current, err := resolver.Run(resolver.Options{})
	if err != nil {
		return fmt.Errorf("suggest: scan: %w", err)
	}

	active := map[int]*types.PortEntry{}
	for _, e := range current.Entries {
		active[e.Port] = e
	}

	var layout []types.RegistryEntry
	seen := map[int]bool{}
	for _, cat := range cats {
		for _, e := range registry.ByCategory(cat) {
			if !seen[e.Port] {
				layout = append(layout, e)
				seen[e.Port] = true
			}
		}
	}

	switch suggestExport {
	case "compose":
		render.PrintSuggestCompose(layout, active)
	case "env":
		render.PrintSuggestEnv(layout, active)
	case "json":
		render.PrintSuggestJSON(layout, active)
	default:
		render.PrintSuggestTable(suggestStack, layout, active, current.Conflicts)
	}

	return nil
}

func runFix(port int) error {
	current, err := resolver.Run(resolver.Options{})
	if err != nil {
		return err
	}

	var conflictPair *types.ConflictPair
	for _, pair := range current.Conflicts {
		if pair.Port == port {
			conflictPair = pair
			break
		}
	}

	if conflictPair == nil {
		render.PrintInfoLine(fmt.Sprintf("port %d has no detected conflict", port))
		return nil
	}

	alt := findNextFree(port, current.Entries)
	render.PrintFixSuggestion(conflictPair, alt)
	return nil
}

func findNextFree(port int, entries []*types.PortEntry) int {
	bound := map[int]bool{}
	for _, e := range entries {
		bound[e.Port] = true
	}
	for _, delta := range []int{10, 1, 100, -10} {
		candidate := port + delta
		if candidate > 0 && candidate <= 65535 && !bound[candidate] {
			return candidate
		}
	}
	return 0
}
