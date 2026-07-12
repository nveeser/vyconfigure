package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	goyaml "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

const defaultMappingFile = "vysync.yaml"

type Repo struct {
	ConfigDirectory string
	MappingFile     string
}

type MappingConfig struct {
	ConfigRoot string         `yaml:"config_root"`
	Mappings   []MappingEntry `yaml:"mappings"`
}

type MappingEntry struct {
	File string `yaml:"file"`
	Path string `yaml:"path"`
}

// WriteConfig writes existing VyOS config to the local filesystem
func (r *Repo) WriteConfig(data map[string]any) error {
	mc, err := r.readMappingConfig()
	if err != nil {
		return err
	}
	splitter := &fileSplitter{mc}

	cfgDir, err := r.configDir()
	if err != nil {
		return err
	}
	for _, e := range splitter.split(data) {
		if e.basename == "" {
			log.Printf("Ignored path: %s", e.yamlPath)
			continue
		}
		if e.contents == nil {
			log.Printf("No Value path: %s", e.yamlPath)
			continue
		}

		p := path.Join(cfgDir, e.basename+".yaml")
		fp, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		encoder := goyaml.NewEncoder(fp)
		defer encoder.Close()
		encoder.SetIndent(2)
		if err := encoder.Encode(e.contents); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) ReadConfig() (map[string]any, error) {
	cfgDir, err := r.configDir()
	if err != nil {
		return nil, err
	}

	mc, err := r.readMappingConfig()
	if err != nil {
		return nil, err
	}
	merger := &fileMerger{mc: mc}

	files, err := os.ReadDir(cfgDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		if f.Name() == r.MappingFile {
			continue
		}
		fp := path.Join(cfgDir, f.Name())
		c, err := os.ReadFile(fp)
		if err != nil {
			return nil, err
		}
		var contents map[string]any
		err = goyaml.Unmarshal(c, &contents)
		if err != nil {
			return nil, fmt.Errorf("error converting unmarshalling YAML: %w", err)
		}
		filename := strings.TrimSuffix(f.Name(), ".yaml")
		merger.add(filename, contents)
	}

	return merger.merge(), nil
}

func (r *Repo) configDir() (string, error) {
	p := r.ConfigDirectory
	if filepath.IsAbs(p) {
		return p, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(wd, p), nil
}

func (r *Repo) readMappingConfig() (*MappingConfig, error) {
	cfgDir, err := r.configDir()
	if err != nil {
		return nil, err
	}
	if r.MappingFile == "" {
		r.MappingFile = defaultMappingFile
	}
	c, err := os.ReadFile(path.Join(cfgDir, r.MappingFile))
	if errors.Is(err, os.ErrNotExist) {
		return &MappingConfig{
			ConfigRoot: r.ConfigDirectory,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	var mc *MappingConfig
	err = yaml.Unmarshal(c, &mc)
	if err != nil {
		return nil, fmt.Errorf("error converting unmarshalling YAML: %w", err)
	}
	mc.ConfigRoot = r.ConfigDirectory
	return mc, nil
}
