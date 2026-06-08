package cmd

import (
	"cmp"
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/pkg/commands"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
	"slices"
)

func show(c *cli.Context) error {
	repo := &config.Repo{
		ConfigDirectory: c.String("config-dir"),
	}
	config, err := repo.ReadConfig()
	if err != nil {
		return err
	}
	cmds, err := commands.FromConfigMap(config, "")
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
