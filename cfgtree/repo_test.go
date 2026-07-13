package cfgtree

import (
	"os"
	"path/filepath"
	"testing"

	diffcmp "github.com/google/go-cmp/cmp"
	"github.com/nveeser/vyconfigure/section"
)

func TestRepo_WriteAndReadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	r := &Repo{
		ConfigDirectory: tmpDir,
		//		SectionMapper:   &section.Mapper{},
	}

	testData := map[string]any{
		"system": map[string]any{
			"host-name":   "vyos",
			"domain-name": "local",
		},
		"interfaces": map[string]any{
			"ethernet": map[string]any{
				"eth0": map[string]any{
					"address": "dhcp",
				},
			},
		},
	}

	err := r.WriteConfigTree(testData)
	if err != nil {
		t.Fatalf("WriteConfigTree() unexpected error: %v", err)
	}

	systemFile := filepath.Join(tmpDir, "system.yaml")
	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		t.Errorf("system.yaml was not created")
	}

	interfacesFile := filepath.Join(tmpDir, "interfaces.yaml")
	if _, err := os.Stat(interfacesFile); os.IsNotExist(err) {
		t.Errorf("interfaces.yaml was not created")
	}

	err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("should be ignored"), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy text file: %v", err)
	}

	readData, err := r.ReadConfigTree()
	if err != nil {
		t.Fatalf("ReadConfigTree() unexpected error: %v", err)
	}

	if len(readData) != 2 {
		t.Errorf("ReadConfigTree() returned %d keys, want 2", len(readData))
	}

	if _, ok := readData["readme"]; ok {
		t.Errorf("ReadConfigTree() should not have read readme.txt")
	}

	if diff := diffcmp.Diff(testData, readData); diff != "" {
		t.Errorf("ReadConfigTree() data mismatch (-want +got):\n%s", diff)
	}
}

func TestRepo_ReadConfig_Error(t *testing.T) {
	r := &Repo{
		ConfigDirectory: "/nonexistent/dir/that/does/not/exist",
		SectionMapper:   &section.Mapper{},
	}
	_, err := r.ReadConfigTree()
	if err == nil {
		t.Error("expected error when reading from non-existent directory, got nil")
	}
}

func TestRepo_ReadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "invalid.yaml"), []byte("invalid-yaml : : :"), 0644)
	if err != nil {
		t.Fatalf("failed to write invalid yaml: %v", err)
	}

	r := &Repo{
		ConfigDirectory: tmpDir,
		SectionMapper:   &section.Mapper{},
	}
	_, err = r.ReadConfigTree()
	if err == nil {
		t.Error("expected error when reading invalid yaml file, got nil")
	}
}

func TestRepo_WriteConfig_Error(t *testing.T) {
	r := &Repo{
		ConfigDirectory: "/nonexistent/dir/that/does/not/exist",
		SectionMapper:   &section.Mapper{},
	}
	err := r.WriteConfigTree(map[string]any{"test": "data"})
	if err == nil {
		t.Error("expected error when writing to non-existent directory, got nil")
	}
}
