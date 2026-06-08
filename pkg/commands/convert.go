package commands

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Entry struct {
	Path  string
	Value string
}

func (c Entry) String() string {
	return join(c.Path, c.Value)
}

func FromConfigMap(config map[string]interface{}, prefix string) ([]Entry, error) {
	var mm mapper
	err := mm.processObj(prefix, config)
	return mm.cmds, err
}

type mapper struct {
	cmds []Entry
}

func (m *mapper) processObj(cmd string, nm any) error {
	switch nm := nm.(type) {
	case map[string]any:
		for k, v := range nm {
			if err := m.processKV(cmd, k, v); err != nil {
				return err
			}
		}
	case []any:
		return m.processSlice(cmd+" ", nm)
	default:
		return fmt.Errorf("invalid input, must be a map or slice of interface: %s %T", cmd, nm)
	}
	return nil
}

func (m *mapper) processKV(cmd string, k string, v any) error {
	cmd = join(cmd, k)
	switch vt := v.(type) {
	case map[string]any:
		if len(vt) == 0 {
			return m.processValue(cmd, "")
		} else {
			return m.processValue(cmd, vt)
		}
	case []any:
		return m.processSlice(cmd, vt)
	default:
		return m.processValue(cmd, v)
	}
}

func (m *mapper) processSlice(cmd string, vt []any) error {
	for _, val := range vt {
		if err := m.processValue(cmd, val); err != nil {
			return err
		}
	}
	return nil
}

func (m *mapper) processValue(cmd string, v any) error {
	switch v := v.(type) {
	case map[string]any:
		return m.processObj(cmd, v)
	case []any:
		return m.processObj(cmd, v)
	case string:
		m.cmds = append(m.cmds, Entry{Path: cmd, Value: v})
	default:
		panic("invalid configuration type: " + reflect.TypeOf(v).String())
	}
	return nil
}

func join(a ...string) string {
	aa := slices.DeleteFunc(a, func(s string) bool {
		return s == ""
	})
	return strings.Join(aa, " ")
}
