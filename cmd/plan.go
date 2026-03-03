package cmd

import (
	"github.com/fatih/color"
	"github.com/nveeser/vyconfigure/pkg/api"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/nveeser/vyconfigure/pkg/convert"
	"github.com/nveeser/vyconfigure/pkg/options"
	diff "github.com/r3labs/diff/v3"
	"github.com/urfave/cli/v2"
)

func plan(c *cli.Context) error {
	o := options.GetOptions(c)
	repo := config.Repo{o.ConfigDirectory}
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
	changelog, err := diff.Diff(rc, lc)
	if err != nil {
		return err
	}

	if len(changelog) > 0 {
		println("Changes to be applied:")
		for _, change := range changelog {
			if change.Type == "create" {
				color.Green("+ set " + change.To.(string))
			}
			if change.Type == "delete" {
				color.Red("- delete " + change.From.(string))
			}
		}
	} else {
		println("No changes to apply.")
	}

	return nil
}
