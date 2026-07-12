package commands

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	cases := []struct {
		name   string
		config map[string]any
		want   []Entry
	}{
		{
			name: "basic",
			config: map[string]any{
				"firewall": map[string]any{
					"ipv6-name": map[string]any{
						"WAN_IN": map[string]any{
							"default-action": "drop",
						},
					},
				},
			},
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN default-action", Value: "drop"},
			},
		},
		{
			name: "empty object",
			config: map[string]any{
				"firewall": map[string]any{
					"ipv6-name": map[string]any{
						"WAN_IN": map[string]any{
							"rule": map[string]any{
								"30": map[string]any{},
							},
						},
					},
				},
			},
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN rule 30", Value: ""},
			},
		},
		{
			name: "quote string",
			config: map[string]any{
				"firewall": map[string]any{
					"ipv6-name": map[string]any{
						"WAN_IN": map[string]any{
							"rule": map[string]any{
								"10": map[string]any{
									"description": "Hello this is a rule",
								},
							},
						},
					},
				},
			},
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN rule 10 description", Value: "Hello this is a rule"},
			},
		},
		{
			name: "slice",
			config: map[string]any{
				"service": map[string]any{
					"mdns": map[string]any{
						"repeater": map[string]any{
							"interface": []any{
								"eth1.10",
								"eth2.20",
								"eth1.50",
							},
						},
					},
				},
			},
			want: []Entry{
				{Path: "test service mdns repeater interface", Value: "eth1.10"},
				{Path: "test service mdns repeater interface", Value: "eth2.20"},
				{Path: "test service mdns repeater interface", Value: "eth1.50"},
			},
		},
		{
			name:   "large",
			config: largeConfig,
			want:   largeCommands,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("FromConfigMap", func(t *testing.T) {
				got, err := FromConfigMap(tc.config, "test")

				assert.NoError(t, err)
				less := func(a, b Entry) bool {
					return a.String() < b.String()
				}
				if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(less)); diff != "" {
					t.Errorf("FromConfigMap() mismatch (-want +got):\n%s", diff)
				}
			})
		})
	}
}

var (
	largeConfig = map[string]any{
		"firewall": map[string]any{
			"ipv6-name": map[string]any{
				"WAN_IN": map[string]any{
					"default-action": "drop",
					"rule": map[string]any{
						"10": map[string]any{
							"description": "Hello this is a rule",
							"action":      "accept",
							"state": map[string]any{
								"established": "enable",
								"related":     "enable",
							},
						},
						"20": map[string]any{
							"action":   "accept",
							"protocol": "ipv6-icmp",
						},
						"30": map[string]any{},
					},
				},
			},
		},
		"service": map[string]any{
			"mdns": map[string]any{
				"repeater": map[string]any{
					"interface": []any{
						"eth1.10",
						"eth2.20",
						"eth1.50",
					},
				},
			},
		},
	}

	largeCommands = []Entry{
		{Path: "test firewall ipv6-name WAN_IN default-action", Value: "drop"},
		{Path: "test firewall ipv6-name WAN_IN rule 10 action", Value: "accept"},
		{Path: "test firewall ipv6-name WAN_IN rule 10 description", Value: "Hello this is a rule"},
		{Path: "test firewall ipv6-name WAN_IN rule 10 state established", Value: "enable"},
		{Path: "test firewall ipv6-name WAN_IN rule 10 state related", Value: "enable"},
		{Path: "test firewall ipv6-name WAN_IN rule 20 action", Value: "accept"},
		{Path: "test firewall ipv6-name WAN_IN rule 20 protocol", Value: "ipv6-icmp"},
		{Path: "test firewall ipv6-name WAN_IN rule 30", Value: ""},
		{Path: "test service mdns repeater interface", Value: "eth1.10"},
		{Path: "test service mdns repeater interface", Value: "eth2.20"},
		{Path: "test service mdns repeater interface", Value: "eth1.50"},
	}
)
