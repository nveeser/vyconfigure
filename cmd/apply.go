package cmd

import (
	"github.com/nveeser/vyconfigure/pkg/api"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/nveeser/vyconfigure/pkg/convert"
	"github.com/nveeser/vyconfigure/pkg/options"
	r3diff "github.com/r3labs/diff/v3"
	"github.com/urfave/cli/v2"
)

func apply(c *cli.Context) error {
	o := options.GetOptions(c)
	repo := &config.Repo{o.ConfigDirectory}

	// get remote config as cmds
	client, err := api.CreateClient(o)
	if err != nil {
		return err
	}
	d, err := client.RetrieveJson(c.Context)
	if err != nil {
		return err
	}

	rc, _ := convert.JsonToCmds(d, "")

	// get local config as cmds
	lc, err := repo.ReadAsCmds()
	if err != nil {
		return err
	}

	// get diff
	changelog, err := r3diff.Diff(rc, lc)
	if err != nil {
		return err
	}

	var toDelete []string
	var toCreate []string
	if len(changelog) > 0 {
		for _, change := range changelog {
			if change.Type == "create" {
				toCreate = append(toCreate, change.To.(string))
			}
			if change.Type == "delete" {
				toDelete = append(toDelete, change.From.(string))
			}
		}
	} else {
		println("No changes to apply.")
		return nil
	}

	dc := convert.CmdsToData(toDelete, "delete")
	cc := convert.CmdsToData(toCreate, "set")

	err = client.ConfigMode().Configure(c.Context, append(dc, cc...)...)
	if err != nil {
		return err
	}

	return nil
}
