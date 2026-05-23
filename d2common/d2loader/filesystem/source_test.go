package filesystem

import (
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2loader/asset/types"
)

func TestSource_Metadata(t *testing.T) {
	root := "testdata"
	source := &Source{Root: root}

	if source.Type() != types.AssetSourceFileSystem {
		t.Errorf("expected type %v, got %v", types.AssetSourceFileSystem, source.Type())
	}

	if source.Path() != root {
		t.Errorf("expected path %q, got %q", root, source.Path())
	}

	if source.String() != root {
		t.Errorf("expected string %q, got %q", root, source.String())
	}
}

func TestSource_Exists(t *testing.T) {
	source := &Source{Root: "testdata"}

	if !source.Exists("test.txt") {
		t.Error("expected test.txt to exist")
	}

	if source.Exists("nonexistent.txt") {
		t.Error("expected nonexistent.txt to not exist")
	}
}

func TestSource_Open(t *testing.T) {
	source := &Source{Root: "testdata"}

	reader, err := source.Open("test.txt")
	if err != nil {
		t.Fatalf("failed to open test.txt: %v", err)
	}

	buf := make([]byte, 11)
	n, err := reader.Read(buf)
	if err != nil {
		t.Errorf("read failed: %v", err)
	}
	if n != 11 {
		t.Errorf("expected to read 11 bytes, got %d", n)
	}
	if string(buf) != "hello world" {
		t.Errorf("expected \"hello world\", got %q", string(buf))
	}

	_, err = source.Open("nonexistent.txt")
	if err == nil {
		t.Error("expected error when opening nonexistent file")
	}
}
