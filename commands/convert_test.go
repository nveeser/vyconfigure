package commands

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/nveeser/vyconfigure/commands/testdata"
	"testing"
)

func TestConvert(t *testing.T) {
	cases := []struct {
		name       string
		configFile string
		config     map[string]any
		want       []Entry
		wantErr    bool
	}{
		{
			name:       "basic",
			configFile: "basic.yaml",
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN default-action", Value: "drop"},
			},
		},
		{
			name:       "empty_object",
			configFile: "empty_object.yaml",
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN rule 30", Value: ""},
			},
		},
		{
			name:       "quote_string",
			configFile: "quote_string.yaml",
			want: []Entry{
				{Path: "test firewall ipv6-name WAN_IN rule 10 description", Value: "Hello this is a rule"},
			},
		},
		{
			name:       "slice",
			configFile: "slice.yaml",

			want: []Entry{
				{Path: "test service mdns repeater interface", Value: "eth1.10"},
				{Path: "test service mdns repeater interface", Value: "eth2.20"},
				{Path: "test service mdns repeater interface", Value: "eth1.50"},
			},
		},
		{
			name:       "non_string",
			configFile: "non_string.yaml",
			wantErr:    true,
		},
		{
			name:       "large",
			configFile: "large.yaml",
			want: []Entry{
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
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("FromConfigMap", func(t *testing.T) {
				config := tc.config
				if config == nil {
					config = testdata.ReadYAML(t, tc.configFile)
				}

				got, err := FromConfigMap(config, "test")

				switch {
				case err != nil && !tc.wantErr:
					t.Errorf("FromConfigMap() error = %v, wantErr %v", err, tc.wantErr)
				case err == nil && tc.wantErr:
					t.Errorf("FromConfigMap() error = %v, wantErr %v", err, tc.wantErr)
				}

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
