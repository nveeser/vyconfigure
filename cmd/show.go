package cmd

import (
	"cmp"
	"slices"

	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Aliases: []string{"sh"},
		Short:   "Shows the local instance.",
		RunE:    show,
	}
}

func show(_ *cobra.Command, _ []string) error {
	repo, err := newRepo()
	if err != nil {
		return err
	}
	cfg, err := repo.ReadConfigTree()
	if err != nil {
		return err
	}
	cmds, err := commands.FromConfigMap(cfg, "")
	if err != nil {
		return err
	}
	slices.SortFunc(cmds, func(a, b commands.Entry) int {
		return cmp.Or(cmp.Compare(a.Path, b.Path), cmp.Compare(a.Value, b.Value))
	})
	for _, cmd := range cmds {
		color.Green(" set " + cmd.String())
	}
	return nil
}
