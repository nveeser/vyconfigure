package cmd

import (
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/pkg/config"
	"github.com/urfave/cli/v2"
)

type action struct {
}

func createClient(c *cli.Context) (*config.Client, error) {
	opts := []vyos.Option{
		vyos.Token(c.String("api-key")),
	}
	if c.Bool("insecure") {
		opts = append(opts, vyos.Insecure())
	}
	if c.Bool("debug") {
		opts = append(opts, vyos.DebugLogging())
	}
	client, err := vyos.NewClient(c.String("host"), opts...)
	if err != nil {
		return nil, err
	}
	cc := &config.Client{
		Client: client,
	}
	return cc, nil
}
