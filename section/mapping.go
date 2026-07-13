package section

import (
	"fmt"
	goyaml "gopkg.in/yaml.v3"
	"iter"
	"maps"
	"os"
	"slices"
)

func NewMapper(configPath string) (*Mapper, error) {
	c, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var mc mapping
	err = goyaml.Unmarshal(c, &mc)
	if err != nil {
		return nil, fmt.Errorf("error converting unmarshalling YAML: %w", err)
	}
	return &Mapper{
		mappings: mc.Mappings,
	}, nil
}

type mapping struct {
	Mappings []MappingEntry `yaml:"Mappings"`
}

type MappingEntry struct {
	File     string `yaml:"File"`
	Path     string `yaml:"Path"`
	KeepPath bool   `yaml:"KeepPath"`
}

func Join(from, to []*Section) iter.Seq2[*Section, *Section] {
	fromMap := indexSection(from)
	toMap := indexSection(to)

	return func(yield func(*Section, *Section) bool) {
		for _, k := range slices.Compact(slices.Sorted(maps.Keys(fromMap))) {
			ff := fromMap[k]
			tt := toMap[k]
			if ff == nil {
				ff = cloneSection(tt)
			}
			if tt == nil {
				tt = cloneSection(ff)
			}
			if !yield(ff, tt) {
				return
			}
		}
	}
}

func cloneSection(s *Section) *Section {
	return &Section{
		Basename:        s.Basename,
		YAMLPath:        s.YAMLPath,
		contents:        nil,
		steps:           slices.Clone(s.steps),
		storeWithPrefix: s.storeWithPrefix,
	}
}

func indexSection(f []*Section) map[string]*Section {
	m := make(map[string]*Section, len(f))
	for _, p := range f {
		m[p.Basename] = p
	}
	return m
}
