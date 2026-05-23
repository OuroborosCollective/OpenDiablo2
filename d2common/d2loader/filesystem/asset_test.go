package filesystem

import (
	"os"
	"testing"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2loader/asset/types"
)

func TestAsset_Metadata(t *testing.T) {
	source := &Source{Root: "testdata"}
	path := "test.txt"
	assetType := types.AssetTypeDataDictionary

	f, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	asset := &Asset{
		assetType: assetType,
		source:    source,
		path:      path,
		file:      f,
	}

	if asset.Type() != assetType {
		t.Errorf("expected type %v, got %v", assetType, asset.Type())
	}

	if asset.Source() != source {
		t.Errorf("expected source %v, got %v", source, asset.Source())
	}

	if asset.Path() != path {
		t.Errorf("expected path %q, got %q", path, asset.Path())
	}

	if asset.String() != path {
		t.Errorf("expected string %q, got %q", path, asset.String())
	}
}

func TestAsset_IO(t *testing.T) {
	source := &Source{Root: "testdata"}
	path := "test.txt"
	assetType := types.AssetTypeDataDictionary

	f, err := os.Open("testdata/test.txt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer f.Close()

	asset := &Asset{
		assetType: assetType,
		source:    source,
		path:      path,
		file:      f,
	}

	// Test Read
	buf := make([]byte, 5)
	n, err := asset.Read(buf)
	if err != nil {
		t.Errorf("read failed: %v", err)
	}
	if n != 5 {
		t.Errorf("expected to read 5 bytes, got %d", n)
	}
	if string(buf) != "hello" {
		t.Errorf("expected \"hello\", got %q", string(buf))
	}

	// Test Seek
	pos, err := asset.Seek(6, 0)
	if err != nil {
		t.Errorf("seek failed: %v", err)
	}
	if pos != 6 {
		t.Errorf("expected position 6, got %d", pos)
	}

	n, err = asset.Read(buf)
	if err != nil {
		t.Errorf("read after seek failed: %v", err)
	}
	if string(buf[:5]) != "world" {
		t.Errorf("expected \"world\", got %q", string(buf[:5]))
	}

	// Test Data
	data, err := asset.Data()
	if err != nil {
		t.Errorf("data failed: %v", err)
	}
	if string(data) != "hello world\n" {
		t.Errorf("expected \"hello world\n\", got %q", string(data))
	}

	// Test Data cache
	data2, err := asset.Data()
	if err != nil {
		t.Errorf("data second call failed: %v", err)
	}
	if &data[0] != &data2[0] {
		t.Errorf("expected cached data")
	}

	// Test Close
	err = asset.Close()
	if err != nil {
		t.Errorf("close failed: %v", err)
	}
}

func TestAsset_Data_NoFile(t *testing.T) {
	asset := &Asset{path: "missing.txt"}
	_, err := asset.Data()
	if err == nil {
		t.Error("expected error when calling Data() on asset with no file")
	}
}
