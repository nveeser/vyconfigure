package cmd

import (
	"github.com/urfave/cli/v2"
)

func sync(c *cli.Context) error {
	repo, client, err := newRepoAndClient(c)
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
