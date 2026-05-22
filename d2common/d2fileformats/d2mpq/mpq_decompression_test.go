package d2mpq

import (
	"bytes"
	"compress/zlib"
	"testing"
)

func TestDecompressMulti_ZLib(t *testing.T) {
	// Create some zlib compressed data
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	originalData := []byte("hello world multi-block compression test data")
	_, _ = w.Write(originalData)
	_ = w.Close()

	// Prepend compression type 2 (ZLib)
	compressedData := append([]byte{2}, b.Bytes()...)

	decompressed, err := decompressMulti(compressedData, uint32(len(originalData)))
	if err != nil {
		t.Fatalf("decompressMulti failed: %v", err)
	}

	if !bytes.Equal(decompressed, originalData) {
		t.Errorf("decompressed data mismatch. got %s, want %s", string(decompressed), string(originalData))
	}
}

func TestDecompressMulti_InvalidType(t *testing.T) {
	_, err := decompressMulti([]byte{0xFF, 0x00, 0x01}, 10)
	if err == nil {
		t.Error("expected error for unknown compression type, got nil")
	}
}

func TestDecompressMulti_EmptyData(t *testing.T) {
	_, err := decompressMulti([]byte{}, 10)
	if err == nil {
		t.Error("expected error for empty data, got nil")
	}
}

func TestDecompressMulti_SizeMismatch(t *testing.T) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	originalData := []byte("data")
	_, _ = w.Write(originalData)
	_ = w.Close()

	compressedData := append([]byte{2}, b.Bytes()...)

	_, err := decompressMulti(compressedData, 100) // Wrong expected size
	if err == nil {
		t.Error("expected error for size mismatch, got nil")
	}
}
