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
		name string
		json string
		yaml string
		want []string
	}{

		{
			name: "basic",
			json: `{ "firewall": { "ipv6-name": { "WAN_IN": { "default-action": "drop" } } } }`,
			yaml: `
firewall:
  ipv6-name: 
    WAN_IN:
      default-action: drop`,
			want: []string{
				"test firewall ipv6-name WAN_IN default-action drop",
			},
		},
		{
			name: "empty object",
			json: `{ "firewall": { "ipv6-name": { "WAN_IN": { "rule": { "30": {} } } } } }`,
			yaml: `
firewall:
  ipv6-name:
    WAN_IN:
      rule:
        "30": {}`,
			want: []string{
				"test firewall ipv6-name WAN_IN rule 30",
			},
		},
		{
			name: "quote string",
			json: `{ "firewall": { "ipv6-name": { "WAN_IN": { "rule": { "10": { "description": "Hello this is a rule" } } } } } } `,
			yaml: `
firewall:
  ipv6-name:
    WAN_IN:
      rule:
        "10":
            description: Hello this is a rule
`,
			want: []string{
				"test firewall ipv6-name WAN_IN rule 10 description \"Hello this is a rule\"",
			},
		},
		{
			name: "slice",
			json: `{ "service": { "mdns": { "repeater": { "interface": [ "eth1.10", "eth2.20", "eth1.50" ] } } } }`,
			yaml: `
service:
    mdns:
        repeater:
            interface:
            - eth1.10
            - eth2.20
            - eth1.50
`,
			want: []string{
				"test service mdns repeater interface eth1.10",
				"test service mdns repeater interface eth2.20",
				"test service mdns repeater interface eth1.50",
			},
		},
		{
			name: "large",
			json: largeJSON,
			yaml: largeYAML,
			want: largeCommands,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("FromJSON", func(t *testing.T) {
				got, err := FromJSON([]byte(tc.json), "test")

				assert.NoError(t, err)
				if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(strings.Compare)); diff != "" {
					t.Errorf("FromJSON() mismatch (-want +got):\n%s", diff)
				}
			})
			t.Run("FromYAML", func(t *testing.T) {
				got, err := FromYAML([]byte(tc.yaml), "test")

				assert.NoError(t, err)
				if diff := cmp.Diff(tc.want, got, cmpopts.SortSlices(strings.Compare)); diff != "" {
					t.Errorf("FromYAML() mismatch (-want +got):\n%s", diff)
				}
			})
		})
	}
}

var (
	largeCommands = []string{
		"test firewall ipv6-name WAN_IN default-action drop",
		"test firewall ipv6-name WAN_IN rule 10 action accept",
		"test firewall ipv6-name WAN_IN rule 10 description \"Hello this is a rule\"",
		"test firewall ipv6-name WAN_IN rule 10 state established enable",
		"test firewall ipv6-name WAN_IN rule 10 state related enable",
		"test firewall ipv6-name WAN_IN rule 20 action accept",
		"test firewall ipv6-name WAN_IN rule 20 protocol ipv6-icmp",
		"test firewall ipv6-name WAN_IN rule 30",
		"test service mdns repeater interface eth1.10",
		"test service mdns repeater interface eth2.20",
		"test service mdns repeater interface eth1.50",
	}

	largeYAML = `
firewall:
 ipv6-name:
   WAN_IN:
     default-action: drop
     rule:
       "10":
           description: Hello this is a rule
           action: accept
           state:
             established: enable
             related: enable
       "20":
           action: accept
           protocol: ipv6-icmp
       "30": {}
service:
   mdns:
       repeater:
           interface:
           - eth1.10
           - eth2.20
           - eth1.50
`

	largeJSON = `
   {
       "firewall": {
           "ipv6-name": {
               "WAN_IN": {
                   "default-action": "drop",
                   "rule": {
                       "10": {
							"description": "Hello this is a rule",
                           "action": "accept",
                           "state": {
                               "established": "enable",
                               "related": "enable"
                           }
                       },
                       "20": {
                           "action": "accept",
                           "protocol": "ipv6-icmp"
                       },
                       "30": {}
                   }
               }
           }
       },
       "service": {
           "mdns": {
               "repeater": {
                   "interface": [
                       "eth1.10",
                       "eth2.20",
                       "eth1.50"
                   ]
               }
           }
       }
   }
   `
)
