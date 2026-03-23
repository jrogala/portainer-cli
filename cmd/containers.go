package cmd

import (
	"fmt"

	"github.com/jrogala/portainer-cli/internal/cmdutil"
	"github.com/jrogala/portainer-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd, startCmd, stopCmd, restartCmd, logsCmd, inspectCmd)

	psCmd.Flags().BoolP("all", "a", false, "Show all containers (including stopped)")
	logsCmd.Flags().Int("tail", 100, "Number of lines to show")
}

var psCmd = &cobra.Command{
	Use:     "ps",
	Aliases: []string{"list", "ls"},
	Short:   "List containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		all, _ := cmd.Flags().GetBool("all")
		entries, err := ops.ListContainers(c, all)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, entries, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME\tIMAGE\tSTATE\tSTATUS")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					e.ID, e.Name, e.Image, e.State, e.Status)
			}
			w.Flush()
		})
		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Start a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := ops.StartContainer(c, args[0]); err != nil {
			return err
		}
		fmt.Printf("Started %s\n", args[0])
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop <name>",
	Short: "Stop a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := ops.StopContainer(c, args[0]); err != nil {
			return err
		}
		fmt.Printf("Stopped %s\n", args[0])
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart <name>",
	Short: "Restart a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := ops.RestartContainer(c, args[0]); err != nil {
			return err
		}
		fmt.Printf("Restarted %s\n", args[0])
		return nil
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs <name>",
	Short: "Show container logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		tail, _ := cmd.Flags().GetInt("tail")
		logs, err := ops.GetLogs(c, args[0], tail)
		if err != nil {
			return err
		}
		fmt.Println(logs)
		return nil
	},
}

var inspectCmd = &cobra.Command{
	Use:   "inspect <name>",
	Short: "Show container details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		detail, err := ops.InspectContainer(c, args[0])
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, detail, func() {
			fmt.Printf("Name:     %s\n", detail.Name)
			fmt.Printf("ID:       %s\n", detail.ID)
			fmt.Printf("Image:    %s\n", detail.Image)
			fmt.Printf("State:    %s\n", detail.State)
			fmt.Printf("Running:  %v\n", detail.Running)
			fmt.Printf("Started:  %s\n", detail.StartedAt)
			fmt.Printf("Restart:  %s\n", detail.RestartPolicy)
			for netName, ip := range detail.Networks {
				fmt.Printf("Network:  %s (%s)\n", netName, ip)
			}
		})
		return nil
	},
}
