package cmd

import (
	"github.com/nveeser/go-vyos/vyos"
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

	it, err := commands.DiffConfigs(rc, lc)
	if err != nil {
		return err
	}

	var reqs []vyos.ConfigRequest
	for change, entry := range it {
		switch change {
		case commands.Added:
			reqs = append(reqs, &vyos.SetRequest{Path: entry.Path, Value: entry.Value})
		case commands.Deleted:
			reqs = append(reqs, &vyos.DeleteRequest{Path: entry.Path, Value: entry.Value})
		}
	}
	if len(reqs) == 0 {
		println("No changes to apply.")
		return nil
	}
	return client.WriteCmds(c.Context, reqs)
}
