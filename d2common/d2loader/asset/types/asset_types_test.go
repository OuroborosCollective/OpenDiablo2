package types

import "testing"

func TestExt2AssetType(t *testing.T) {
	tests := []struct {
		ext      string
		expected AssetType
	}{
		{"json", AssetTypeJSON},
		{".json", AssetTypeJSON},
		{"JSON", AssetTypeJSON},
		{"tbl", AssetTypeStringTable},
		{"txt", AssetTypeDataDictionary},
		{"dat", AssetTypePalette},
		{"pl2", AssetTypePaletteTransform},
		{"cof", AssetTypeCOF},
		{"dc6", AssetTypeDC6},
		{"dcc", AssetTypeDCC},
		{"ds1", AssetTypeDS1},
		{"dt1", AssetTypeDT1},
		{"wav", AssetTypeWAV},
		{"d2", AssetTypeD2},
		{"unknown", AssetTypeUnknown},
		{"", AssetTypeUnknown},
		{".", AssetTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := Ext2AssetType(tt.ext)
			if got != tt.expected {
				t.Errorf("Ext2AssetType(%q) = %v; want %v", tt.ext, got, tt.expected)
			}
		})
	}
}
