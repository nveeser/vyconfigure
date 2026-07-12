package config

import (
	"fmt"
	"github.com/nveeser/vyconfigure/section"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	goyaml "gopkg.in/yaml.v3"
)

type Repo struct {
	ConfigDirectory string
	SectionMapper   *section.Mapper
	Ignore          []string
}

// WriteSection writes existing VyOS config section to the local filesystem
func (r *Repo) WriteSection(e *section.Section) error {
	cfgDir, err := r.configDir()
	if err != nil {
		return err
	}
	if e.Basename == "" {
		log.Printf("Ignored path: %s", e.YAMLPath)
		return nil
	}
	if e.StoredContents() == nil {
		log.Printf("No Value path: %s", e.YAMLPath)
		return nil
	}

	p := path.Join(cfgDir, e.Basename+".yaml")
	fp, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := goyaml.NewEncoder(fp)
	defer encoder.Close()
	encoder.SetIndent(2)
	if err := encoder.Encode(e.StoredContents()); err != nil {
		return err
	}
	return nil
}

// WriteConfig writes existing VyOS config to the local filesystem
func (r *Repo) WriteConfig(data map[string]any) error {
	for _, e := range r.SectionMapper.Split(data) {
		err := r.WriteSection(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) ReadSections() ([]*section.Section, error) {
	cfgDir, err := r.configDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(cfgDir)
	if err != nil {
		return nil, err
	}
	var sections []*section.Section
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		basename := strings.TrimSuffix(f.Name(), ".yaml")
		if slices.Contains(r.Ignore, basename) {
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
		sections = append(sections, r.SectionMapper.FromFile(basename, contents))
	}
	return sections, nil
}
func (r *Repo) ReadConfig() (map[string]any, error) {
	sections, err := r.ReadSections()
	if err != nil {
		return nil, err
	}
	return r.SectionMapper.Merge(sections), nil
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
