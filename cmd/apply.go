package cmd

import (
	"github.com/nveeser/vyconfigure/pkg/commands"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
)

func apply(c *cli.Context) error {
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
	var adds, dels []string
	for _, change := range changes {
		switch change.Type {
		case commands.Added:
			adds = append(adds, change.Command)
		case commands.Deleted:
			dels = append(dels, change.Command)
		}
	}
	return client.WriteCmds(c.Context, adds, dels)
}
