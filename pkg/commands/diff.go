package commands

import r3diff "github.com/r3labs/diff/v3"

type ChangeType string

const (
	Added   ChangeType = "ADDED"
	Deleted ChangeType = "DELETED"
)

type Change struct {
	Type    ChangeType
	Command string
}

func DiffConfigs(from, to map[string]any) ([]Change, error) {
	var m1 mapper
	err := m1.mapObj("", from)
	if err != nil {
		return nil, err
	}
	var m2 mapper
	err = m2.mapObj("", to)
	if err != nil {
		return nil, err
	}

	// get diff
	changelog, err := r3diff.Diff(m1.cmds, m2.cmds)
	if err != nil {
		return nil, err
	}
	var diff []Change
	for _, change := range changelog {
		if change.Type == "create" {
			diff = append(diff, Change{Added, change.To.(string)})
		}
		if change.Type == "delete" {
			diff = append(diff, Change{Deleted, change.From.(string)})
		}
	}
	return diff, nil
}
