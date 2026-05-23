package cmd

import (
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/pkg/commands"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
	"sort"
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
	sort.Strings(cmds)
	for _, cmd := range cmds {
		color.Green(" set " + cmd)
	}

	return nil
}
