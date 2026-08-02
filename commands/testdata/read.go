package testdata

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// PackageDir holds the absolute path to the package's directory.
var testDataPath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not get current filename")
	}
	testDataPath = filepath.Dir(filename)
}

func Read(t *testing.T, filename string) []byte {
	t.Helper()
	d, err := os.ReadFile(filepath.Join(testDataPath, filename))
	if err != nil {
		t.Logf("PackageDir: %s", testDataPath)
		t.Fatalf("Error reading filename: %s, ", filepath.Join("./testdata", filename))
	}
	return d
}
func ReadYAML(t *testing.T, filename string) map[string]any {
	t.Helper()
	data := Read(t, filename)
	var out map[string]any
	err := yaml.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("Error parsing YAML: %s, ", filepath.Join("./testdata", filename))
	}
	return out
}
