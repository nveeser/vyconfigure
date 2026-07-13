package cmd

import (
	"cmp"
	"fmt"
	"github.com/fatih/color"
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/spf13/cobra"
	"slices"
)

func newApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "apply",
		Aliases: []string{"a", "push"},
		Short:   "Applies the current configuration.",
		RunE:    apply,
	}
}

func apply(cmd *cobra.Command, _ []string) error {
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
	it, err := commands.DiffConfigs(remote, local)
	if err != nil {
		return err
	}

	var reqs []vyos.ConfigRequest
	for change, entry := range it {
		switch change {
		case commands.Added:
			color.Green("  set " + entry.String())
			reqs = append(reqs, &vyos.SetRequest{Path: entry.Path, Value: entry.Value})
		case commands.Deleted:
			color.Red("  delete " + entry.String())
			reqs = append(reqs, &vyos.DeleteRequest{Path: entry.Path, Value: entry.Value})
		}
	}

	if len(reqs) == 0 {
		color.Blue("No changes to apply.")
		return nil
	}
	slices.SortFunc(reqs, compareConfigRequest)
	if err := client.ConfigMode().Configure(cmd.Context(), reqs...); err != nil {
		return fmt.Errorf("error calling Configure: %w", err)
	}
	return nil
}

func compareConfigRequest(a, b vyos.ConfigRequest) int {
	return cmp.Compare(rankConfigRequest(a), rankConfigRequest(b))
}

func rankConfigRequest(r vyos.ConfigRequest) int {
	switch r.(type) {
	case *vyos.DeleteRequest:
		return -1
	case *vyos.SetRequest:
		return 1
	default:
		return 0
	}
}
