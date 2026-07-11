package cmd

import (
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
)

func sync(c *cli.Context) error {
	repo := &config.Repo{
		ConfigDirectory: c.String("config-dir"),
	}
	client, err := createClient(c)
	if err != nil {
		return err
	}

	d, err := client.ReadConfig(c.Context)
	if err != nil {
		return err
	}

	err = repo.WriteConfig(d)
	if err != nil {
		return err
	}
	return nil
}
