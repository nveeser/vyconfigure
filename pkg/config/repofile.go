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
	keepPrefix := false
	if mapping, ok := findMapping(rm.mc, basename); ok {
		yamlPath = mapping.Path
		keepPrefix = mapping.KeepPath
	}
	rm.repoFiles = append(rm.repoFiles, &repoFile{
		basename:   basename,
		yamlPath:   yamlPath,
		steps:      strings.Split(yamlPath, "."),
		contents:   contents,
		keepPrefix: keepPrefix,
	})
}

func (rm *fileMerger) merge() map[string]any {
	out := map[string]any{}
	for _, repoFile := range rm.repoFiles {
		var contents any = repoFile.contents
		if !repoFile.keepPrefix {
			contents = addPrefix(repoFile.contents, repoFile.steps)
		}
		out = mergeMap(out, contents.(map[string]interface{}))
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
		var contents any = nil
		sections = append(sections, &repoFile{
			basename:   entry.File,
			yamlPath:   entry.Path,
			steps:      strings.Split(entry.Path, "."),
			contents:   contents,
			keepPrefix: entry.KeepPath,
		})
	}
	slices.SortFunc(sections, compareFileEntry)
	for _, sf := range sections {
		data, sf.contents = splitMap(data, sf.steps, sf.keepPrefix)
	}
	for k, v := range data {
		sections = append(sections, &repoFile{
			basename: k,
			yamlPath: k,
			steps:    strings.Split(k, "."),
			contents: v,
		})
	}
	return sections
}

type repoFile struct {
	basename   string
	yamlPath   string
	steps      []string
	keepPrefix bool
	contents   any
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
