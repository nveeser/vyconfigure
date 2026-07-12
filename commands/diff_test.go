package commands

import (
	"cmp"
	"testing"

	diffcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type diffEntry struct {
	Type  ChangeType
	Entry Entry
}

func TestDiffConfigs(t *testing.T) {
	tests := []struct {
		name     string
		from     any
		to       any
		wantErr  bool
		wantDiff []diffEntry
	}{
		{
			name: "no changes",
			from: map[string]any{
				"firewall": map[string]any{
					"default-action": "drop",
				},
			},
			to: map[string]any{
				"firewall": map[string]any{
					"default-action": "drop",
				},
			},
			wantErr:  false,
			wantDiff: nil,
		},
		{
			name: "added config",
			from: map[string]any{},
			to: map[string]any{
				"firewall": map[string]any{
					"default-action": "drop",
				},
			},
			wantDiff: []diffEntry{
				{
					Type:  Added,
					Entry: Entry{Path: "firewall default-action", Value: "drop"},
				},
			},
		},
		{
			name: "deleted config",
			from: map[string]any{
				"firewall": map[string]any{
					"default-action": "drop",
				},
			},
			to: map[string]any{},
			wantDiff: []diffEntry{
				{
					Type:  Deleted,
					Entry: Entry{Path: "firewall default-action", Value: "drop"},
				},
			},
		},
		{
			name: "modified config",
			from: map[string]any{
				"firewall": map[string]any{
					"default-action": "drop",
				},
			},
			to: map[string]any{
				"firewall": map[string]any{
					"default-action": "accept",
				},
			},
			wantDiff: []diffEntry{
				{
					Type:  Deleted,
					Entry: Entry{Path: "firewall default-action", Value: "drop"},
				},
				{
					Type:  Added,
					Entry: Entry{Path: "firewall default-action", Value: "accept"},
				},
			},
		},
		{
			name:    "invalid from config",
			from:    123,
			to:      map[string]any{},
			wantErr: true,
		},
		{
			name:    "invalid to config",
			from:    map[string]any{},
			to:      123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, err := DiffConfigs(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DiffConfigs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			var gotDiff []diffEntry
			for ct, e := range seq {
				gotDiff = append(gotDiff, diffEntry{Type: ct, Entry: e})
			}

			less := func(a, b diffEntry) bool {
				return cmp.Or(cmp.Compare(a.Type, b.Type),
					cmp.Compare(a.Entry.String(), b.Entry.String())) < 0
			}
			if diff := diffcmp.Diff(tt.wantDiff, gotDiff, cmpopts.SortSlices(less)); diff != "" {
				t.Errorf("DiffConfigs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
