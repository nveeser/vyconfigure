package section

import (
	"testing"

	diffcmp "github.com/google/go-cmp/cmp"
)

func TestSplitMap(t *testing.T) {
	tests := []struct {
		name       string
		path       []string
		input      map[string]any
		wantRemain map[string]any
		wantSub    any
	}{
		{
			name: "Split single level key",
			path: []string{"foo"},
			input: map[string]any{
				"foo": "bar",
				"abc": 123,
			},
			wantRemain: map[string]any{
				"abc": 123,
			},
			wantSub: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "Split nested key",
			path: []string{"foo", "bar"},
			input: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
					"qux": 456,
				},
				"abc": 123,
			},
			wantRemain: map[string]any{
				"foo": map[string]any{
					"qux": 456,
				},
				"abc": 123,
			},
			wantSub: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
		},
		{
			name: "Split non-existent key",
			path: []string{"nonexistent"},
			input: map[string]any{
				"foo": "bar",
			},
			wantRemain: map[string]any{
				"foo": "bar",
			},
			wantSub: nil,
		},
		{
			name: "Split nested non-existent key",
			path: []string{"foo", "nonexistent"},
			input: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			wantRemain: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			wantSub: map[string]any{
				"foo": nil,
			},
		},
		{
			name: "Split path through non-map",
			path: []string{"foo", "bar"},
			input: map[string]any{
				"foo": "not-a-map",
			},
			wantRemain: map[string]any{
				"foo": "not-a-map",
			},
			wantSub: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := cloneMap(tt.input)
			entry := Section{steps: tt.path}
			gotRemain, gotSub := SplitMap(inputCopy, entry.steps, true)
			if diff := diffcmp.Diff(tt.wantRemain, gotRemain); diff != "" {
				t.Errorf("Extract() gotRemain mismatch (-want +got):\n%s", diff)
			}
			if diff := diffcmp.Diff(tt.wantSub, gotSub); diff != "" {
				t.Errorf("Extract() gotSub mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMergeMap(t *testing.T) {
	tests := []struct {
		name     string
		dest     map[string]any
		src      map[string]any
		expected map[string]any
	}{
		{
			name:     "both nil",
			dest:     nil,
			src:      nil,
			expected: map[string]any{},
		},
		{
			name: "dest nil",
			dest: nil,
			src: map[string]any{
				"foo": "bar",
			},
			expected: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "disjoint maps",
			dest: map[string]any{
				"foo": "bar",
			},
			src: map[string]any{
				"baz": 123,
			},
			expected: map[string]any{
				"foo": "bar",
				"baz": 123,
			},
		},
		{
			name: "overwrite simple value",
			dest: map[string]any{
				"foo": "bar",
			},
			src: map[string]any{
				"foo": "baz",
			},
			expected: map[string]any{
				"foo": "baz",
			},
		},
		{
			name: "deep Merge nested maps",
			dest: map[string]any{
				"nested": map[string]any{
					"a": 1,
					"b": 2,
				},
			},
			src: map[string]any{
				"nested": map[string]any{
					"b": 20,
					"c": 3,
				},
			},
			expected: map[string]any{
				"nested": map[string]any{
					"a": 1,
					"b": 20,
					"c": 3,
				},
			},
		},
		{
			name: "overwrite map with simple value",
			dest: map[string]any{
				"foo": map[string]any{
					"bar": 1,
				},
			},
			src: map[string]any{
				"foo": "simple",
			},
			expected: map[string]any{
				"foo": "simple",
			},
		},
		{
			name: "overwrite simple value with map",
			dest: map[string]any{
				"foo": "simple",
			},
			src: map[string]any{
				"foo": map[string]any{
					"bar": 1,
				},
			},
			expected: map[string]any{
				"foo": map[string]any{
					"bar": 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeMap(tt.dest, tt.src)
			if diff := diffcmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("mergeMap() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func cloneMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			res[k] = cloneMap(vm)
		} else {
			res[k] = v
		}
	}
	return res
}
