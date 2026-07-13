package cmd

import (
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s", "pull"},
		Short:   "Syncs configuration to the current directory through the HTTP API.",
		RunE:    sync,
	}
}

func sync(cmd *cobra.Command, _ []string) error {
	repo, client, err := newRepoAndClient(cmd.Context())
	if err != nil {
		return err
	}
	d, err := client.ReadConfigTree(cmd.Context())
	if err != nil {
		return err
	}
	return repo.WriteConfigTree(d)
}
