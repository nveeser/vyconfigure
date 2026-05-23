package cmd

import (
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/pkg/commands"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
)

func plan(c *cli.Context) error {
	repo := &config.Repo{
		ConfigDirectory: c.String("config-dir"),
	}
	// get remote config as cmds
	client, err := createClient(c)
	if err != nil {
		return err
	}
	rc, err := client.ReadConfig(c.Context)
	if err != nil {
		return err
	}
	lc, err := repo.ReadConfig()
	if err != nil {
		return err
	}
	changes, err := commands.DiffConfigs(rc, lc)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		println("No changes to apply.")
		return nil
	}
	for _, change := range changes {
		switch change.Type {
		case commands.Added:
			color.Green("+ set " + change.Command)
		case commands.Deleted:
			color.Red("- delete " + change.Command)
		}
	}
	return nil
}
