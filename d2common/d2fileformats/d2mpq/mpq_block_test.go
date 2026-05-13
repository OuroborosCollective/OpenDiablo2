package d2mpq

import (
	"testing"
)

func TestBlock_HasFlag(t *testing.T) {
	tests := []struct {
		name     string
		flags    FileFlag
		flag     FileFlag
		expected bool
	}{
		{
			name:     "Flag present",
			flags:    FileEncrypted | FileCompress,
			flag:     FileEncrypted,
			expected: true,
		},
		{
			name:     "Flag not present",
			flags:    FileEncrypted | FileCompress,
			flag:     FileImplode,
			expected: false,
		},
		{
			name:     "Multiple flags present",
			flags:    FileEncrypted | FileCompress | FileExists,
			flag:     FileExists,
			expected: true,
		},
		{
			name:     "No flags set",
			flags:    0,
			flag:     FileEncrypted,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{Flags: tt.flags}
			if got := b.HasFlag(tt.flag); got != tt.expected {
				t.Errorf("Block.HasFlag() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBlock_calculateEncryptionSeed(t *testing.T) {
	tests := []struct {
		name                 string
		fileName             string
		filePosition         uint32
		uncompressedFileSize uint32
		expectedSeed         uint32
	}{
		{
			name:                 "(listfile)",
			fileName:             "(listfile)",
			filePosition:         100,
			uncompressedFileSize: 500,
			expectedSeed:         0x2D2F0B0C,
		},
		{
			name:                 "armor.txt with path",
			fileName:             `data\global\excel\armor.txt`,
			filePosition:         100,
			uncompressedFileSize: 500,
			expectedSeed:         0x2A266025,
		},
		{
			name:                 "armor.txt without path",
			fileName:             "armor.txt",
			filePosition:         100,
			uncompressedFileSize: 500,
			expectedSeed:         0x2A266025,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				FilePosition:         tt.filePosition,
				UncompressedFileSize: tt.uncompressedFileSize,
			}
			b.calculateEncryptionSeed(tt.fileName)
			if b.EncryptionSeed != tt.expectedSeed {
				t.Errorf("Block.calculateEncryptionSeed() = 0x%08X, want 0x%08X", b.EncryptionSeed, tt.expectedSeed)
			}
		})
	}
}
