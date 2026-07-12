package cmd

import (
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/config"
	"github.com/nveeser/vyconfigure/section"
	"github.com/urfave/cli/v2"
	"path"
	"strings"
)

const defaultMappingFile = "vysync.yaml"

func newRepo(c *cli.Context) (*config.Repo, error) {
	mapper, err := section.NewMapper(path.Join(c.String("config-dir"), defaultMappingFile))
	if err != nil {
		return nil, err
	}
	return &config.Repo{
		ConfigDirectory: c.String("config-dir"),
		SectionMapper:   mapper,
		Ignore:          []string{strings.TrimSuffix(defaultMappingFile, ".yaml")},
	}, nil
}

func newRepoAndClient(c *cli.Context) (*config.Repo, *config.Client, error) {
	repo, err := newRepo(c)
	if err != nil {
		return nil, nil, err
	}
	opts := []vyos.Option{
		vyos.Token(c.String("api-key")),
		vyos.Timeout(0),
	}
	if c.Bool("insecure") {
		opts = append(opts, vyos.Insecure())
	}
	if c.Bool("debug") {
		opts = append(opts, vyos.DebugLogging())
	}
	client, err := vyos.NewClient(c.String("host"), opts...)
	if err != nil {
		return nil, nil, err
	}
	cc := &config.Client{
		Client:        client,
		SectionMapper: repo.SectionMapper,
	}
	return repo, cc, nil
}
