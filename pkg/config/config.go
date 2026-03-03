package config

import (
	"os"
	"path"
	"strings"

	"github.com/nveeser/vyconfigure/pkg/convert"
	"sigs.k8s.io/yaml"
)

type Repo struct {
	ConfigDirectory string
}

// Write writes existing vyos config to the local filesystem
func (r *Repo) Write(data map[string]any) error {
	for k := range data {
		y, err := yaml.Marshal(data[k])
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		p := path.Join(wd, r.ConfigDirectory, k+".yaml")
		err = os.WriteFile(p, y, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadAsCmds reads all yaml configuration and converts it into vyos "set" commands
func (r *Repo) ReadAsCmds() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	p := path.Join(wd, r.ConfigDirectory)
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		fp := path.Join(wd, r.ConfigDirectory, f.Name())

		c, err := os.ReadFile(fp)
		if err != nil {
			return nil, err
		}

		configPath := strings.TrimSuffix(f.Name(), ".yaml")
		cmds, err := convert.YamlToCmds(c, configPath+" ")
		if err != nil {
			return nil, err
		}
		res = append(res, cmds...)
	}

	return res, nil
}
