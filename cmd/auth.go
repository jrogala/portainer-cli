package cmd

import (
	"fmt"
	"time"

	"github.com/jrogala/portainer-cli/client"
	"github.com/jrogala/portainer-cli/config"
	"github.com/jrogala/portainer-cli/internal/cmdutil"
	"github.com/jrogala/portainer-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd, stacksCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Portainer",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		jwt, err := client.Authenticate(config.URL(), username, password)
		if err != nil {
			return err
		}

		token := &config.Token{
			JWT:       jwt,
			ExpiresAt: time.Now().Add(8 * time.Hour),
		}
		if err := config.SaveToken(token); err != nil {
			return err
		}
		fmt.Println("Authenticated successfully!")
		return nil
	},
}

var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "List stacks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		entries, err := ops.ListStacks(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, entries, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME\tSTATUS")
			for _, e := range entries {
				fmt.Fprintf(w, "%d\t%s\t%s\n", e.ID, e.Name, e.Status)
			}
			w.Flush()
		})
		return nil
	},
}
