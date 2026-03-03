package api

import (
	"context"
	"encoding/json"
	"github.com/ganawaj/go-vyos/vyos"
	"github.com/nveeser/vyconfigure/pkg/options"
)

type Client struct {
	*vyos.Client
}

func CreateClient(o *options.Options) (*Client, error) {
	client, err := vyos.NewClient(o.Host, vyos.Token(o.ApiKey), vyos.Insecure(), vyos.DebugLogging())
	if err != nil {
		return nil, err
	}
	c := &Client{
		Client: client,
	}
	return c, nil
}

func (c *Client) RetrieveJson(ctx context.Context) ([]byte, error) {
	data, err := c.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (c *Client) Retrieve(ctx context.Context) (map[string]any, error) {
	data, err := c.ConfigMode().Show(ctx, "")
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	// delete login configuration under "system" due to complexities with encrypted passwords
	for key := range data {
		if key == "system" {
			users := data[key]
			u := users.(map[string]any)
			delete(u, "login")
		}
	}
	return data, nil
}
