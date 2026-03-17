package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/aymantoumi/portlens/internal/render"
	"github.com/aymantoumi/portlens/internal/resolver"
)

var freeCount int
var freeFormat string

var freeCmd = &cobra.Command{
	Use:   "free [MIN-MAX]",
	Short: "Find available ports in a range",
	Long: `Find ports that are not bound by any process or container and are not
reserved by portlens convention. Optionally provide a range like 4000-5000.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFree,
}

func init() {
	freeCmd.Flags().IntVar(&freeCount, "count", 5, "Number of free ports to return")
	freeCmd.Flags().StringVar(&freeFormat, "format", "table", "Output format: table, json, list")
}

func runFree(cmd *cobra.Command, args []string) error {
	min, max := 3000, 9999

	if len(args) == 1 {
		parts := strings.SplitN(args[0], "-", 2)
		if len(parts) != 2 {
			return fmt.Errorf("free: range must be in format MIN-MAX, got %q", args[0])
		}
		var err error
		min, err = strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("free: invalid min port %q: %w", parts[0], err)
		}
		max, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("free: invalid max port %q: %w", parts[1], err)
		}
		if min >= max {
			return fmt.Errorf("free: min (%d) must be less than max (%d)", min, max)
		}
		if min < 1 || max > 65535 {
			return fmt.Errorf("free: port range must be within 1-65535")
		}
	}

	result, err := resolver.FreeInRange(min, max, freeCount)
	if err != nil {
		return fmt.Errorf("free: %w", err)
	}

	switch freeFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	case "list":
		for _, p := range result.Available {
			fmt.Println(p)
		}
		return nil
	default:
		render.PrintFreeResult(result)
		return nil
	}
}
