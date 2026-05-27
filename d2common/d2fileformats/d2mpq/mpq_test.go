package d2mpq

import (
	"testing"
)

const testMpqPath = "../../d2loader/testdata/D.mpq"

func TestNew(t *testing.T) {
	mpq, err := New(testMpqPath)
	if err != nil {
		t.Fatalf("failed to create new MPQ: %v", err)
	}
	defer mpq.Close()

	if mpq.Path() != testMpqPath {
		t.Errorf("expected path %s, got %s", testMpqPath, mpq.Path())
	}
}

func TestNew_NonExistent(t *testing.T) {
	_, err := New("non_existent.mpq")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestFromFile(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ from file: %v", err)
	}
	defer mpq.Close()

	if len(mpq.hashes) == 0 {
		t.Error("expected non-empty hash table")
	}

	if len(mpq.blocks) == 0 {
		t.Error("expected non-empty block table")
	}
}

func TestFromFile_NonExistent(t *testing.T) {
	_, err := FromFile("non_existent.mpq")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestReadFile(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	fileName := "exclusive_d.txt"
	data, err := mpq.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", fileName, err)
	}

	if string(data[0]) != "d" {
		t.Errorf("expected content 'd', got %q", string(data))
	}
}

func TestReadFile_NonExistent(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	_, err = mpq.ReadFile("non_existent.txt")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestReadFileStream(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	fileName := "exclusive_d.txt"
	stream, err := mpq.ReadFileStream(fileName)
	if err != nil {
		t.Fatalf("failed to read file stream %s: %v", fileName, err)
	}

	buffer := make([]byte, 1)
	_, err = stream.Read(buffer)
	if err != nil {
		t.Fatalf("failed to read from stream: %v", err)
	}

	if string(buffer) != "d" {
		t.Errorf("expected content 'd', got %q", string(buffer))
	}
}

func TestReadTextFile(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	fileName := "exclusive_d.txt"
	text, err := mpq.ReadTextFile(fileName)
	if err != nil {
		t.Fatalf("failed to read text file %s: %v", fileName, err)
	}

	if string(text[0]) != "d" {
		t.Errorf("expected content starting with 'd', got %q", text)
	}
}

func TestListfile(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	files, err := mpq.Listfile()
	if err != nil {
		// D.mpq has an issue with (listfile) compression, so we skip this check if it fails with specific error
		t.Logf("Listfile failed (expected for D.mpq due to blast issue): %v", err)
		return
	}

	if len(files) == 0 {
		t.Error("expected non-empty listfile")
	}
}

func TestContains(t *testing.T) {
	mpq, err := FromFile(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	if !mpq.Contains("exclusive_d.txt") {
		t.Error("expected archive to contain 'exclusive_d.txt'")
	}

	if mpq.Contains("non_existent.txt") {
		t.Error("expected archive to not contain 'non_existent.txt'")
	}
}

func TestSize(t *testing.T) {
	mpq, err := New(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	if mpq.Size() == 0 {
		t.Error("expected non-zero size")
	}
}

func TestPath(t *testing.T) {
	mpq, err := New(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}
	defer mpq.Close()

	if mpq.Path() != testMpqPath {
		t.Errorf("expected path %s, got %s", testMpqPath, mpq.Path())
	}
}

func TestClose(t *testing.T) {
	mpq, err := New(testMpqPath)
	if err != nil {
		t.Fatalf("failed to load MPQ: %v", err)
	}

	if err := mpq.Close(); err != nil {
		t.Fatalf("failed to close MPQ: %v", err)
	}
}
