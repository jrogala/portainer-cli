// Package cmdutil provides shared helpers for CLI commands.
package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jrogala/portainer-cli/client"
	"github.com/jrogala/portainer-cli/config"
	"github.com/spf13/cobra"
)

// NewClient creates an authenticated Portainer client from config.
func NewClient() (*client.Client, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, err
	}
	if token.IsExpired() {
		return nil, fmt.Errorf("token expired, run 'portainer login'")
	}
	return client.New(config.URL(), token.JWT, config.EndpointID()), nil
}

// PrintJSON encodes v as indented JSON to stdout.
func PrintJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// IsJSON returns true if the --json persistent flag is set on the command's root.
func IsJSON(cmd *cobra.Command) bool {
	v, _ := cmd.Root().PersistentFlags().GetBool("json")
	return v
}

// NewTabWriter creates a tabwriter for aligned table output.
func NewTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// Render outputs data as JSON if --json is set, otherwise calls tableFunc.
func Render(cmd *cobra.Command, data any, tableFunc func()) {
	if IsJSON(cmd) {
		_ = PrintJSON(data)
		return
	}
	tableFunc()
}

// ExitErr prints an error to stderr and exits.
func ExitErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
