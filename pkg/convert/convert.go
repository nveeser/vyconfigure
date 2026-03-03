package convert

import (
	"encoding/json"
	"errors"
	"github.com/ganawaj/go-vyos/vyos"
	"strings"

	"sigs.k8s.io/yaml"
)

// YamlToCmds
func YamlToCmds(config []byte, prefix string) ([]string, error) {
	j, _ := yaml.YAMLToJSON(config)
	var nestedMap map[string]interface{}
	err := json.Unmarshal(j, &nestedMap)
	if err != nil {
		return nil, err
	}

	var mm mapper
	return mm.toCmds(nestedMap, prefix)
}

// JsonToCmds
func JsonToCmds(config []byte, prefix string) ([]string, error) {
	var nestedMap map[string]interface{}
	err := json.Unmarshal(config, &nestedMap)
	if err != nil {
		return nil, err
	}
	var mm mapper
	return mm.toCmds(nestedMap, prefix)
}

// CmdsToData converts a list of commands into a format the HTTP API understands
func CmdsToData(cmds []string, op string) []vyos.ConfigRequest {
	var res []vyos.ConfigRequest
	for _, c := range cmds {
		switch op {
		case "delete":
			res = append(res, &vyos.DeleteRequest{c})
		case "add":
			res = append(res, &vyos.SetRequest{c})
		}
	}
	return res
}

type mapper struct {
	cmds []string
}

// toCmds maps an interface to an array of VyOS commands
func (m *mapper) toCmds(nm any, prefix string) ([]string, error) {
	err := m.mapObj(true, nm, prefix)
	return m.cmds, err
}

func (m *mapper) mapObj(top bool, nm any, prefix string) error {
	switch nm := nm.(type) {
	case map[string]any:
		for k, v := range nm {
			if err := m.mapKV(top, prefix, k, v); err != nil {
				return err
			}
		}
	case []any:
		for _, v := range nm {
			cmd := m.buildCmd(isArray, prefix, "")
			if err := m.assign(cmd, v, false, false); err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid input, must be a map or slice of interface")
	}
	return nil
}

func (m *mapper) mapKV(top bool, prefix string, k string, v any) error {
	cmd := m.buildCmd(top, prefix, k)

	// this is pretty ugly, basically when building the cmds we only care about the key if the value is {}
	r, _ := json.Marshal(v)
	res := string(r)

	switch {
	case res == "{}":
		if err := m.assign(cmd, k, isScalar, ignoreValue); err != nil {
			return err
		}
	case strings.HasPrefix(res, "[") && strings.HasSuffix(res, "]"):
		for _, val := range v.([]interface{}) {
			if err := m.assign(cmd, val, isArray, keepValue); err != nil {
				return err
			}
		}
	default:
		if err := m.assign(cmd, v, isScalar, keepValue); err != nil {
			return err
		}
	}
	return nil
}

const (
	ignoreValue = true
	keepValue   = false
	isArray     = true
	isScalar    = false
)

func (m *mapper) assign(cmd string, v interface{}, array bool, ignoreValue bool) error {
	switch v := v.(type) {
	case map[string]interface{}, []interface{}:
		if err := m.mapObj(false, v, cmd); err != nil {
			return err
		}
	case string:
		if array || !ignoreValue {
			m.cmds = append(m.cmds, cmd+" "+v)
		} else {
			m.cmds = append(m.cmds, cmd)
		}
	default:
		m.cmds = append(m.cmds, cmd+" "+v.(string))
	}

	return nil
}

// buildCmd
func (m *mapper) buildCmd(array bool, prefix, value string) string {
	if array {
		prefix += value
	} else {
		prefix += " " + value
	}

	return prefix
}
