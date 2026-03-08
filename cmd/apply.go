package cmd

import (
	"cmp"
	"slices"

	"github.com/nveeser/go-vyos/vyos"
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

	var reqs []vyos.ConfigRequest
	if len(changelog) > 0 {
		for _, change := range changelog {
			if change.Type == "create" {
				reqs = append(reqs, &vyos.SetRequest{change.To.(string)})
			}
			if change.Type == "delete" {
				reqs = append(reqs, &vyos.DeleteRequest{change.From.(string)})
			}
		}
	} else {
		println("No changes to apply.")
		return nil
	}

	slices.SortFunc(reqs, compareConfigRequest)
	return client.ConfigMode().Configure(c.Context, reqs...)
}

func compareConfigRequest(a, b vyos.ConfigRequest) int {
	return cmp.Compare(rankConfigRequest(a), rankConfigRequest(b))
}

func rankConfigRequest(r vyos.ConfigRequest) int {
	switch r.(type) {
	case *vyos.DeleteRequest:
		return -1
	case *vyos.SetRequest:
		return 1
	default:
		return 0
	}
}
