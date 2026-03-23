// Package cmd implements the portainer-cli commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/jrogala/portainer-cli/config"
	"github.com/jrogala/portainer-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portainer",
	Short: "CLI for managing Docker containers via Portainer",
	Long: `CLI for managing Docker containers via Portainer API.

Quick examples:
  portainer login -u user -p pass                 Authenticate
  portainer ps                                     List running containers
  portainer ps -a                                  List all containers (including stopped)
  portainer start <name>                           Start a container
  portainer stop <name>                            Stop a container
  portainer restart <name>                         Restart a container
  portainer logs <name>                            Show container logs
  portainer logs <name> --tail 50                  Show last 50 log lines
  portainer inspect <name>                         Show container details
  portainer stacks                                 List stacks`,
}

func Execute() {
	config.Init()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output raw JSON responses")
}

// Re-export for subcommands in this package.
var newClient = cmdutil.NewClient
