package commands

import (
	"iter"
)

type ChangeType string

const (
	Added   ChangeType = "ADDED"
	Deleted ChangeType = "DELETED"
)

func DiffConfigs(from, to map[string]any) (iter.Seq2[ChangeType, Entry], error) {
	var fm mapper
	err := fm.processObj("", from)
	if err != nil {
		return nil, err
	}
	var tm mapper
	err = tm.processObj("", to)
	if err != nil {
		return nil, err
	}

	return func(yield func(ChangeType, Entry) bool) {
		toMap := index(tm.cmds)
		fromMap := index(fm.cmds)

		for _, e := range fm.cmds {
			if _, exists := toMap[e.String()]; !exists {
				if !yield(Deleted, e) {
					return
				}
			}
		}
		for _, e := range tm.cmds {
			if _, exists := fromMap[e.String()]; !exists {
				if !yield(Added, e) {
					return
				}
			}
		}
	}, nil
}

func index(f []Entry) map[string]Entry {
	m := make(map[string]Entry, len(f))
	for _, item := range f {
		m[item.String()] = item
	}
	return m
}
