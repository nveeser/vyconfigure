package cmd

import (
	"github.com/nveeser/vyconfigure/pkg/api"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/nveeser/vyconfigure/pkg/options"
	"github.com/urfave/cli/v2"
)

func sync(c *cli.Context) error {
	o := options.GetOptions(c)
	repo := &config.Repo{o.ConfigDirectory}

	client, err := api.CreateClient(o)
	if err != nil {
		return err
	}

	d, err := client.Retrieve(c.Context)
	if err != nil {
		return err
	}
	err = repo.Write(d)
	if err != nil {
		return err
	}

	return nil
}
