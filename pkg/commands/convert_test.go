package commands

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	cases := []struct {
		name   string
		config map[string]any
		want   []string
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
			want: []string{
				"test firewall ipv6-name WAN_IN default-action drop",
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
			want: []string{
				"test firewall ipv6-name WAN_IN rule 30",
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
			want: []string{
				"test firewall ipv6-name WAN_IN rule 10 description 'Hello this is a rule'",
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
			want: []string{
				"test service mdns repeater interface eth1.10",
				"test service mdns repeater interface eth2.20",
				"test service mdns repeater interface eth1.50",
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
				if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(strings.Compare)); diff != "" {
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

	largeCommands = []string{
		"test firewall ipv6-name WAN_IN default-action drop",
		"test firewall ipv6-name WAN_IN rule 10 action accept",
		"test firewall ipv6-name WAN_IN rule 10 description 'Hello this is a rule'",
		"test firewall ipv6-name WAN_IN rule 10 state established enable",
		"test firewall ipv6-name WAN_IN rule 10 state related enable",
		"test firewall ipv6-name WAN_IN rule 20 action accept",
		"test firewall ipv6-name WAN_IN rule 20 protocol ipv6-icmp",
		"test firewall ipv6-name WAN_IN rule 30",
		"test service mdns repeater interface eth1.10",
		"test service mdns repeater interface eth2.20",
		"test service mdns repeater interface eth1.50",
	}
)
