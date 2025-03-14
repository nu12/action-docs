package helper

import (
	"os"
	"testing"
)

func TestHash(t *testing.T) {
	var data = "data"
	result := Hash(data)
	expected := "8d777f385d3dfec8815d20f7496026dc"
	if expected != result {
		t.Errorf("error: %s", "Hash doesn't match")
	}
}

func TestScanPatternWorkflow(t *testing.T) {
	// create a directory and a file inside it
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/a.yml", []byte(""), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := os.WriteFile(dir+"/b.yml", []byte(""), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	files, err := ScanPattern(dir, ".yml", false)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("error: %s", "Should have 2 files")
	}

	if files[0] != dir+"/a.yml" {
		t.Errorf("error: %s", "File doesn't match")
	}

	if files[1] != dir+"/b.yml" {
		t.Errorf("error: %s", "File doesn't match")
	}
}

func TestScanPatternAction(t *testing.T) {
	// create a directory and a file inside it
	dir := t.TempDir()
	if err := os.Mkdir(dir+"/a", 0755); err != nil {
		t.Errorf("failed to create directory: %v", err)
	}
	if err := os.Mkdir(dir+"/b", 0755); err != nil {
		t.Errorf("failed to create directory: %v", err)
	}
	if err := os.Mkdir(dir+"/c", 0755); err != nil {
		t.Errorf("failed to create directory: %v", err)
	}
	if err := os.Mkdir(dir+"/c/d", 0755); err != nil {
		t.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(dir+"/a/action.yml", []byte(""), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := os.WriteFile(dir+"/b/action.yml", []byte(""), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := os.WriteFile(dir+"/c/d/action.yml", []byte(""), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	files, err := ScanPattern(dir, "action.yml", true)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("error: %s", "Should have 2 files")
	}

	if files[0] != dir+"/a/action.yml" {
		t.Errorf("error: %s", "File doesn't match")
	}

	if files[1] != dir+"/b/action.yml" {
		t.Errorf("error: %s", "File doesn't match")
	}

	if files[2] != dir+"/c/d/action.yml" {
		t.Errorf("error: %s", "File doesn't match")
	}
}
