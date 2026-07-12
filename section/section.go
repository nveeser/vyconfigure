package section

import (
	"cmp"
	"slices"
	"strings"
)

// Mapper wraps a slice of MappingEntry instance and
// provides methods to split the full config into sections
// to store or translate the contents of the filesystem into Sections
// and then merge it into a single configuration.
type Mapper struct {
	mappings []MappingEntry
}

func (m *Mapper) FromFile(basename string, contents any) *Section {
	yamlPath := basename
	keepPrefix := false
	if mapping, ok := findMapping(m, basename); ok {
		yamlPath = mapping.Path
		keepPrefix = mapping.KeepPath
	}
	return &Section{
		Basename:        basename,
		YAMLPath:        yamlPath,
		steps:           strings.Split(yamlPath, "."),
		contents:        contents,
		storeWithPrefix: keepPrefix,
	}
}

func (m *Mapper) Merge(sections []*Section) map[string]any {
	out := map[string]any{}
	for _, section := range sections {
		out = mergeMap(out, section.Contents().(map[string]any))
	}
	return out
}

func findMapping(mc *Mapper, basename string) (MappingEntry, bool) {
	var zero MappingEntry
	if mc == nil {
		return zero, false
	}
	for _, entry := range mc.mappings {
		if entry.File == basename {
			return entry, true
		}
	}
	return zero, false
}

func (m *Mapper) Split(data map[string]any) []*Section {
	if m == nil {
		m = &Mapper{}
	}
	var sections []*Section
	for _, entry := range m.mappings {
		var contents any = nil
		sections = append(sections, &Section{
			Basename:        entry.File,
			YAMLPath:        entry.Path,
			steps:           strings.Split(entry.Path, "."),
			contents:        contents,
			storeWithPrefix: entry.KeepPath,
		})
	}

	slices.SortFunc(sections, compareFileEntry)

	for _, sf := range sections {
		data, sf.contents = SplitMap(data, sf.steps, sf.storeWithPrefix)
	}
	for k, v := range data {
		sections = append(sections, &Section{
			Basename: k,
			YAMLPath: k,
			steps:    strings.Split(k, "."),
			contents: v,
		})
	}
	return sections
}

type Section struct {
	Basename        string
	YAMLPath        string
	contents        any
	steps           []string
	storeWithPrefix bool
}

func (s *Section) Ignore() bool { return s.contents == nil }

func (s *Section) StoredContents() any { return s.contents }

func (s *Section) Contents() any {
	if s.storeWithPrefix {
		return s.contents
	}
	return addPrefix(s.contents, s.steps)
}

func compareFileEntry(a, b *Section) int {
	return -cmp.Compare(len(a.steps), len(b.steps))
}

func addPrefix(m any, path []string) map[string]any {
	for _, step := range slices.Backward(path) {
		m = map[string]any{
			step: m,
		}
	}
	return m.(map[string]any)
}

func SplitMap(m map[string]any, path []string, keepPrefix bool) (aa map[string]any, sub any) {
	// foo.bar.baz => foo, bar.baz
	currKey, remain := path[0], path[1:]
	v, ok := m[currKey]
	if !ok {
		return m, nil
	}
	if len(remain) == 0 {
		delete(m, currKey)
		if keepPrefix {
			return m, map[string]any{currKey: v}
		}
		return m, v
	}
	if vv, ok := v.(map[string]any); ok {
		_, bb := SplitMap(vv, remain, keepPrefix)
		if keepPrefix {
			return m, map[string]any{currKey: bb}
		}
		return m, bb
	}
	return m, nil
}

func mergeMap(dest, src map[string]any) map[string]any {
	if dest == nil {
		dest = make(map[string]any)
	}
	for k, v := range src {
		if destVal, ok := dest[k]; ok {
			if destMap, ok1 := destVal.(map[string]any); ok1 {
				if srcMap, ok2 := v.(map[string]any); ok2 {
					dest[k] = mergeMap(destMap, srcMap)
					continue
				}
			}
		}
		dest[k] = v
	}
	return dest
}
