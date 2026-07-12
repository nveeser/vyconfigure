package config

import (
	"context"
	"github.com/nveeser/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/section"
)

type Client struct {
	*vyos.Client
	SectionMapper *section.Mapper
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
	data, _ = section.SplitMap(data, []string{"system", "login"}, false)
	return data, nil
}

func (c *Client) ReadSections(ctx context.Context) ([]*section.Section, error) {
	data, err := c.ReadConfig(ctx)
	if err != nil {
		return nil, err
	}
	return c.SectionMapper.Split(data), nil
}
