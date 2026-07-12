package cmd

import (
	"cmp"
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/commands"
	"github.com/urfave/cli/v2"
	"slices"
)

func apply(c *cli.Context) error {
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
	it, err := commands.DiffConfigs(remote, local)
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
	slices.SortFunc(reqs, compareConfigRequest)
	err = client.ConfigMode().Configure(c.Context, reqs...)
	//var httpErr *vyos.HTTPError
	//if errors.As(err, &httpErr) {
	//	httpErr.
	//}
	if err != nil {
		return err
	}
	return nil
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
