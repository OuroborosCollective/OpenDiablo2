package types

import "testing"

func TestExt2SourceType(t *testing.T) {
	tests := []struct {
		ext      string
		expected SourceType
	}{
		{"mpq", AssetSourceMPQ},
		{".mpq", AssetSourceMPQ},
		{"MPQ", AssetSourceMPQ},
		{".MPQ", AssetSourceMPQ},
		{"txt", AssetSourceUnknown},
		{"", AssetSourceUnknown},
		{".", AssetSourceUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := Ext2SourceType(tt.ext)
			if got != tt.expected {
				t.Errorf("Ext2SourceType(%q) = %v; want %v", tt.ext, got, tt.expected)
			}
		})
	}
}
