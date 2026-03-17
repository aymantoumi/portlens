package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/portmap/internal/conflict"
	"github.com/yourusername/portmap/internal/render"
	"github.com/yourusername/portmap/internal/resolver"
	"github.com/yourusername/portmap/internal/types"
)

var scanJSON bool
var scanSkipDocker bool
var scanSkipProc bool
var scanCategory string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all active port bindings",
	Long: `Scan Docker containers, system processes, and systemd services for active
port bindings. Identifies each service, enriches with registry metadata,
and reports any conflicts.`,
	RunE: runScan,
}

func init() {
	scanCmd.Flags().BoolVar(&scanJSON, "json", false, "Output as JSON")
	scanCmd.Flags().BoolVar(&scanSkipDocker, "skip-docker", false, "Skip Docker socket scan")
	scanCmd.Flags().BoolVar(&scanSkipProc, "skip-proc", false, "Skip /proc scan")
	scanCmd.Flags().StringVar(&scanCategory, "category", "", "Filter by category (e.g. database, api, web)")
}

func runScan(cmd *cobra.Command, args []string) error {
	opts := resolver.Options{
		SkipDocker: scanSkipDocker,
		SkipProc:   scanSkipProc,
	}

	result, err := resolver.Run(opts)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	if scanCategory != "" {
		var filtered []*types.PortEntry
		var filteredConflicts []*types.ConflictPair
		filteredPorts := map[int]bool{}

		for _, e := range result.Entries {
			if e.Category == scanCategory {
				filtered = append(filtered, e)
				filteredPorts[e.Port] = true
			}
		}

		for _, c := range result.Conflicts {
			if filteredPorts[c.Port] {
				filteredConflicts = append(filteredConflicts, c)
			}
		}

		result.Entries = filtered
		result.Conflicts = filteredConflicts
	}

	if scanJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	render.PrintScanHeader(result.Duration, len(result.Entries), len(result.Conflicts))
	render.PrintScanTable(result.Entries)
	render.PrintScanSummary(result)

	for _, pair := range result.Conflicts {
		s := conflict.Suggest(pair)
		render.PrintConflictBlock(s)
	}

	if len(result.Conflicts) == 0 {
		render.PrintCommandHints()
	}

	if len(result.Conflicts) > 0 {
		os.Exit(1)
	}
	return nil
}
