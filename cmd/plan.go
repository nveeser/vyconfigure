package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/urfave/cli/v2"
)

func plan(c *cli.Context) error {
	repo, client, err := newRepoAndClient(c)
	if err != nil {
		return err
	}
	remote, err := client.ReadConfig(c.Context)
	if err != nil {
		return err
	}
	local, err := repo.ReadConfig()
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
		return nil
	}
	return nil
}
