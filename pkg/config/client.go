package config

import (
	"context"
	"github.com/nveeser/go-vyos/vyos"
)

type Client struct {
	*vyos.Client
}

func (c *Client) ReadConfig(ctx context.Context) (map[string]any, error) {
	data, err := c.ConfigMode().Show(ctx, "")
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	// extract the system.login path from the full config and drop it.
	data, _ = splitMap(data, []string{"system", "login"}, false)
	return data, nil
}
