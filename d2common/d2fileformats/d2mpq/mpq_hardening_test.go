package d2mpq

import (
	"bytes"
	"compress/zlib"
	"testing"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2compression"
)

func TestDecompressMulti_ZLib_Correct(t *testing.T) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	originalData := []byte("expected decompressed content")
	_, _ = w.Write(originalData)
	_ = w.Close()

	compressedData := append([]byte{2}, b.Bytes()...)

	decompressed, err := decompressMulti(compressedData, uint32(len(originalData)))
	if err != nil {
		t.Fatalf("decompressMulti failed: %v", err)
	}

	if !bytes.Equal(decompressed, originalData) {
		t.Errorf("got %s, want %s", string(decompressed), string(originalData))
	}
}

func TestHuffmanDecompress_Hardening(t *testing.T) {
	// 0x01 is Huffman. The rest is invalid bitstream data.
	data := []byte{1, 0xFF, 0xFF, 0xFF, 0xFF}

	// Should not panic and should handle invalid data gracefully
	// We wrap in a check to ensure it doesn't crash the test runner.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("HuffmanDecompress panicked with invalid data: %v", r)
		}
	}()

	res := d2compression.HuffmanDecompress(data)

	// With random garbage, it either returns nil or some garbage bytes.
	// The primary goal is stability.
	t.Logf("Decompressed length: %d", len(res))
}

func TestDecompressMulti_Hardened_Errors(t *testing.T) {
	// Test unknown compression type
	_, err := decompressMulti([]byte{0x99, 0x00}, 10)
	if err == nil {
		t.Error("expected error for unknown compression type 0x99")
	}

	// Test empty data
	_, err = decompressMulti([]byte{}, 0)
	if err == nil {
		t.Error("expected error for empty data")
	}
}
