package config

import (
	"cmp"
	"slices"
	"strings"
)

// fileMerger collects the contents of the filesystem (filename, yaml content)
// and merges it into a single configuration.
type fileMerger struct {
	mc        *MappingConfig
	repoFiles []*repoFile
}

func (rm *fileMerger) add(basename string, contents any) {
	yamlPath := basename
	if mapping, ok := findMapping(rm.mc, basename); ok {
		yamlPath = mapping.Path
	}
	rm.repoFiles = append(rm.repoFiles, newRepoFile(basename, yamlPath, contents))
}

func (rm *fileMerger) merge() map[string]any {
	out := map[string]any{}
	for _, repoFile := range rm.repoFiles {
		out = mergeMap(out, addPrefix(repoFile.contents, repoFile.steps))
	}
	return out
}

func findMapping(mc *MappingConfig, basename string) (MappingEntry, bool) {
	var zero MappingEntry
	if mc == nil {
		return zero, false
	}
	for _, entry := range mc.Mappings {
		if entry.File == basename {
			return entry, true
		}
	}
	return zero, false
}

// fileSplitter takes a single configuration and splits it each repofile
// (filename, yaml content) to be written to the filesystem
type fileSplitter struct {
	mc *MappingConfig
}

func (rm *fileSplitter) split(data map[string]any) []*repoFile {
	var sections []*repoFile
	for _, entry := range rm.mc.Mappings {
		sections = append(sections, newRepoFile(entry.File, entry.Path, nil))
	}
	slices.SortFunc(sections, compareFileEntry)
	for _, sf := range sections {
		data, sf.contents = splitMap(data, sf.steps, false)
	}
	for k, v := range data {
		sections = append(sections, newRepoFile(k, k, v))
	}
	return sections
}

type repoFile struct {
	basename string
	yamlPath string
	steps    []string
	contents any
}

func newRepoFile(basename string, yamlPath string, contents any) *repoFile {
	return &repoFile{
		basename: basename,
		yamlPath: yamlPath,
		steps:    strings.Split(yamlPath, "."),
		contents: contents,
	}
}

func compareFileEntry(a, b *repoFile) int {
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

func splitMap(m map[string]any, path []string, keepPrefix bool) (aa map[string]any, sub any) {
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
		_, bb := splitMap(vv, remain, keepPrefix)
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
