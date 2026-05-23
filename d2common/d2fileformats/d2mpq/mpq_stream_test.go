package d2mpq

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// MockFile implements io.ReadWriteSeeker for testing
type MockFile struct {
	*bytes.Reader
}

func (m *MockFile) Close() error { return nil }

func TestStream_MultiBlockRead(t *testing.T) {
	// Create a mock MPQ stream that spans multiple blocks
	// Block size is 512 bytes (default)
	blockSize := uint32(512)

	mpq := &MPQ{
		header: Header{
			BlockSize: 0, // 0x200 << 0 = 512
		},
	}

	block := &Block{
		UncompressedFileSize: 1500, // Roughly 3 blocks
		FilePosition: 0,
		Flags: 0, // No compression for simple test
	}

	data := make([]byte, 1500)
	for i := range data {
		data[i] = byte(i % 256)
	}

	mpq.file = nil // We'll mock the Read/Seek if needed, but for now let's see if we can use a real file or bytes.Reader

	// Actually, Stream.Read calls v.MPQ.file.Read
	// We need a way to mock that. MPQ.file is *os.File.
	// This is hard to mock without changing the struct.

	// Since I can't easily mock *os.File without a temporary file, let's use a temporary file.
	tmpFile, err := os.CreateTemp("", "mpq_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}

	mpq.file = tmpFile

	s := &Stream{
		MPQ:   mpq,
		Block: block,
		Size:  blockSize,
		Index: 0xFFFFFFFF,
	}

	buffer := make([]byte, 1500)
	n, err := s.Read(buffer, 0, 1500)
	if err != nil && err != io.EOF {
		t.Errorf("unexpected error: %v", err)
	}

	if n != 1500 {
		t.Errorf("expected 1500 bytes, got %d", n)
	}

	if !bytes.Equal(buffer, data) {
		t.Errorf("data mismatch")
	}
}
