package d2mpq

import (
	"testing"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data/d2compression"
)

func TestHuffmanDecompress_Invalid(t *testing.T) {
	// Passing random garbage to huffman should not panic the whole app
	data := []byte{1, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA}
	res := d2compression.HuffmanDecompress(data)
	// It might return nil or some garbage, but it shouldn't panic.
	t.Logf("Decompressed size: %d", len(res))
}
