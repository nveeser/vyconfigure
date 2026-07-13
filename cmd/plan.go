package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/spf13/cobra"
)

func newPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "plan",
		Aliases: []string{"d", "diff"},
		Short:   "Shows the diff between the current directory and VyOS instance.",
		RunE:    plan,
	}
}

func plan(cmd *cobra.Command, _ []string) error {
	repo, client, err := newRepoAndClient(cmd.Context())
	if err != nil {
		return err
	}

	remote, err := client.ReadConfigTree(cmd.Context())
	if err != nil {
		return err
	}
	local, err := repo.ReadConfigTree()
	if err != nil {
		return err
	}
	changes, err := commands.DiffConfigs(remote, local)
	if err != nil {
		return err
	}
	var diffs bool
	for change, entry := range changes {
		diffs = true
		switch change {
		case commands.Added:
			color.Green("  set " + entry.String())
		case commands.Deleted:
			color.Red("  delete " + entry.String())
		}
	}
	if !diffs {
		fmt.Printf("No changes to apply.\n")
	}
	return nil
}
