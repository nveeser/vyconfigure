package cmd

import (
	"cmp"
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/urfave/cli/v2"
	"slices"
)

func show(c *cli.Context) error {
	repo, err := newRepo(c)
	if err != nil {
		return err
	}
	cfg, err := repo.ReadConfig()
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
