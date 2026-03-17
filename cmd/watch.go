package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/aymantoumi/portlens/internal/conflict"
	"github.com/aymantoumi/portlens/internal/render"
	"github.com/aymantoumi/portlens/internal/resolver"
	"github.com/aymantoumi/portlens/internal/types"
)

var watchInterval int
var watchOnce bool

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Live monitor for port binding changes",
	Long: `Poll for port binding changes and print one log line per event.
New bindings are shown in green (+), released bindings in yellow (-),
and conflicts in red (!). Press Ctrl+C to exit.`,
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().IntVar(&watchInterval, "interval", 2, "Poll interval in seconds")
	watchCmd.Flags().BoolVar(&watchOnce, "once", false, "Poll once and exit")
}

func runWatch(cmd *cobra.Command, args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	render.PrintWatchHeader(watchInterval)

	prev := map[int]*types.PortEntry{}
	knownConflicts := map[int]bool{}

	poll := func() error {
		result, err := resolver.Run(resolver.Options{})
		if err != nil {
			return err
		}

		current := map[int]*types.PortEntry{}
		for _, e := range result.Entries {
			current[e.Port] = e
		}

		ts := time.Now().Format("15:04:05")

		for port, e := range current {
			if _, existed := prev[port]; !existed {
				render.PrintWatchEvent(ts, "+", port, e)
			}
		}

		for port, e := range prev {
			if _, stillExists := current[port]; !stillExists {
				render.PrintWatchEvent(ts, "-", port, e)
			}
		}

		for _, pair := range result.Conflicts {
			if !knownConflicts[pair.Port] {
				knownConflicts[pair.Port] = true
				render.PrintWatchEvent(ts, "!", pair.Port, pair.A)
				s := conflict.Suggest(pair)
				render.PrintConflictBlock(s)
			}
		}

		activePorts := map[int]bool{}
		for _, pair := range result.Conflicts {
			activePorts[pair.Port] = true
		}
		for p := range knownConflicts {
			if !activePorts[p] {
				delete(knownConflicts, p)
			}
		}

		prev = current
		return nil
	}

	if err := poll(); err != nil {
		return err
	}

	if watchOnce {
		return nil
	}

	ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return nil
		case <-ticker.C:
			if err := poll(); err != nil {
				render.PrintErrorLine(err.Error())
			}
		}
	}
}
