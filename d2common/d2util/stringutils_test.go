package d2util

import (
	"testing"
)

func TestUtf16BytesToString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:    "empty slice",
			input:   []byte{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "simple string",
			input:   []byte{0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00},
			want:    "Hello",
			wantErr: false,
		},
		{
			name:    "odd length byte slice",
			input:   []byte{0x48},
			want:    "",
			wantErr: true,
		},
		{
			name:    "odd length byte slice longer",
			input:   []byte{0x48, 0x00, 0x65},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Utf16BytesToString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Utf16BytesToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Utf16BytesToString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
