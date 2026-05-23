package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"sigs.k8s.io/yaml"
)

// FromYAML
func FromYAML(config []byte, prefix string) ([]string, error) {
	j, err := yaml.YAMLToJSON(config)
	if err != nil {
		return nil, fmt.Errorf("error converting YAML to JSON: %w", err)
	}
	return FromJSON(j, prefix)
}

// FromJSON
func FromJSON(config []byte, prefix string) ([]string, error) {
	var nestedMap map[string]any
	err := json.Unmarshal(config, &nestedMap)
	if err != nil {
		return nil, err
	}
	return FromConfigMap(nestedMap, prefix)
}

func FromConfigMap(config map[string]interface{}, prefix string) ([]string, error) {
	var mm mapper
	err := mm.mapObj(prefix, config)
	return mm.cmds, err
}

type mapper struct {
	cmds []string
}

func (m *mapper) mapObj(cmd string, nm any) error {
	switch nm := nm.(type) {
	case map[string]any:
		for k, v := range nm {
			if err := m.mapKV(cmd, k, v); err != nil {
				return err
			}
		}
	case []any:
		return m.mapSlice(cmd+" ", nm)
	default:
		return errors.New("invalid input, must be a map or slice of interface")
	}
	return nil
}

func (m *mapper) mapKV(cmd string, k string, v any) error {
	cmd = join(cmd, k)
	switch vt := v.(type) {
	case map[string]any:
		if len(vt) == 0 {
			return m.mapValue(cmd, "")
		} else {
			return m.mapValue(cmd, v)
		}

	case []any:
		return m.mapSlice(cmd, vt)
	default:
		return m.mapValue(cmd, v)
	}
}

func (m *mapper) mapSlice(cmd string, vt []any) error {
	for _, val := range vt {
		if err := m.mapValue(cmd, val); err != nil {
			return err
		}
	}
	return nil
}

func (m *mapper) mapValue(cmd string, v any) error {
	switch v := v.(type) {
	case map[string]any, []any:
		return m.mapObj(cmd, v)

	case string:
		if strings.ContainsFunc(v, unicode.IsSpace) {
			var b strings.Builder
			b.WriteRune('"')
			b.WriteString(v)
			b.WriteRune('"')
			v = b.String()
		}
		m.cmds = append(m.cmds, join(cmd, v))

	default:
		m.cmds = append(m.cmds, cmd+" "+v.(string))
	}
	return nil
}

func join(a ...string) string {
	aa := slices.DeleteFunc(a, func(s string) bool {
		return s == ""
	})
	return strings.Join(aa, " ")
}
