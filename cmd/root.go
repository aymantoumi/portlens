package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version string

func Execute(v string) {
	version = v
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "portlens",
	Short: "Port visibility and conflict detection for Docker and system services",
	Long: `portlens scans your machine for active port bindings across Docker containers,
system processes, and systemd services. It identifies what each port belongs to,
detects conflicts, and suggests conventional port layouts.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(freeCmd)
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(watchCmd)
}
