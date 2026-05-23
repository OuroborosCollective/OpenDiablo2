package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2loader/asset/types"
)

func TestSource_Type(t *testing.T) {
	s := &Source{}
	if s.Type() != types.AssetSourceFileSystem {
		t.Errorf("expected %v, got %v", types.AssetSourceFileSystem, s.Type())
	}
}

func TestSource_PathAndString(t *testing.T) {
	path := "/test/path"
	s := &Source{Root: path}
	if s.Path() != path {
		t.Errorf("expected %v, got %v", path, s.Path())
	}
	if s.String() != path {
		t.Errorf("expected %v, got %v", path, s.String())
	}
}

func TestSource_Open(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "d2source_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	fileName := "test.txt"
	content := "hello world"
	filePath := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	s := &Source{Root: tempDir}

	// Test success
	r, err := s.Open(fileName)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != content {
		t.Errorf("expected %q, got %q", content, string(data))
	}

	// Test failure
	_, err = s.Open("nonexistent.txt")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestSource_Exists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "d2exists_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	fileName := "exists.txt"
	filePath := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	s := &Source{Root: tempDir}

	tests := []struct {
		name     string
		subPath  string
		expected bool
	}{
		{"Existing file", fileName, true},
		{"Non-existent file", "missing.txt", false},
		{"Root directory", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.Exists(tt.subPath)
			if got != tt.expected {
				t.Errorf("Exists(%q) = %v; want %v", tt.subPath, got, tt.expected)
			}
		})
	}
}
